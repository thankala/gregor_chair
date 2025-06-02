package main

import (
	"os"

	"github.com/anthdm/hollywood/actor"
	AssemblyTask1Service "github.com/thankala/gregor_chair/assembly_task_1/services"
	AssemblyTask2Service "github.com/thankala/gregor_chair/assembly_task_2/services"
	AssemblyTask3Service "github.com/thankala/gregor_chair/assembly_task_3/services"
	AssemblyTask4Service "github.com/thankala/gregor_chair/assembly_task_4/services"
	AssemblyTask5Service "github.com/thankala/gregor_chair/assembly_task_5/services"
	AssemblyTask6Service "github.com/thankala/gregor_chair/assembly_task_6/services"
	AssemblyTask7Service "github.com/thankala/gregor_chair/assembly_task_7/services"
	AssemblyTask8Service "github.com/thankala/gregor_chair/assembly_task_8/services"
	OrchestratorService "github.com/thankala/gregor_chair/orchestrator/services"
	"github.com/thankala/gregor_chair_common/configuration"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/logger"
	"github.com/thankala/gregor_chair_common/services"
)

type Orchestrator struct {
	Robot1Controller     *controllers.RobotController
	Robot2Controller     *controllers.RobotController
	Robot3Controller     *controllers.RobotController
	Workbench1Controller *controllers.WorkbenchController
	Workbench2Controller *controllers.WorkbenchController

	children map[enums.Task]*actor.PID
}

func NewOrchestrator(
	robot1Controller *controllers.RobotController,
	robot2Controller *controllers.RobotController,
	robot3Controller *controllers.RobotController,
	workbench1Controller *controllers.WorkbenchController,
	workbench2Controller *controllers.WorkbenchController,
) actor.Producer {
	return func() actor.Receiver {
		return &Orchestrator{
			children:             make(map[enums.Task]*actor.PID),
			Robot1Controller:     robot1Controller,
			Robot2Controller:     robot2Controller,
			Robot3Controller:     robot3Controller,
			Workbench1Controller: workbench1Controller,
			Workbench2Controller: workbench2Controller,
		}
	}
}

func (c *Orchestrator) Receive(ctx *actor.Context) {
	switch event := ctx.Message().(type) {
	case actor.Initialized:
		c.children[enums.Orchestrator] = ctx.SpawnChild(
			services.NewOrchestratorActor[OrchestratorService.OrchestratorActor](
				OrchestratorService.NewOrchestratorActor(
					[]controllers.WorkbenchController{*c.Workbench1Controller, *c.Workbench2Controller},
					[]controllers.RobotController{*c.Robot1Controller, *c.Robot2Controller, *c.Robot3Controller}),
				nil),
			enums.Orchestrator.String(),
		)
		c.children[enums.AssemblyTask1] = ctx.SpawnChild(services.NewAssemblyTaskActor(
			AssemblyTask1Service.NewAssemblyTask1Actor(*c.Robot1Controller), nil),
			enums.AssemblyTask1.String(),
		)
		c.children[enums.AssemblyTask2] = ctx.SpawnChild(services.NewAssemblyTaskActor(
			AssemblyTask2Service.NewAssemblyTask2Actor(*c.Robot2Controller), nil),
			enums.AssemblyTask2.String(),
		)
		c.children[enums.AssemblyTask3] = ctx.SpawnChild(services.NewAssemblyTaskActor(
			AssemblyTask3Service.NewAssemblyTask3Actor(*c.Robot2Controller), nil),
			enums.AssemblyTask3.String(),
		)
		c.children[enums.AssemblyTask4] = ctx.SpawnChild(services.NewAssemblyTaskActor(
			AssemblyTask4Service.NewAssemblyTask4Actor(*c.Robot1Controller), nil),
			enums.AssemblyTask4.String(),
		)
		c.children[enums.AssemblyTask5] = ctx.SpawnChild(services.NewAssemblyTaskActor(
			AssemblyTask5Service.NewAssemblyTask5Actor(*c.Robot2Controller), nil),
			enums.AssemblyTask5.String(),
		)
		c.children[enums.AssemblyTask6] = ctx.SpawnChild(services.NewAssemblyTaskActor(
			AssemblyTask6Service.NewAssemblyTask6Actor(*c.Robot3Controller), nil),
			enums.AssemblyTask6.String(),
		)
		c.children[enums.AssemblyTask7] = ctx.SpawnChild(services.NewAssemblyTaskActor(
			AssemblyTask7Service.NewAssemblyTask7Actor(*c.Robot3Controller), nil),
			enums.AssemblyTask7.String(),
		)
		c.children[enums.AssemblyTask8] = ctx.SpawnChild(services.NewAssemblyTaskActor(
			AssemblyTask8Service.NewAssemblyTask8Actor(*c.Robot3Controller), nil),
			enums.AssemblyTask8.String(),
		)
	case actor.Started:
		pid := c.children[enums.Orchestrator]
		ctx.Send(pid, &events.OrchestratorEvent{
			Destination: enums.Orchestrator,
			Type:        enums.AssemblyStarted,
		})
	case actor.Stopped:
		break
	case *events.OrchestratorEvent:
		destination := c.children[event.Destination]
		ctx.Send(destination, event)
	case *events.AssemblyTaskEvent:
		destination := c.children[event.Destination]
		ctx.Send(destination, event)
	default:
		logger.Get().Warn("Unknown event received", event)
		return
	}
}

