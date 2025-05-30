package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
)

type AssemblyTask5Actor struct {
	robot controllers.RobotController
}

func (a *AssemblyTask5Actor) Task() enums.Task {
	return enums.AssemblyTask5
}

func (a *AssemblyTask5Actor) Steps() interfaces.StepHandlers[AssemblyTask5Actor] {
	return interfaces.StepHandlers[AssemblyTask5Actor]{
		enums.Step1: a.requestFixtureAtW2F1,
		enums.Step2: a.moveToW2AndRequestPickup,
		enums.Step3: a.pickup,
		enums.Step4: a.requestFixtureAtW1F2,
		enums.Step5: a.attachComposite,
	}
}

func NewAssemblyTask5Actor(robot controllers.RobotController) *AssemblyTask5Actor {
	return &AssemblyTask5Actor{robot}
}

func (a *AssemblyTask5Actor) requestFixtureAtW2F1(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(event.Destination); err != nil {
		ctx.Send(ctx.PID(), event)
		return
	}
	// a.robot.MoveToWorkbench(enums.Workbench2)
	// logger.Get().Info("Fixture requested: ", enums.Fi "Workbench", "Robot", a.robot.Key(), "Task", a.Task())
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

func (a *AssemblyTask5Actor) moveToW2AndRequestPickup(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToWorkbench(enums.Workbench2)
	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.FixtureRequested,
		Step:        enums.Step3,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench2,
		Fixture:     enums.Fixture1,
		Expected:    []enums.Stage{enums.SeatPlateAttached},
		IsPickup:    true,
	})
}

func (a *AssemblyTask5Actor) pickup(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.PickupItemFromWorkbench(event.Component, enums.Workbench2)
	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: enums.AssemblyTask5,
		Step:        enums.Step4,
	})
}

func (a *AssemblyTask5Actor) requestFixtureAtW1F2(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.FixtureRequested,
		Step:        enums.Step5,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture2,
		Expected:    []enums.Stage{enums.LiftAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask5Actor) attachComposite(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToWorkbench(enums.Workbench1)
	a.robot.Press()
	item := a.robot.ReleaseItem()
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
		Destination: enums.AssemblyTask6,
		Step:        enums.Step1,
	})
}

// ctx.Send(ctx.PID(), &events.OrchestratorEvent{
// 		Source:      a.Task(),
// 		Destination: enums.Orchestrator,
// 		Type:        enums.FixtureRequested,
// 		Step:        enums.Step3,
// 		Caller:      a.robot.Key(),
// 		Workbench:   enums.Workbench1,
// 		Fixture:     enums.Fixture2,
// 		Expected:    []enums.Stage{enums.LiftAttached},
// 		IsPickup:    false,
// 	})
