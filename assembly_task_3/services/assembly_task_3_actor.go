package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
)

type AssemblyTask3Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask3Actor(robot controllers.RobotController) *AssemblyTask3Actor {
	return &AssemblyTask3Actor{robot: robot}
}

func (a *AssemblyTask3Actor) Task() enums.Task {
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

func (a *AssemblyTask3Actor) processStep1RequestW2F1(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(event.Destination); err != nil {
		ctx.Send(ctx.PID(), event)
		return
	}
	a.robot.MoveToWorkbench(enums.Workbench2)
	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.FixtureRequested,
		Step:        enums.Step2,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench2,
		Fixture:     enums.Fixture1,
		Expected:    []enums.Stage{enums.SeatPlateAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask3Actor) processStep2AttachW2F1(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.PickAndInsert()
	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentAttached,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench2,
		Fixture:     enums.Fixture1,
		Component:   enums.SeatScrews,
	})

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.FixtureRequested,
		Step:        enums.Step3,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench2,
		Fixture:     enums.Fixture1,
		Expected:    []enums.Stage{enums.SeatScrewsAttached},
		IsPickup:    true,
	})
}

func (a *AssemblyTask3Actor) processStep3PickupW2F1(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.PickupItemFromWorkbench(event.Component, enums.Workbench2)
	a.robot.MoveToWorkbench(enums.Workbench1)
	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.FixtureRequested,
		Step:        enums.Step4,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture2,
		Expected:    []enums.Stage{enums.LiftAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask3Actor) processStep4AttachW1F2(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.ScrewPickAndFasten()
	item := a.robot.PlaceItem()
	a.robot.ClearCurrentTask()

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentAttached,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture2,
		Component:   item,
	})

	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: enums.AssemblyTask5,
		Step:        enums.Step1,
	})
}
