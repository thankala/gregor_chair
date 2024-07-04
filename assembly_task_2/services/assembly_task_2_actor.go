package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/messages"
)

type AssemblyTask2Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask2Actor(robot controllers.RobotController) *AssemblyTask2Actor {
	return &AssemblyTask2Actor{robot: robot}
}

func (a *AssemblyTask2Actor) Task() enums.AssemblyTask {
	return enums.AssemblyTask2
}

func (a *AssemblyTask2Actor) Steps() interfaces.StepHandlers[AssemblyTask2Actor] {
	return interfaces.StepHandlers[AssemblyTask2Actor]{
		enums.Step1: a.processStep1GetCastors,
		enums.Step2: a.processStep2AttachW1F2,
		enums.Step3: a.processStep3GetLift,
		enums.Step4: a.processStep4AttachW1F2,
	}
}

func (a *AssemblyTask2Actor) processStep1GetCastors(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(msg.Task); err != nil {
		ctx.Send(ctx.PID(), msg)
		return
	}
	a.robot.MoveToStorage(enums.StorageB4)
	a.robot.PickupItemFromStorage(enums.StorageB4)

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.FixtureRequested,
		Task:        a.Task(),
		Step:        enums.Step2,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture2,
		Expected:    []enums.Stage{enums.BaseAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask2Actor) processStep2AttachW1F2(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	a.robot.MoveToWorkbench(enums.Workbench1)
	a.robot.PickAndPlace()
	item := a.robot.PlaceItem()

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.ComponentAttached,
		Task:        a.Task(),
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture2,
		Component:   item,
	})

	ctx.Send(ctx.PID(), &messages.AssemblyTaskMessage{
		Event:       enums.AssemblyTaskEvent,
		Source:      a.Task().String(),
		Destination: enums.AssemblyTask2.String(),
		Task:        a.Task(),
		Step:        enums.Step3,
	})
}

func (a *AssemblyTask2Actor) processStep3GetLift(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	a.robot.MoveToStorage(enums.StorageB5)
	a.robot.PickupItemFromStorage(enums.StorageB5)

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.FixtureRequested,
		Task:        a.Task(),
		Step:        enums.Step4,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture2,
		Expected:    []enums.Stage{enums.CastorsAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask2Actor) processStep4AttachW1F2(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	a.robot.MoveToWorkbench(enums.Workbench1)
	a.robot.PickAndFlipAndPress()
	item := a.robot.PlaceItem()
	a.robot.ClearCurrentTask()

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.ComponentAttached,
		Task:        a.Task(),
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture2,
		Component:   item,
	})

	ctx.Send(ctx.PID(), &messages.AssemblyTaskMessage{
		Event:       enums.AssemblyTaskEvent,
		Source:      a.Task().String(),
		Destination: enums.AssemblyTask3.String(),
		Task:        enums.AssemblyTask3,
		Step:        enums.Step1,
	})
}
