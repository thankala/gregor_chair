package services

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/messages"
	"github.com/thankala/gregor_chair_common/models"
	"slices"
)

type Coordinator2Actor struct {
	workbench controllers.WorkbenchController
}

func NewCoordinator2Actor(workbench controllers.WorkbenchController) *Coordinator2Actor {
	return &Coordinator2Actor{
		workbench: workbench,
	}
}

func (a *Coordinator2Actor) Coordinator() enums.Coordinator {
	return enums.Coordinator2
}

func (a *Coordinator2Actor) RequestFixture(msg *messages.CoordinatorMessage) {
	a.workbench.PushRequest(models.Request{
		Task:     msg.Task,
		Step:     msg.Step,
		Caller:   msg.Caller,
		Expected: msg.Expected,
		IsPickup: msg.IsPickup,
	}, msg.Fixture)
}

func (a *Coordinator2Actor) PlaceComponent(msg *messages.CoordinatorMessage) {
	a.workbench.SetItem(msg.Task, msg.Caller, msg.Fixture, msg.Component)
	a.workbench.SetFixtureOwner(enums.NoneAssemblyTask, msg.Caller, msg.Fixture)
}

func (a *Coordinator2Actor) AttachComponent(msg *messages.CoordinatorMessage) {
	a.workbench.AttachItem(msg.Task, msg.Caller, msg.Fixture, msg.Component)
	a.workbench.SetFixtureOwner(enums.NoneAssemblyTask, msg.Caller, msg.Fixture)
}

func (a *Coordinator2Actor) Process(ctx *actor.Context) {
	fixtures := a.workbench.GetFixturesContent()
	a.workbench.SetLEDs(fixtures)
	for _, fixture := range fixtures {
		request := a.workbench.PeekRequest(fixture.Fixture)
		for request != nil {
			popRequest := a.workbench.PopRequest(fixture.Fixture)
			request = a.workbench.PeekRequest(fixture.Fixture)

			if slices.Contains(popRequest.Expected, fixture.Component.Stage()) {
				if popRequest.Task == enums.AssemblyTask4 && popRequest.Step == enums.Step2 {
					a.workbench.SetLED(enums.Fixture1, "ASSEMBLING")
				}
				if popRequest.Task == enums.AssemblyTask3 && popRequest.Step == enums.Step2 {
					a.workbench.SetLED(enums.Fixture1, "ASSEMBLING")
				}

				a.workbench.SetFixtureOwner(popRequest.Task, popRequest.Caller, fixture.Fixture)
				component := enums.NoneComponent
				if popRequest.IsPickup {
					component = a.workbench.GetItem(popRequest.Task, popRequest.Caller, fixture.Fixture)
					a.workbench.SetFixtureOwner(enums.NoneAssemblyTask, popRequest.Caller, fixture.Fixture)
				}
				ctx.Send(ctx.PID(), &messages.AssemblyTaskMessage{
					Event:       enums.AssemblyTaskEvent,
					Source:      a.Coordinator().String(),
					Destination: popRequest.Task.String(),
					Task:        popRequest.Task,
					Step:        popRequest.Step,
					Component:   component,
				})
			} else {
				a.workbench.PushRequest(*popRequest, fixture.Fixture)
			}
		}
	}
}
