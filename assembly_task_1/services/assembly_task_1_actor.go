package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/logger"
)

type AssemblyTask1Actor struct {
	robot        controllers.RobotController
	numberOfRuns int
}

func NewAssemblyTask1Actor(robot controllers.RobotController) *AssemblyTask1Actor {
	return &AssemblyTask1Actor{robot: robot, numberOfRuns: 1}
}

func (a *AssemblyTask1Actor) Task() enums.Task {
	return enums.AssemblyTask1
}

func (a *AssemblyTask1Actor) Steps() interfaces.StepHandlers[AssemblyTask1Actor] {
	return interfaces.StepHandlers[AssemblyTask1Actor]{
		enums.Step1: a.requestFixtureAtW1F1,
		enums.Step2: a.getBaseAndPlace,
		enums.Step3: a.getLegsAndAttach,
	}
}

func (a *AssemblyTask1Actor) requestFixtureAtW1F1(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(event.Destination, a.numberOfRuns); err != nil {
		ctx.Send(ctx.PID(), event)
		return
	}

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.FixtureRequested,
		Step:        enums.Step2,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture1,
		Expected:    []enums.Stage{enums.Initial},
	})
	logger.Get().Info("Fixture requested", "Task", a.Task(), "Caller", a.robot.Key(), "Workbench", enums.Workbench1, "Fixture", enums.Fixture1, "Chair", a.numberOfRuns)
}

func (a *AssemblyTask1Actor) getBaseAndPlace(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToStorage(enums.StorageB2, a.numberOfRuns)
	a.robot.PickupItemFromStorage(enums.StorageB2, a.numberOfRuns)
	a.robot.MoveToWorkbench(enums.Workbench1, a.numberOfRuns)
	a.robot.Place(a.numberOfRuns)
	item := a.robot.ReleaseItem(a.numberOfRuns)

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentPlaced,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture1,
		Component:   item,
	})
	logger.Get().Info("Component placed", "Task", a.Task(), "Caller", a.robot.Key(), "Workbench", enums.Workbench1, "Fixture", enums.Fixture1, "Component", item.String(), "Chair", a.numberOfRuns)

	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: a.Task(),
		Step:        enums.Step3,
	})
}

func (a *AssemblyTask1Actor) getLegsAndAttach(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToStorage(enums.StorageB1, a.numberOfRuns)
	a.robot.PickupItemFromStorage(enums.StorageB1, a.numberOfRuns)
	a.robot.MoveToWorkbench(enums.Workbench1, a.numberOfRuns)
	a.robot.Screw(a.numberOfRuns)
	item := a.robot.ReleaseItem(a.numberOfRuns)

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentAttached,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture1,
		Component:   item,
	})

	logger.Get().Info("Component attached", "Task", a.Task(), "Caller", a.robot.Key(), "Workbench", enums.Workbench1, "Fixture", enums.Fixture1, "Component", item.String(), "Chair", a.numberOfRuns)
	a.robot.ClearCurrentTask(a.numberOfRuns)
	a.numberOfRuns = a.numberOfRuns + 1
	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: enums.AssemblyTask4,
		Step:        enums.Step1,
	})
}
