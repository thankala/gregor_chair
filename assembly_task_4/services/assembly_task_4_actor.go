package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/messages"
)

type AssemblyTask4Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask4Actor(robot controllers.RobotController) *AssemblyTask4Actor {
	return &AssemblyTask4Actor{robot: robot}
}

func (a *AssemblyTask4Actor) Task() enums.AssemblyTask {
	return enums.AssemblyTask4
}

func (a *AssemblyTask4Actor) Steps() interfaces.StepHandlers[AssemblyTask4Actor] {
	return interfaces.StepHandlers[AssemblyTask4Actor]{
		enums.Step1: a.processStep1GetSeat,
		enums.Step2: a.processStep2PlaceW2F1,
		enums.Step3: a.processStep3GetSeatPlate,
		enums.Step4: a.processStep4AttachW2F1,
	}
}

func (a *AssemblyTask4Actor) processStep1GetSeat(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(msg.Task); err != nil {
		ctx.Send(ctx.PID(), msg)
		return
	}
	a.robot.MoveToConveyorBelt(enums.ConveyorBelt2)
	a.robot.PickupItemFromConveyorBelt(enums.ConveyorBelt2)
	a.robot.MoveToWorkbench(enums.Workbench2)

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator2.String(),
		Type:        enums.RequestFixture,
		Task:        a.Task(),
		Step:        enums.Step2,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture1,
		Expected:    []enums.Stage{enums.Initial},
		IsPickup:    false,
	})
}

func (a *AssemblyTask4Actor) processStep2PlaceW2F1(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	item := a.robot.ReleaseItem()
	a.robot.PickAndPlace()
	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator2.String(),
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

func (a *AssemblyTask4Actor) processStep3GetSeatPlate(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.MoveToStorage(enums.StorageB3)
	a.robot.PickupItemFromStorage(enums.StorageB3)
	a.robot.MoveToWorkbench(enums.Workbench2)

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator2.String(),
		Type:        enums.RequestFixture,
		Task:        a.Task(),
		Step:        enums.Step4,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture1,
		Expected:    []enums.Stage{enums.InitialSeat},
		IsPickup:    false,
	})
}

func (a *AssemblyTask4Actor) processStep4AttachW2F1(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	item := a.robot.ReleaseItem()
	a.robot.ClearCurrentTask()
	a.robot.ScrewPickAndFasten()
	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator2.String(),
		Type:        enums.AttachComponent,
		Task:        a.Task(),
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture1,
		Component:   item,
	})

	ctx.Send(ctx.PID(), &messages.AssemblyTaskMessage{
		Event:       enums.AssemblyTaskEvent,
		Source:      a.Task().String(),
		Destination: enums.AssemblyTask1.String(),
		Task:        enums.AssemblyTask1,
		Step:        enums.Step1,
	})
}
