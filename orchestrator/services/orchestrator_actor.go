package services

import (
	"slices"

	"fmt"

	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/logger"
	"github.com/thankala/gregor_chair_common/models"
)

type OrchestratorActor struct {
	workbenches    map[enums.Workbench]*controllers.WorkbenchController
	numberOfChairs int
}

func NewOrchestratorActor(workbenchControllers ...controllers.WorkbenchController) *OrchestratorActor {
	workbenches := make(map[enums.Workbench]*controllers.WorkbenchController)
	for _, workbench := range workbenchControllers {
		workbenches[workbench.Key()] = &workbench
	}
	return &OrchestratorActor{
		workbenches:    workbenches,
		numberOfChairs: 0,
	}
}

func (a *OrchestratorActor) Orchestrator() enums.Task {
	return enums.Orchestrator
}

func (a *OrchestratorActor) FixtureRequested(event *events.OrchestratorEvent) {
	a.workbenches[event.Workbench].PushRequest(models.Request{
		Task:     event.Source,
		Step:     event.Step,
		Type:     event.Type,
		Caller:   event.Caller,
		Expected: event.Expected,
		IsPickup: event.IsPickup,
	}, event.Fixture)
}

func (a *OrchestratorActor) ComponentPlaced(event *events.OrchestratorEvent) {
	a.workbenches[event.Workbench].SetItem(event.Source, event.Caller, event.Fixture, event.Component)
	a.workbenches[event.Workbench].SetFixtureOwner(event.Source, event.Caller, event.Fixture)
}

func (a *OrchestratorActor) ComponentAttached(event *events.OrchestratorEvent) {
	a.workbenches[event.Workbench].AttachItem(event.Source, event.Caller, event.Fixture, event.Component)
	a.workbenches[event.Workbench].SetFixtureOwner(enums.NoneTask, event.Caller, event.Fixture)
}

func (a *OrchestratorActor) Process(ctx *actor.Context) {
	for _, workbench := range a.workbenches {
		fixtures := workbench.GetFixturesContent()
		workbench.SetLEDs(fixtures)
		for _, fixture := range fixtures {
			request := workbench.PopRequest(fixture.Fixture)
			if request == nil {
				continue
			}
			if !slices.Contains(request.Expected, fixture.Component.Stage()) {
				ctx.Send(ctx.PID(), &events.OrchestratorEvent{
					Source:      request.Task,
					Destination: a.Orchestrator(),
					Type:        enums.FixtureRequested,
					Step:        request.Step,
					Caller:      request.Caller,
					Workbench:   workbench.Key(),
					Fixture:     fixture.Fixture,
					Expected:    request.Expected,
					IsPickup:    request.IsPickup,
				})
				continue
			}

			workbench.SetFixtureOwner(request.Task, request.Caller, fixture.Fixture)
			component := enums.NoneComponent

			if request.Type == enums.FixtureRequested {
				workbench.SetLED(fixture.Fixture, "ASSEMBLING")
			}

			if request.IsPickup {
				component = workbench.GetItem(request.Task, request.Caller, fixture.Fixture)
				workbench.SetFixtureOwner(enums.NoneTask, request.Caller, fixture.Fixture)
			}

			ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
				Source:      a.Orchestrator(),
				Destination: request.Task,
				Step:        request.Step,
				Component:   component,
			})
		}

		if len(fixtures) == 3 && fixtures[2].Component.Stage() == enums.Completed {
			a.numberOfChairs = a.numberOfChairs + 1
			logger.Get().Info("Chair assembled", "Number of chair:", fmt.Sprint(a.numberOfChairs))
			workbench.RemoveCompletedItem()
		}

		if workbench.CanRotate() {
			fixtures := workbench.RotateFixtures()
			//LEDs
			workbench.SetLED(enums.Fixture1, "FREE")

			if fixtures[1].Component.Stage() == enums.LegsAttached {
				ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
					Source:      a.Orchestrator(),
					Destination: enums.AssemblyTask2,
					Step:        enums.Step1,
				})
			}
		}
	}
}
