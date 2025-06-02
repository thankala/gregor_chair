package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
)

type AssemblyTask4Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask4Actor(robot controllers.RobotController) *AssemblyTask4Actor {
	return &AssemblyTask4Actor{robot: robot}
}

func (a *AssemblyTask4Actor) Task() enums.Task {
	return enums.AssemblyTask4
}

func (a *AssemblyTask4Actor) Steps() interfaces.StepHandlers[AssemblyTask4Actor] {
	return interfaces.StepHandlers[AssemblyTask4Actor]{
		enums.Step1: a.requestFixtureAtW2F1,
		enums.Step2: a.getSeatAndPlace,
		enums.Step3: a.getSeatPlateAndAttach,
	}
}

func (a *AssemblyTask4Actor) requestFixtureAtW2F1(event *events.AssemblyTaskEvent, ctx *actor.Context) {
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
		Workbench:   enums.Workbench2,
		Fixture:     enums.Fixture1,
		Expected:    []enums.Stage{enums.Initial},
	})
}

func (a *AssemblyTask4Actor) getSeatAndPlace(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToConveyorBelt(enums.ConveyorBelt1)
	a.robot.PickupItemFromConveyorBelt(enums.ConveyorBelt1)
	a.robot.MoveToWorkbench(enums.Workbench2)
	a.robot.Place()
	item := a.robot.ReleaseItem()

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentPlaced,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench2,
		Fixture:     enums.Fixture1,
		Component:   item,
	})

	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: a.Task(),
		Step:        enums.Step3,
	})
}

func (a *AssemblyTask4Actor) getSeatPlateAndAttach(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToStorage(enums.StorageB3)
	a.robot.PickupItemFromStorage(enums.StorageB3)
	a.robot.MoveToWorkbench(enums.Workbench2)
	a.robot.Screw()
	item := a.robot.ReleaseItem()
	a.robot.ClearCurrentTask()

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentAttached,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench2,
		Fixture:     enums.Fixture1,
		Component:   item,
	})

	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: enums.AssemblyTask1,
		Step:        enums.Step1,
	})
}
