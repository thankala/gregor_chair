package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/messages"
)

type AssemblyTask3Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask3Actor(robot controllers.RobotController) *AssemblyTask3Actor {
	return &AssemblyTask3Actor{robot: robot}
}

func (a *AssemblyTask3Actor) Task() enums.AssemblyTask {
	return enums.AssemblyTask3
}

func (a *AssemblyTask3Actor) Steps() interfaces.StepHandlers[AssemblyTask3Actor] {
	return interfaces.StepHandlers[AssemblyTask3Actor]{
		enums.Step1: a.processStep1RequestW2F1,
		enums.Step2: a.processStep2AttachW2F1,
		enums.Step3: a.processStep3PickupW2F1,
		enums.Step4: a.processStep4AttachW1F2,
	}
}

func (a *AssemblyTask3Actor) processStep1RequestW2F1(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(msg.Task); err != nil {
		ctx.Send(ctx.PID(), msg)
		return
	}
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
		Expected:    []enums.Stage{enums.SeatPlateAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask3Actor) processStep2AttachW2F1(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	a.robot.PickAndInsert()
	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator2.String(),
		Type:        enums.AttachComponent,
		Task:        a.Task(),
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture1,
		Component:   enums.SeatScrews,
	})

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator2.String(),
		Type:        enums.RequestFixture,
		Task:        a.Task(),
		Step:        enums.Step3,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture1,
		Expected:    []enums.Stage{enums.SeatScrewsAttached},
		IsPickup:    true,
	})
}

func (a *AssemblyTask3Actor) processStep3PickupW2F1(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	a.robot.PickupItemFromWorkbench(msg.Component, enums.Workbench2)
	a.robot.MoveToWorkbench(enums.Workbench1)
	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.RequestFixture,
		Task:        a.Task(),
		Step:        enums.Step4,
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture2,
		Expected:    []enums.Stage{enums.LiftAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask3Actor) processStep4AttachW1F2(msg *messages.AssemblyTaskMessage, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(msg.Task)
	a.robot.ScrewPickAndFasten()
	item := a.robot.ReleaseItem()
	a.robot.ClearCurrentTask()

	ctx.Send(ctx.PID(), &messages.CoordinatorMessage{
		Event:       enums.CoordinatorEvent,
		Source:      a.Task().String(),
		Destination: enums.Coordinator1.String(),
		Type:        enums.AttachComponent,
		Task:        a.Task(),
		Caller:      a.robot.Key(),
		Fixture:     enums.Fixture2,
		Component:   item,
	})
	ctx.Send(ctx.PID(), &messages.AssemblyTaskMessage{
		Event:       enums.AssemblyTaskEvent,
		Source:      a.Task().String(),
		Destination: enums.AssemblyTask5.String(),
		Task:        enums.AssemblyTask5,
		Step:        enums.Step1,
	})
}
