package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/messages"
)

type AssemblyTask1Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask1Actor(robot controllers.RobotController) *AssemblyTask1Actor {
	return &AssemblyTask1Actor{robot: robot}
}

func (a *AssemblyTask1Actor) Task() enums.AssemblyTask {
	return enums.AssemblyTask1
}

func (a *AssemblyTask1Actor) Steps() interfaces.StepHandlers[AssemblyTask1Actor] {
	return interfaces.StepHandlers[AssemblyTask1Actor]{
		enums.Step1: a.processStep1GetLegs,
		enums.Step2: a.processStep2PlaceW1F1,
		enums.Step3: a.processStep3GetBase,
		enums.Step4: a.processStep4AttachW1F1,
	}
}

func (a *AssemblyTask1Actor) processStep1GetLegs(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	//todo check here if i need to set source as this task before sending it back to the queue
	if err := a.robot.SetCurrentTask(msg.Task); err != nil {
		ctx.Send(ctx.PID(), msg)
		return
	}
	a.robot.MoveToStorage(enums.StorageB1)
	a.robot.PickupItemFromStorage(enums.StorageB1)

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.RequestFixture,
		Task:        a.Task(),
		Step:        enums.Step2,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture1,
		Expected:    []enums.Stage{enums.Initial},
		IsPickup:    false,
	})
}

func (a *AssemblyTask1Actor) processStep2PlaceW1F1(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	a.robot.MoveToWorkbench(enums.Workbench1)
	item := a.robot.PlaceItem()

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.PlaceComponent,
		Task:        a.Task(),
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture1,
		Component:   item,
	})

	ctx.Send(ctx.PID(), &messages.AssemblyTaskMessage{
		Event:       enums.AssemblyTaskEvent,
		Source:      a.Task().String(),
		Destination: a.Task().String(),
		Task:        a.Task(),
		Step:        enums.Step3,
	})
}

func (a *AssemblyTask1Actor) processStep3GetBase(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	a.robot.MoveToStorage(enums.StorageB2)
	a.robot.PickupItemFromStorage(enums.StorageB2)

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.RequestFixture,
		Task:        a.Task(),
		Step:        enums.Step4,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture1,
		Expected:    []enums.Stage{enums.LegsAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask1Actor) processStep4AttachW1F1(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	a.robot.MoveToWorkbench(enums.Workbench1)
	a.robot.PickAndPlace()
	a.robot.PickAndInsert()
	item := a.robot.PlaceItem()
	a.robot.ClearCurrentTask()

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.AttachComponent,
		Task:        a.Task(),
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture1,
		Component:   item,
	})

	ctx.Send(ctx.PID(), &messages.AssemblyTaskMessage{
		Event:       enums.AssemblyTaskEvent,
		Source:      a.Task().String(),
		Destination: enums.AssemblyTask4.String(),
		Task:        enums.AssemblyTask4,
		Step:        enums.Step1,
	})
}
