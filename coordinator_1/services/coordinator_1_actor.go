package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/logger"
	"github.com/thankala/gregor_chair_common/messages"
	"github.com/thankala/gregor_chair_common/models"
	"slices"
	"time"
)

type Coordinator1Actor struct {
	workbench controllers.WorkbenchController
}

func NewCoordinator1Actor(workbench controllers.WorkbenchController) *Coordinator1Actor {
	return &Coordinator1Actor{
		workbench: workbench,
	}
}

func (a *Coordinator1Actor) Coordinator() enums.Coordinator {
	return enums.Coordinator1
}

func (a *Coordinator1Actor) RequestFixture(msg *messages.CoordinatorMessage) {
	a.workbench.PushRequest(models.Request{
		Task:     msg.Task,
		Step:     msg.Step,
		Caller:   msg.Caller,
		Expected: msg.Expected,
		IsPickup: msg.IsPickup,
	}, msg.Fixture)
}

func (a *Coordinator1Actor) PlaceComponent(msg *messages.CoordinatorMessage) {
	a.workbench.SetItem(msg.Task, msg.Caller, msg.Fixture, msg.Component)
	a.workbench.SetFixtureOwner(enums.NoneAssemblyTask, msg.Caller, msg.Fixture)
}

func (a *Coordinator1Actor) AttachComponent(msg *messages.CoordinatorMessage) {
	a.workbench.AttachItem(msg.Task, msg.Caller, msg.Fixture, msg.Component)
	a.workbench.SetFixtureOwner(enums.NoneAssemblyTask, msg.Caller, msg.Fixture)
}

func (a *Coordinator1Actor) Process(ctx *actor.Context) {
	fixtures := a.workbench.GetFixturesContent()
	for _, fixture := range fixtures {
		request := a.workbench.PopRequest(fixture.Fixture)
		if request == nil {
			continue
		}

		a.workbench.SetLEDs(fixtures)
		if slices.Contains(request.Expected, fixture.Component.Stage()) {
			a.workbench.SetFixtureOwner(request.Task, request.Caller, fixture.Fixture)
			component := enums.NoneComponent
			if request.IsPickup {
				component = a.workbench.GetItem(request.Task, request.Caller, fixture.Fixture)
				a.workbench.SetFixtureOwner(enums.NoneAssemblyTask, request.Caller, fixture.Fixture)
			}

			ctx.Send(ctx.PID(), &messages.AssemblyTaskMessage{
				Event:       enums.AssemblyTaskEvent,
				Source:      a.Coordinator().String(),
				Destination: request.Task.String(),
				Task:        request.Task,
				Step:        request.Step,
				Component:   component},
			)
		} else {
			a.workbench.PushRequest(*request, fixture.Fixture)
		}
	}

	if (len(fixtures) == 3) && (fixtures[2].Component.Stage() == enums.Completed) {
		logger.Get().Info("One chair has been assembled %v\n", time.Now())
		a.workbench.RemoveCompletedItem()
	}

	if a.workbench.CanRotate() {
		fixtures := a.workbench.RotateFixtures()
		//LEDs
		a.workbench.SetLEDs(fixtures)

		if fixtures[1].Component.Stage() == enums.BaseAttached {
			a.workbench.SetLED(fixtures[1].Fixture, "ASSEMBLING")
			ctx.Send(ctx.PID(), &messages.AssemblyTaskMessage{
				Event:       enums.AssemblyTaskEvent,
				Source:      a.Coordinator().String(),
				Destination: enums.AssemblyTask2.String(),
				Task:        enums.AssemblyTask2,
				Step:        enums.Step1,
			})
		}

		if fixtures[2].Component.Stage() == enums.ScrewsAttached {
			a.workbench.SetLED(fixtures[2].Fixture, "ASSEMBLING")
			ctx.Send(ctx.PID(), &messages.AssemblyTaskMessage{
				Event:       enums.AssemblyTaskEvent,
				Source:      a.Coordinator().String(),
				Destination: enums.AssemblyTask6.String(),
				Task:        enums.AssemblyTask6,
				Step:        enums.Step1,
			})
		}
	}

}