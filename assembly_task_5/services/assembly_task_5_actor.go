package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
)

type AssemblyTask5Actor struct {
	robot controllers.RobotController
}

func (a *AssemblyTask5Actor) Task() enums.Task {
	return enums.AssemblyTask5
}

func (a *AssemblyTask5Actor) Steps() interfaces.StepHandlers[AssemblyTask5Actor] {
	return interfaces.StepHandlers[AssemblyTask5Actor]{
		enums.Step1: a.processStep1RequestAttachW1F2,
		enums.Step2: a.processStep2AttachW1F2,
	}
}

func NewAssemblyTask5Actor(robot controllers.RobotController) *AssemblyTask5Actor {
	return &AssemblyTask5Actor{robot}
}

func (a *AssemblyTask5Actor) processStep1RequestAttachW1F2(event *events.AssemblyTaskEvent, ctx *actor.Context) {
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
		Fixture:     enums.Fixture2,
		Expected:    []enums.Stage{enums.SeatAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask5Actor) processStep2AttachW1F2(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.ScrewPickAndFasten()
	a.robot.ClearCurrentTask()

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentAttached,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture2,
		Component:   enums.Screws,
	})
}
