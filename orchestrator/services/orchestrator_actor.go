package services

import (
	"fmt"
	"slices"

	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/logger"
	"github.com/thankala/gregor_chair_common/models"
)

type OrchestratorActor struct {
	workbenchControllers map[enums.Workbench]*controllers.WorkbenchController
	robotControllers     map[enums.Robot]*controllers.RobotController
	numberOfChairs       int
}

func NewOrchestratorActor(workbenchControllers []controllers.WorkbenchController, robotControllers []controllers.RobotController) *OrchestratorActor {
	workbenches := make(map[enums.Workbench]*controllers.WorkbenchController)
	robots := make(map[enums.Robot]*controllers.RobotController)
	for _, workbenchController := range workbenchControllers {
		workbenchController.ResetState()
		workbenchController.ResetLEDs()
		workbenches[workbenchController.Key()] = &workbenchController
	}
	for _, robotController := range robotControllers {
		robotController.ResetState()
		robotController.ResetRobot()
		robots[robotController.Key()] = &robotController
	}
	return &OrchestratorActor{
		workbenchControllers: workbenches,
		robotControllers:     robots,
		numberOfChairs:       1,
	}
}

func (o *OrchestratorActor) Orchestrator() enums.Task {
	return enums.Orchestrator
}

func (o *OrchestratorActor) Process(ctx *actor.Context, event *events.OrchestratorEvent) {
	workbench := o.workbenchControllers[event.Workbench]
	fixtures := workbench.GetFixturesContent()
	fixture := workbench.FindFixtureContentByFixture(fixtures, event.Fixture)

	switch event.Type {
	case enums.ComponentPlaced:
		o.handleComponentPlaced(workbench, event)
	case enums.ComponentAttached:
		o.handleComponentAttached(workbench, event)
	case enums.ComponentPickedUp:
		o.handleComponentPickedUp(workbench, event)
	case enums.FixtureRequested:
		o.handleFixtureRequested(ctx, workbench, event, fixture)
	}

	o.handleChairCompleted(workbench)
	o.handleRotation(ctx, workbench)
}

func (o *OrchestratorActor) StartAssembly(ctx *actor.Context, event *events.OrchestratorEvent) {
	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      o.Orchestrator(),
		Destination: enums.AssemblyTask1,
		Step:        enums.Step1,
	})
}

func (o *OrchestratorActor) handleComponentPlaced(workbench *controllers.WorkbenchController, event *events.OrchestratorEvent) {
	workbench.SetItem(event.Source, event.Caller, event.Fixture, event.Component)
	workbench.SetFixtureOwner(event.Source, event.Caller, event.Fixture)
}

func (o *OrchestratorActor) handleComponentAttached(workbench *controllers.WorkbenchController, event *events.OrchestratorEvent) {
	workbench.AttachItem(event.Source, event.Caller, event.Fixture, event.Component)
	workbench.SetFixtureOwner(enums.NoneTask, event.Caller, event.Fixture)
}

func (o *OrchestratorActor) handleComponentPickedUp(workbench *controllers.WorkbenchController, event *events.OrchestratorEvent) {
	workbench.ReleaseItem(event.Source, event.Caller, event.Fixture)
	workbench.SetFixtureOwner(enums.NoneTask, event.Caller, event.Fixture)
}

func (o *OrchestratorActor) handleFixtureRequested(ctx *actor.Context, workbench *controllers.WorkbenchController, event *events.OrchestratorEvent, fixture models.FixtureContent) {
	if !slices.Contains(event.Expected, fixture.Component.Stage()) {
		ctx.Send(ctx.PID(), event)
		return
	}
	logger.Get().Info("Fixture request granted", "Task", event.Source, "Caller", event.Caller, "Workbench", event.Workbench, "Fixture", event.Fixture, "Chair", o.numberOfChairs)
	component := workbench.SetFixtureOwner(event.Source, event.Caller, fixture.Fixture)
	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      o.Orchestrator(),
		Destination: event.Source,
		Step:        event.Step,
		Component:   component,
	})
}

func (o *OrchestratorActor) handleRotation(ctx *actor.Context, workbench *controllers.WorkbenchController) {
	if workbench.CanRotate() {
		fixtures := workbench.RotateFixtures()
		workbench.SetLED(enums.Fixture1, "FREE")
		if fixtures[1].Component.Stage() == enums.LegsAttached {
			ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
				Source:      o.Orchestrator(),
				Destination: enums.AssemblyTask2,
				Step:        enums.Step1,
			})
		}
	}
}

func (a *OrchestratorActor) handleChairCompleted(workbench *controllers.WorkbenchController) {
	if workbench.HasAssembledFinished() {
		logger.Get().Info("Chair assembled", "Number of chair:", fmt.Sprint(a.numberOfChairs))
		a.numberOfChairs++
		workbench.RemoveCompletedItem()
	}
}
