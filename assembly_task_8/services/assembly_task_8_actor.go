package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/logger"
)

type AssemblyTask8Actor struct {
	robot        controllers.RobotController
	numberOfRuns int
}

func NewAssemblyTask8Actor(robot controllers.RobotController) *AssemblyTask8Actor {
	return &AssemblyTask8Actor{robot: robot, numberOfRuns: 1}
}

func (a *AssemblyTask8Actor) Task() enums.Task {
	return enums.AssemblyTask8
}

func (a *AssemblyTask8Actor) Steps() interfaces.StepHandlers[AssemblyTask8Actor] {
	return interfaces.StepHandlers[AssemblyTask8Actor]{
		enums.Step1: a.requestFixtureAtW1F3,
		enums.Step2: a.getRightArmAndAttach,
	}
}

func (a *AssemblyTask8Actor) requestFixtureAtW1F3(event *events.AssemblyTaskEvent, ctx *actor.Context) {
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
		Fixture:     enums.Fixture3,
		Expected:    []enums.Stage{enums.BackAttached, enums.LeftArmAttached},
	})
	logger.Get().Info("Fixture requested", "Task", a.Task(), "Caller", a.robot.Key(), "Workbench", enums.Workbench1, "Fixture", enums.Fixture3, "Chair", a.numberOfRuns)
}

func (a *AssemblyTask8Actor) getRightArmAndAttach(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToStorage(enums.StorageB6R, a.numberOfRuns)
	a.robot.PickupItemFromStorage(enums.StorageB6R, a.numberOfRuns)
	a.robot.MoveToWorkbench(enums.Workbench1, a.numberOfRuns)
	a.robot.Screw(a.numberOfRuns)
	item := a.robot.ReleaseItem(a.numberOfRuns)

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentAttached,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture3,
		Component:   item,
	})
	logger.Get().Info("Component picked up", "Task", a.Task(), "Caller", a.robot.Key(), "Workbench", enums.Workbench1, "Fixture", enums.Fixture3, "Component", item.String(), "Chair", a.numberOfRuns)
	a.numberOfRuns = a.numberOfRuns + 1
	a.robot.ClearCurrentTask(a.numberOfRuns)
}
