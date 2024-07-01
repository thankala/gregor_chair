package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/messages"
)

type AssemblyTask5Actor struct {
	robot controllers.RobotController
}

func (a *AssemblyTask5Actor) Task() enums.AssemblyTask {
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

func (a *AssemblyTask5Actor) processStep1RequestAttachW1F2(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(msg.Task); err != nil {
		ctx.Send(ctx.PID(), msg)
		return
	}
	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.RequestFixture,
		Task:        a.Task(),
		Step:        enums.Step2,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture2,
		Expected:    []enums.Stage{enums.SeatAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask5Actor) processStep2AttachW1F2(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	a.robot.ScrewPickAndFasten()
	a.robot.ClearCurrentTask()

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.AttachComponent,
		Task:        a.Task(),
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture2,
		Component:   enums.Screws,
	})
}