func main() {
	var robot1HttpClient interfaces.HttpClient
	var robot2HttpClient interfaces.HttpClient
	var robot3HttpClient interfaces.HttpClient
	var workbench1HttpClient interfaces.HttpClient
	var workbench2HttpClient interfaces.HttpClient

	if os.Getenv("ROBOT_1_URL") != "" {
		robot1HttpClient = services.NewHttpClient(os.Getenv("ROBOT_1_URL"))
	}

	if os.Getenv("ROBOT_2_URL") != "" {
		robot2HttpClient = services.NewHttpClient(os.Getenv("ROBOT_2_URL"))
	}

	if os.Getenv("ROBOT_3_URL") != "" {
		robot3HttpClient = services.NewHttpClient(os.Getenv("ROBOT_3_URL"))
	}

	if os.Getenv("WORKBENCH_1_URL") != "" {
		workbench1HttpClient = services.NewHttpClient(os.Getenv("WORKBENCH_1_URL"))
	}

	if os.Getenv("WORKBENCH_2_URL") != "" {
		workbench2HttpClient = services.NewHttpClient(os.Getenv("WORKBENCH_2_URL"))
	}

	redisStorer := services.NewRedisStore(
		configuration.WithStoreAddr("localhost:6379"),
		configuration.WithStorePassword(""),
		configuration.WithStoreDb("0"),
	)

	robot1Controller := controllers.NewRobotController(
		redisStorer,
		robot1HttpClient,
		configuration.WithRobotKey(enums.Robot1),
		configuration.WithStorages(
			*configuration.NewStorageConfiguration(
				enums.StorageB1,
				enums.Position1,
				enums.Legs,
				configuration.NewLocation(float64(210), float64(90), float64(20), float64(0)),
			),
			*configuration.NewStorageConfiguration(
				enums.StorageB2,
				enums.Position1,
				enums.Base,
				configuration.NewLocation(float64(200), float64(180), float64(20), float64(0)),
			),
			*configuration.NewStorageConfiguration(
				enums.StorageB3,
				enums.Position2,
				enums.SeatPlate,
				configuration.NewLocation(float64(290), float64(0), float64(20), float64(0))),
		),
		configuration.WithWorkbenches(
			*configuration.NewWorkbenchConfiguration(
				enums.Workbench1,
				enums.Position1,
				enums.Fixture1,
				configuration.NewLocation(float64(260), float64(110), float64(90), float64(0)),
			),
			*configuration.NewWorkbenchConfiguration(
				enums.Workbench2,
				enums.Position2,
				enums.Fixture1,
				configuration.NewLocation(float64(270), float64(-90), float64(40), float64(0)),
			),
		),
		configuration.WithConveyorBelts(
			*configuration.NewConveyorBeltConfiguration(
				enums.ConveyorBelt1,
				enums.Position2,
				enums.Seat,
				false,
				configuration.NewLocation(float64(230), float64(-120), float64(40), float64(0)),
			),
		),
	)

	robot2Controller := controllers.NewRobotController(
		redisStorer,
		robot2HttpClient,
		configuration.WithRobotKey(enums.Robot2),
		configuration.WithStorages(
			*configuration.NewStorageConfiguration(
				enums.StorageB4,
				enums.Position1,
				enums.Castors,
				configuration.NewLocation(float64(170), float64(-160), float64(20), float64(0)),
			),
			*configuration.NewStorageConfiguration(
				enums.StorageB5,
				enums.Position1,
				enums.Lift,
				configuration.NewLocation(float64(290), float64(20), float64(20), float64(0)),
			),
		),
		configuration.WithWorkbenches(
			*configuration.NewWorkbenchConfiguration(
				enums.Workbench1,
				enums.Position1,
				enums.Fixture2,
				configuration.NewLocation(float64(260), float64(-110), float64(90), float64(0)),
			),
			*configuration.NewWorkbenchConfiguration(
				enums.Workbench2,
				enums.Position2,
				enums.Fixture1,
				configuration.NewLocation(float64(270), float64(80), float64(40), float64(0)),
			),
		),
	)

	robot3Controller := controllers.NewRobotController(
		redisStorer,
		robot3HttpClient,
		configuration.WithRobotKey(enums.Robot3),
		configuration.WithStorages(
			*configuration.NewStorageConfiguration(
				enums.StorageB6L,
				enums.Position1,
				enums.LeftArm,
				configuration.NewLocation(float64(210), float64(170), float64(20), float64(0)),
			),
			*configuration.NewStorageConfiguration(
				enums.StorageB6R,
				enums.Position1,
				enums.RightArm,
				configuration.NewLocation(float64(160), float64(180), float64(20), float64(0)),
			),
		),
		configuration.WithWorkbenches(
			*configuration.NewWorkbenchConfiguration(
				enums.Workbench1,
				enums.Position1,
				enums.Fixture3,
				configuration.NewLocation(float64(260), float64(0), float64(90), float64(0)),
			),
		),
		configuration.WithConveyorBelts(
			*configuration.NewConveyorBeltConfiguration(
				enums.ConveyorBelt2,
				enums.Position1,
				enums.Back,
				false,
				configuration.NewLocation(float64(220), float64(190), float64(20), float64(0)),
			),
			*configuration.NewConveyorBeltConfiguration(
				enums.ConveyorBelt3,
				enums.Position1,
				enums.NoneComponent,
				true,
				configuration.NewLocation(float64(220), float64(190), float64(20), float64(0)),
			),
		),
	)

	workbench1Controller := controllers.NewWorkbenchController(
		redisStorer,
		workbench1HttpClient,
		configuration.WithWorkbenchKey(enums.Workbench1),
		configuration.WithFixture(
			*configuration.NewFixtureConfiguration(
				enums.Fixture1,
				[]string{enums.Robot1.String()}),
			*configuration.NewFixtureConfiguration(
				enums.Fixture2,
				[]string{enums.Robot2.String()}),
			*configuration.NewFixtureConfiguration(
				enums.Fixture3,
				[]string{enums.Robot3.String()}),
		),
	)

	workbench2Controller := controllers.NewWorkbenchController(
		redisStorer,
		workbench2HttpClient,
		configuration.WithWorkbenchKey(enums.Workbench2),
		configuration.WithFixture(
			*configuration.NewFixtureConfiguration(
				enums.Fixture1,
				[]string{enums.Robot1.String(), enums.Robot2.String()},
			),
		),
	)

	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	engine.Spawn(NewOrchestrator(
		robot1Controller,
		robot2Controller,
		robot3Controller,
		workbench1Controller,
		workbench2Controller), "local")
	<-make(chan struct{})
}
