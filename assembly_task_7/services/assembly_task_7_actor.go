package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/messages"
)

type AssemblyTask7Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask7Actor(robot controllers.RobotController) *AssemblyTask7Actor {
	return &AssemblyTask7Actor{robot: robot}
}

func (a *AssemblyTask7Actor) Task() enums.AssemblyTask {
	return enums.AssemblyTask7
}

func (a *AssemblyTask7Actor) Steps() interfaces.StepHandlers[AssemblyTask7Actor] {
	return interfaces.StepHandlers[AssemblyTask7Actor]{
		enums.Step1: a.processStep1GetLeftArm,
		enums.Step2: a.processStep2AttachW1F3,
	}
}

func (a *AssemblyTask7Actor) processStep1GetLeftArm(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(msg.Task); err != nil {
		ctx.Send(ctx.PID(), msg)
		return
	}

	a.robot.MoveToStorage(enums.StorageB6L)
	a.robot.PickupItemFromStorage(enums.StorageB6L)
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
		Expected:    []enums.Stage{enums.BackAttached, enums.RightArmAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask7Actor) processStep2AttachW1F3(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	a.robot.PickAndInsert()
	a.robot.ScrewPickAndFasten()
	item := a.robot.ReleaseItem()

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

	a.robot.ClearCurrentTask()
}
