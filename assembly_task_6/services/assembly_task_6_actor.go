package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/messages"
)

type AssemblyTask6Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask6Actor(robot controllers.RobotController) *AssemblyTask6Actor {
	return &AssemblyTask6Actor{robot: robot}
}

func (a *AssemblyTask6Actor) Task() enums.AssemblyTask {
	return enums.AssemblyTask6
}

func (a *AssemblyTask6Actor) Steps() interfaces.StepHandlers[AssemblyTask6Actor] {
	return interfaces.StepHandlers[AssemblyTask6Actor]{
		enums.Step1: a.processStep1GetBack,
		enums.Step2: a.processStep2AttachW1F3,
	}
}

func (a *AssemblyTask6Actor) processStep1GetBack(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(msg.Task); err != nil {
		ctx.Send(ctx.PID(), msg)
		return
	}
	a.robot.MoveToConveyorBelt(enums.ConveyorBelt2)
	a.robot.PickupItemFromConveyorBelt(enums.ConveyorBelt2)
	a.robot.MoveToWorkbench(enums.Workbench1)

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.RequestFixture,
		Task:        a.Task(),
		Step:        enums.Step2,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture3,
		Expected:    []enums.Stage{enums.ScrewsAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask6Actor) processStep2AttachW1F3(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	a.robot.PickAndInsert()
	a.robot.ScrewPickAndFasten()
	item := a.robot.PlaceItem()
	a.robot.ClearCurrentTask()

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.AttachComponent,
		Task:        a.Task(),
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture3,
		Component:   item,
	})

	ctx.Send(ctx.PID(), &messages.AssemblyTaskMessage{
		Event:       enums.AssemblyTaskEvent,
		Source:      a.Task().String(),
		Destination: enums.AssemblyTask7.String(),
		Task:        enums.AssemblyTask7,
		Step:        enums.Step1,
	})

	ctx.Send(ctx.PID(), &messages.AssemblyTaskMessage{
		Event:       enums.AssemblyTaskEvent,
		Source:      a.Task().String(),
		Destination: enums.AssemblyTask8.String(),
		Task:        enums.AssemblyTask8,
		Step:        enums.Step1,
	})
}
