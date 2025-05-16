package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/interfaces"
)

type AssemblyTask8Actor struct {
	robot controllers.RobotController
}

func NewAssemblyTask8Actor(robot controllers.RobotController) *AssemblyTask8Actor {
	return &AssemblyTask8Actor{robot: robot}
}

func (a *AssemblyTask8Actor) Task() enums.Task {
	return enums.AssemblyTask8
}

func (a *AssemblyTask8Actor) Steps() interfaces.StepHandlers[AssemblyTask8Actor] {
	return interfaces.StepHandlers[AssemblyTask8Actor]{
		enums.Step1: a.processStep1GetRightArm,
		enums.Step2: a.processStep2AttachW1F3,
	}
}

func (a *AssemblyTask8Actor) processStep1GetRightArm(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	if err := a.robot.SetCurrentTask(event.Destination); err != nil {
		ctx.Send(ctx.PID(), event)
		return
	}

	a.robot.MoveToStorage(enums.StorageB6R)
	a.robot.PickupItemFromStorage(enums.StorageB6R)
	a.robot.MoveToWorkbench(enums.Workbench1)

	ctx.Send(ctx.PID(), &events.OrchestratorEvent{
		Source:      a.Task(),
		Destination: enums.Orchestrator,
		Type:        enums.FixtureRequested,
		Step:        enums.Step2,
		Caller:      a.robot.Key(),
		Workbench:   enums.Workbench1,
		Fixture:     enums.Fixture3,
		Expected:    []enums.Stage{enums.BackAttached, enums.LeftArmAttached},
		IsPickup:    false,
	})
}

func (a *AssemblyTask8Actor) processStep2AttachW1F3(event *events.AssemblyTaskEvent, ctx *actor.Context) {
	a.robot.ValidateCurrentTask(event.Destination)
	a.robot.PickAndFlipAndPress()
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
