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
		enums.Step1: a.requestFixtureAtW1F1,
		enums.Step2: a.getBaseAndPlace,
		enums.Step3: a.getLegsAndAttach,
	}
}

func (a *AssemblyTask1Actor) requestFixtureAtW1F1(event *events.AssemblyTaskEvent, ctx *actor.Context) {
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
	})
}

func (a *AssemblyTask1Actor) getBaseAndPlace(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToStorage(enums.StorageB2)
	a.robot.PickupItemFromStorage(enums.StorageB2)
	a.robot.MoveToWorkbench(enums.Workbench1)
	a.robot.Place()
	item := a.robot.ReleaseItem()

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

func (a *AssemblyTask1Actor) getLegsAndAttach(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToStorage(enums.StorageB1)
	a.robot.PickupItemFromStorage(enums.StorageB1)
	a.robot.MoveToWorkbench(enums.Workbench1)
	a.robot.Screw()
	item := a.robot.ReleaseItem()
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
