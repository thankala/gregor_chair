package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/messages"
)

type AssemblyTask8Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask8Actor(robot controllers.RobotController) *AssemblyTask8Actor {
	return &AssemblyTask8Actor{robot: robot}
}

func (a *AssemblyTask8Actor) Task() enums.AssemblyTask {
	return enums.AssemblyTask8
}

func (a *AssemblyTask8Actor) Steps() interfaces.StepHandlers[AssemblyTask8Actor] {
	return interfaces.StepHandlers[AssemblyTask8Actor]{
		enums.Step1: a.processStep1GetRightArm,
		enums.Step2: a.processStep2AttachW1F3,
	}
}

func (a *AssemblyTask8Actor) processStep1GetRightArm(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(msg.Task); err != nil {
		ctx.Send(ctx.PID(), msg)
		return
	}

	a.robot.MoveToStorage(enums.StorageB6R)
	a.robot.PickupItemFromStorage(enums.StorageB6R)
	a.robot.MoveToWorkbench(enums.Workbench1)

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.FixtureRequested,
		Task:        msg.Task,
		Step:        enums.Step2,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture3,
		Expected:    []enums.Stage{enums.BackAttached, enums.LeftArmAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask8Actor) processStep2AttachW1F3(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	a.robot.PickAndFlipAndPress()
	item := a.robot.PlaceItem()

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.ComponentAttached,
		Task:        a.Task(),
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture3,
		Component:   item,
	})

	a.robot.ClearCurrentTask()
}
