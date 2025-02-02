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
	Coordinator1Service "github.com/thankala/gregor_chair/coordinator_1/services"
	Coordinator2Service "github.com/thankala/gregor_chair/coordinator_2/services"
	"github.com/thankala/gregor_chair_common/configuration"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/logger"
	"github.com/thankala/gregor_chair_common/messages"
	"github.com/thankala/gregor_chair_common/services"
)

type Coordinator struct {
	Robot1Controller     *controllers.RobotController
	Robot2Controller     *controllers.RobotController
	Robot3Controller     *controllers.RobotController
	Workbench1Controller *controllers.WorkbenchController
	Workbench2Controller *controllers.WorkbenchController

	children map[string]*actor.PID
}

func NewCoordinator(
	robot1Controller *controllers.RobotController,
	robot2Controller *controllers.RobotController,
	robot3Controller *controllers.RobotController,
	workbench1Controller *controllers.WorkbenchController,
	workbench2Controller *controllers.WorkbenchController,
) actor.Producer {
	return func() actor.Receiver {
		return &Coordinator{
			children:             make(map[string]*actor.PID),
			Robot1Controller:     robot1Controller,
			Robot2Controller:     robot2Controller,
			Robot3Controller:     robot3Controller,
			Workbench1Controller: workbench1Controller,
			Workbench2Controller: workbench2Controller,
		}
	}
}

