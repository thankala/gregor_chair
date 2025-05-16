package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
)

type AssemblyTask1Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask1Actor(robot controllers.RobotController) *AssemblyTask1Actor {
	return &AssemblyTask1Actor{robot: robot}
}

func (a *AssemblyTask1Actor) Task() enums.Task {
	return enums.AssemblyTask1
}

func (a *AssemblyTask1Actor) Steps() interfaces.StepHandlers[AssemblyTask1Actor] {
	return interfaces.StepHandlers[AssemblyTask1Actor]{
		enums.Step1: a.processStep1RequestFixture,
		enums.Step2: a.processStep2GetLegsAndPlaceW1F1,
		enums.Step3: a.processStep3GetBase,
		enums.Step4: a.processStep4AttachW1F1,
	}
}

func (a *AssemblyTask1Actor) processStep1RequestFixture(event *events.AssemblyTaskEvent, ctx *actor.Context) {
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
		Fixture:     enums.Fixture1,
		Expected:    []enums.Stage{enums.Initial},
		IsPickup:    false,
	})
}

func (a *AssemblyTask1Actor) processStep2GetLegsAndPlaceW1F1(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToStorage(enums.StorageB1)
	a.robot.PickupItemFromStorage(enums.StorageB1)
	a.robot.MoveToWorkbench(enums.Workbench1)
	item := a.robot.PlaceItem()

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentPlaced,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture1,
		Component:   item,
	})

	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: a.Task(),
		Step:        enums.Step3,
	})
}

func (a *AssemblyTask1Actor) processStep3GetBase(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToStorage(enums.StorageB2)
	a.robot.PickupItemFromStorage(enums.StorageB2)

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.FixtureRequested,
		Step:        enums.Step4,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture1,
		Expected:    []enums.Stage{enums.LegsAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask1Actor) processStep4AttachW1F1(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToWorkbench(enums.Workbench1)
	//a.robot.PickAndPlace()
	a.robot.PickAndInsert()
	item := a.robot.PlaceItem()
	a.robot.ClearCurrentTask()

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentAttached,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture1,
		Component:   item,
	})

	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: enums.AssemblyTask4,
		Step:        enums.Step1,
	})
}
