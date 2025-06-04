package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/logger"
)

type AssemblyTask2Actor struct {
	robot        controllers.RobotController
	numberOfRuns int
}

func NewAssemblyTask2Actor(robot controllers.RobotController) *AssemblyTask2Actor {
	return &AssemblyTask2Actor{robot: robot, numberOfRuns: 1}
}

func (a *AssemblyTask2Actor) Task() enums.Task {
	return enums.AssemblyTask2
}

func (a *AssemblyTask2Actor) Steps() interfaces.StepHandlers[AssemblyTask2Actor] {
	return interfaces.StepHandlers[AssemblyTask2Actor]{
		enums.Step1: a.requestFixtureAtW1F2,
		enums.Step2: a.getCastorsAndAttach,
	}
}

func (a *AssemblyTask2Actor) requestFixtureAtW1F2(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(event.Destination, a.numberOfRuns); err != nil {
		ctx.Send(ctx.PID(), event)
		return
	}

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Workbench:   enums.Workbench1,
		Type:        enums.FixtureRequested,
		Step:        enums.Step2,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture2,
		Expected:    []enums.Stage{enums.LegsAttached},
	})
	logger.Get().Info("Fixture requested", "Task", a.Task(), "Caller", a.robot.Key(), "Workbench", enums.Workbench1, "Fixture", enums.Fixture2, "Chair", a.numberOfRuns)
}

func (a *AssemblyTask2Actor) getCastorsAndAttach(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToStorage(enums.StorageB4, a.numberOfRuns)
	a.robot.PickupItemFromStorage(enums.StorageB4, a.numberOfRuns)
	a.robot.MoveToWorkbench(enums.Workbench1, a.numberOfRuns)
	a.robot.Press(a.numberOfRuns)
	item := a.robot.ReleaseItem(a.numberOfRuns)

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Workbench:   enums.Workbench1,
		Type:        enums.ComponentAttached,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture2,
		Component:   item,
	})
	logger.Get().Info("Component attached", "Task", a.Task(), "Caller", a.robot.Key(), "Workbench", enums.Workbench1, "Fixture", enums.Fixture2, "Component", item.String(), "Chair", a.numberOfRuns)
	a.robot.ClearCurrentTask(a.numberOfRuns)
	a.numberOfRuns = a.numberOfRuns + 1

	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: enums.AssemblyTask3,
		Step:        enums.Step1,
	})
}
