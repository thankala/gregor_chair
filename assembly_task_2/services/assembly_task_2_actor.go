package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
)

type AssemblyTask2Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask2Actor(robot controllers.RobotController) *AssemblyTask2Actor {
	return &AssemblyTask2Actor{robot: robot}
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
	if err := a.robot.SetCurrentTask(event.Destination); err != nil {
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
}

func (a *AssemblyTask2Actor) getCastorsAndAttach(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToStorage(enums.StorageB4)
	a.robot.PickupItemFromStorage(enums.StorageB4)
	a.robot.MoveToWorkbench(enums.Workbench1)
	a.robot.Press()
	item := a.robot.ReleaseItem()
	a.robot.ClearCurrentTask()

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Workbench:   enums.Workbench1,
		Type:        enums.ComponentAttached,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture2,
		Component:   item,
	})

	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: enums.AssemblyTask3,
		Step:        enums.Step1,
	})
}
