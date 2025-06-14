package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/logger"
)

type AssemblyTask5Actor struct {
	robot        controllers.RobotController
	numberOfRuns int
}

func (a *AssemblyTask5Actor) Task() enums.Task {
	return enums.AssemblyTask5
}

func (a *AssemblyTask5Actor) Steps() interfaces.StepHandlers[AssemblyTask5Actor] {
	return interfaces.StepHandlers[AssemblyTask5Actor]{
		enums.Step1: a.requestFixtureAtW2F1,
		enums.Step2: a.moveToW2AndPickup,
		enums.Step3: a.requestFixtureAtW1F2,
		enums.Step4: a.attachComposite,
	}
}

func NewAssemblyTask5Actor(robot controllers.RobotController) *AssemblyTask5Actor {
	return &AssemblyTask5Actor{robot: robot, numberOfRuns: 1}
}

func (a *AssemblyTask5Actor) requestFixtureAtW2F1(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(event.Destination, a.numberOfRuns); err != nil {
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
		Expected:    []enums.Stage{enums.SeatPlateAttached},
	})
	logger.Get().Info("Fixture requested", "Task", a.Task(), "Caller", a.robot.Key(), "Workbench", enums.Workbench2, "Fixture", enums.Fixture1, "Chair", a.numberOfRuns)
}

func (a *AssemblyTask5Actor) moveToW2AndPickup(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToWorkbench(enums.Workbench2, a.numberOfRuns)
	a.robot.PickupItemFromWorkbench(event.Component, enums.Workbench2, a.numberOfRuns)
	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentPickedUp,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench2,
		Fixture:     enums.Fixture1,
		Expected:    []enums.Stage{enums.SeatPlateAttached},
	})
	logger.Get().Info("Component picked up", "Task", a.Task(), "Caller", a.robot.Key(), "Workbench", enums.Workbench2, "Fixture", enums.Fixture1, "Component", event.Component.String(), "Chair", a.numberOfRuns)
	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: enums.AssemblyTask5,
		Step:        enums.Step3,
	})
}

func (a *AssemblyTask5Actor) requestFixtureAtW1F2(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.FixtureRequested,
		Step:        enums.Step4,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture2,
		Expected:    []enums.Stage{enums.LiftAttached},
	})
	logger.Get().Info("Fixture requested", "Task", a.Task(), "Caller", a.robot.Key(), "Workbench", enums.Workbench1, "Fixture", enums.Fixture2, "Chair", a.numberOfRuns)
}

func (a *AssemblyTask5Actor) attachComposite(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.MoveToWorkbench(enums.Workbench1, a.numberOfRuns)
	a.robot.Press(a.numberOfRuns)
	item := a.robot.ReleaseItem(a.numberOfRuns)

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentAttached,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture2,
		Component:   item,
	})
	logger.Get().Info("Component attached", "Task", a.Task(), "Caller", a.robot.Key(), "Workbench", enums.Workbench1, "Fixture", enums.Fixture2, "Component", item.String(), "Chair", a.numberOfRuns)
	a.robot.ClearCurrentTask(a.numberOfRuns)
	a.numberOfRuns = a.numberOfRuns + 1

	ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
		Source:      a.Task(),
		Destination: enums.AssemblyTask6,
		Step:        enums.Step1,
	})

}
