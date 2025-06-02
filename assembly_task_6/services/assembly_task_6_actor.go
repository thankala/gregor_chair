package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
)

type AssemblyTask6Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask6Actor(robot controllers.RobotController) *AssemblyTask6Actor {
	return &AssemblyTask6Actor{robot: robot}
}

func (a *AssemblyTask6Actor) Task() enums.Task {
	return enums.AssemblyTask6
}

func (a *AssemblyTask6Actor) Steps() interfaces.StepHandlers[AssemblyTask6Actor] {
	return interfaces.StepHandlers[AssemblyTask6Actor]{
		enums.Step1: a.requestFixtureW1F3,
		enums.Step2: a.getBackAndAttach,
	}
}

func (a *AssemblyTask6Actor) requestFixtureW1F3(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(event.Destination); err != nil {
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
		Expected:    []enums.Stage{enums.SeatAttached},
	})
}

func (a *AssemblyTask6Actor) getBackAndAttach(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToConveyorBelt(enums.ConveyorBelt2)
	a.robot.PickupItemFromConveyorBelt(enums.ConveyorBelt2)
	a.robot.MoveToWorkbench(enums.Workbench1)
	a.robot.Screw()
	item := a.robot.ReleaseItem()
	a.robot.ClearCurrentTask()

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentAttached,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture3,
		Component:   item,
	})

	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: enums.AssemblyTask7,
		Step:        enums.Step1,
	})

	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: enums.AssemblyTask8,
		Step:        enums.Step1,
	})
}
