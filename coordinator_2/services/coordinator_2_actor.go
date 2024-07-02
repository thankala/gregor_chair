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
		request := a.workbench.PopRequest(fixture.Fixture)
		for request != nil {
			leftover := a.processRequest(ctx, request, fixture)
			request = nil
			if leftover != nil {
				request = a.workbench.PopRequest(fixture.Fixture)
				a.workbench.PushRequest(*leftover, fixture.Fixture)
			}
		}
	}
}

func (a *Coordinator2Actor) processRequest(ctx *actor.Context, request *models.Request, fixture models.FixtureContent) *models.Request {
	if slices.Contains(request.Expected, fixture.Component.Stage()) {
		if request.Task == enums.AssemblyTask4 && request.Step == enums.Step2 {
			a.workbench.SetLED(enums.Fixture1, "ASSEMBLING")
		}
		if request.Task == enums.AssemblyTask3 && request.Step == enums.Step2 {
			a.workbench.SetLED(enums.Fixture1, "ASSEMBLING")
		}

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
			Component:   component,
		})
		return nil
	} else {
		return request
	}
}