func (c *Coordinator) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		c.children[enums.Coordinator1.String()] = ctx.SpawnChild(services.NewCoordinatorActor[Coordinator1Service.Coordinator1Actor](
			Coordinator1Service.NewCoordinator1Actor(*c.Workbench1Controller), nil),
			enums.Coordinator1.String(),
		)
		c.children[enums.Coordinator2.String()] = ctx.SpawnChild(services.NewCoordinatorActor[Coordinator2Service.Coordinator2Actor](
			Coordinator2Service.NewCoordinator2Actor(*c.Workbench2Controller), nil),
			enums.Coordinator2.String(),
		)
		c.children[enums.AssemblyTask1.String()] = ctx.SpawnChild(services.NewAssemblyTaskActor[AssemblyTask1Service.AssemblyTask1Actor](
			AssemblyTask1Service.NewAssemblyTask1Actor(*c.Robot1Controller), nil),
			enums.AssemblyTask1.String(),
		)
		c.children[enums.AssemblyTask2.String()] = ctx.SpawnChild(services.NewAssemblyTaskActor[AssemblyTask2Service.AssemblyTask2Actor](
			AssemblyTask2Service.NewAssemblyTask2Actor(*c.Robot2Controller), nil),
			enums.AssemblyTask2.String(),
		)
		c.children[enums.AssemblyTask3.String()] = ctx.SpawnChild(services.NewAssemblyTaskActor[AssemblyTask3Service.AssemblyTask3Actor](
			AssemblyTask3Service.NewAssemblyTask3Actor(*c.Robot2Controller), nil),
			enums.AssemblyTask3.String(),
		)
		c.children[enums.AssemblyTask4.String()] = ctx.SpawnChild(services.NewAssemblyTaskActor[AssemblyTask4Service.AssemblyTask4Actor](
			AssemblyTask4Service.NewAssemblyTask4Actor(*c.Robot1Controller), nil),
			enums.AssemblyTask4.String(),
		)
		c.children[enums.AssemblyTask5.String()] = ctx.SpawnChild(services.NewAssemblyTaskActor[AssemblyTask5Service.AssemblyTask5Actor](
			AssemblyTask5Service.NewAssemblyTask5Actor(*c.Robot2Controller), nil),
			enums.AssemblyTask5.String(),
		)
		c.children[enums.AssemblyTask6.String()] = ctx.SpawnChild(services.NewAssemblyTaskActor[AssemblyTask6Service.AssemblyTask6Actor](
			AssemblyTask6Service.NewAssemblyTask6Actor(*c.Robot3Controller), nil),
			enums.AssemblyTask6.String(),
		)
		c.children[enums.AssemblyTask7.String()] = ctx.SpawnChild(services.NewAssemblyTaskActor[AssemblyTask7Service.AssemblyTask7Actor](
			AssemblyTask7Service.NewAssemblyTask7Actor(*c.Robot3Controller), nil),
			enums.AssemblyTask7.String(),
		)
		c.children[enums.AssemblyTask8.String()] = ctx.SpawnChild(services.NewAssemblyTaskActor[AssemblyTask8Service.AssemblyTask8Actor](
			AssemblyTask8Service.NewAssemblyTask8Actor(*c.Robot3Controller), nil),
			enums.AssemblyTask8.String(),
		)
	case actor.Started:
		pid := c.children[enums.AssemblyTask1.String()]
		ctx.Send(pid, &messages.AssemblyTaskMessage{
			Source:      enums.NoneAssemblyTask.String(),
			Destination: enums.AssemblyTask1.String(),
			Task:        enums.AssemblyTask1,
			Step:        enums.Step1,
		})
	case actor.Stopped:
		break
	case *messages.CoordinatorMessage:
		destination := c.children[msg.Destination]
		ctx.Send(destination, msg)
	case *messages.AssemblyTaskMessage:
		destination := c.children[msg.Destination]
		ctx.Send(destination, msg)
	default:
		logger.Get().Warn("Unknown event received", msg)
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
		configuration.WithRobotKey(enums.Robot1.String()),
		configuration.WithStorages(
			*configuration.NewStorageConfiguration(enums.StorageB1, enums.Position1, enums.Legs),
			*configuration.NewStorageConfiguration(enums.StorageB2, enums.Position1, enums.Base),
			*configuration.NewStorageConfiguration(enums.StorageB3, enums.Position2, enums.SeatPlate),
			*configuration.NewStorageConfiguration(enums.StorageB7A, enums.Position1, enums.NoneComponent),
			*configuration.NewStorageConfiguration(enums.StorageB7B, enums.Position2, enums.NoneComponent),
		),
		configuration.WithWorkbenches(
			*configuration.NewWorkbenchConfiguration(enums.Workbench1, enums.Position1, enums.Fixture1),
			*configuration.NewWorkbenchConfiguration(enums.Workbench2, enums.Position2, enums.Fixture1),
		),
		configuration.WithConveyorBelts(
			*configuration.NewConveyorBeltConfiguration(
				enums.ConveyorBelt1,
				enums.Position1,
				enums.NoneComponent,
				false,
			),
			*configuration.NewConveyorBeltConfiguration(
				enums.ConveyorBelt2,
				enums.Position2,
				enums.Seat,
				false,
			),
		),
	)

	robot2Controller := controllers.NewRobotController(
		redisStorer,
		robot2HttpClient,
		configuration.WithRobotKey(enums.Robot2.String()),
		configuration.WithStorages(
			*configuration.NewStorageConfiguration(enums.StorageB4, enums.Position1, enums.Castors),
			*configuration.NewStorageConfiguration(enums.StorageB5, enums.Position1, enums.Lift),
		),
		configuration.WithWorkbenches(
			*configuration.NewWorkbenchConfiguration(enums.Workbench1, enums.Position1, enums.Fixture2),
			*configuration.NewWorkbenchConfiguration(enums.Workbench2, enums.Position2, enums.Fixture1),
		),
	)

	robot3Controller := controllers.NewRobotController(
		redisStorer,
		robot3HttpClient,
		configuration.WithRobotKey(enums.Robot3.String()),
		configuration.WithStorages(
			*configuration.NewStorageConfiguration(enums.StorageB6L, enums.Position1, enums.LeftArm),
			*configuration.NewStorageConfiguration(enums.StorageB6R, enums.Position1, enums.RightArm),
			*configuration.NewStorageConfiguration(enums.StorageB8C, enums.Position1, enums.NoneComponent),
			*configuration.NewStorageConfiguration(enums.StorageB8D, enums.Position1, enums.NoneComponent),
			*configuration.NewStorageConfiguration(enums.StorageB8E, enums.Position1, enums.NoneComponent),
		),
		configuration.WithWorkbenches(
			*configuration.NewWorkbenchConfiguration(enums.Workbench1, enums.Position1, enums.Fixture3),
		),
		configuration.WithConveyorBelts(
			*configuration.NewConveyorBeltConfiguration(
				enums.ConveyorBelt2,
				enums.Position1,
				enums.Back,
				false,
			),
			*configuration.NewConveyorBeltConfiguration(enums.ConveyorBelt3,
				enums.Position1,
				enums.NoneComponent,
				true,
			),
		),
	)

	workbench1Controller := controllers.NewWorkbenchController(
		redisStorer,
		workbench1HttpClient,
		configuration.WithWorkbenchKey(enums.Workbench1.String()),
		configuration.WithFixture(
			*configuration.NewFixtureConfiguration(enums.Fixture1, []string{enums.Robot1.String()}),
			*configuration.NewFixtureConfiguration(enums.Fixture2, []string{enums.Robot2.String()}),
			*configuration.NewFixtureConfiguration(enums.Fixture3, []string{enums.Robot3.String()}),
		),
		configuration.WithStateMapping(map[enums.Fixture]map[enums.Stage]string{
			enums.Fixture1: {
				enums.Initial:      "FREE",
				enums.LegsAttached: "ASSEMBLING",
				enums.BaseAttached: "COMPLETED",
			},
			enums.Fixture2: {
				enums.Initial:         "FREE",
				enums.BaseAttached:    "ASSEMBLING",
				enums.CastorsAttached: "ASSEMBLING",
				enums.LiftAttached:    "PENDING",
				enums.SeatAttached:    "ASSEMBLING",
				enums.ScrewsAttached:  "COMPLETED",
			},
			enums.Fixture3: {
				enums.Initial:        "FREE",
				enums.ScrewsAttached: "ASSEMBLING",
				// enums.BackAttached:     "COMPLETED",
				enums.LeftArmAttached:  "ASSEMBLING",
				enums.RightArmAttached: "ASSEMBLING",
				enums.Completed:        "COMPLETED",
			},
		}),
	)

	workbench2Controller := controllers.NewWorkbenchController(
		redisStorer,
		workbench2HttpClient,
		configuration.WithWorkbenchKey(enums.Workbench2.String()),
		configuration.WithFixture(
			*configuration.NewFixtureConfiguration(
				enums.Fixture1,
				[]string{enums.Robot1.String(), enums.Robot2.String()},
			),
		),
		configuration.WithStateMapping(map[enums.Fixture]map[enums.Stage]string{
			enums.Fixture1: {
				enums.Initial:            "FREE",
				enums.InitialSeat:        "ASSEMBLING",
				enums.SeatPlateAttached:  "PENDING",
				enums.SeatScrewsAttached: "ASSEMBLING",
			},
		}),
	)

	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	engine.Spawn(NewCoordinator(
		robot1Controller,
		robot2Controller,
		robot3Controller,
		workbench1Controller,
		workbench2Controller), "coordinator")
	<-make(chan struct{})
}
