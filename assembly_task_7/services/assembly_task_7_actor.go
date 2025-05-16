package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
)

type AssemblyTask7Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask7Actor(robot controllers.RobotController) *AssemblyTask7Actor {
	return &AssemblyTask7Actor{robot: robot}
}

func (a *AssemblyTask7Actor) Task() enums.Task {
	return enums.AssemblyTask7
}

func (a *AssemblyTask7Actor) Steps() interfaces.StepHandlers[AssemblyTask7Actor] {
	return interfaces.StepHandlers[AssemblyTask7Actor]{
		enums.Step1: a.processStep1GetLeftArm,
		enums.Step2: a.processStep2AttachW1F3,
	}
}

func (a *AssemblyTask7Actor) processStep1GetLeftArm(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(event.Destination); err != nil {
		ctx.Send(ctx.PID(), event)
		return
	}

	a.robot.MoveToStorage(enums.StorageB6L)
	a.robot.PickupItemFromStorage(enums.StorageB6L)
	a.robot.MoveToWorkbench(enums.Workbench1)

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.FixtureRequested,
		Step:        enums.Step2,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture3,
		Expected:    []enums.Stage{enums.BackAttached, enums.RightArmAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask7Actor) processStep2AttachW1F3(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.PickAndInsert()
	a.robot.ScrewPickAndFasten()
	item := a.robot.PlaceItem()

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.ComponentAttached,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture3,
		Component:   item,
	})

	a.robot.ClearCurrentTask()
}
