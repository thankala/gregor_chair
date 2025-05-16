package services

import (
	"slices"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/events"
	"github.com/thankala/gregor_chair_common/logger"
	"github.com/thankala/gregor_chair_common/models"
)

type OrchestratorActor struct {
	workbenches map[enums.Workbench]*controllers.WorkbenchController
}

func NewOrchestratorActor(workbenchControllers ...controllers.WorkbenchController) *OrchestratorActor {
	workbenches := make(map[enums.Workbench]*controllers.WorkbenchController)
	for _, workbench := range workbenchControllers {
		workbenches[workbench.Key()] = &workbench
	}
	return &OrchestratorActor{
		workbenches: workbenches,
	}
}

func (a *OrchestratorActor) Orchestrator() enums.Task {
	return enums.Orchestrator
}

func (a *OrchestratorActor) FixtureRequested(event *events.OrchestratorEvent) {
	a.workbenches[event.Workbench].PushRequest(models.Request{
		Task:     event.Source,
		Step:     event.Step,
		Caller:   event.Caller,
		Expected: event.Expected,
		IsPickup: event.IsPickup,
	}, event.Fixture)
}

func (a *OrchestratorActor) ComponentPlaced(event *events.OrchestratorEvent) {
	a.workbenches[event.Workbench].SetItem(event.Source, event.Caller, event.Fixture, event.Component)
	a.workbenches[event.Workbench].SetFixtureOwner(enums.NoneTask, event.Caller, event.Fixture)
}

func (a *OrchestratorActor) ComponentAttached(event *events.OrchestratorEvent) {
	a.workbenches[event.Workbench].AttachItem(event.Source, event.Caller, event.Fixture, event.Component)
	a.workbenches[event.Workbench].SetFixtureOwner(enums.NoneTask, event.Caller, event.Fixture)
}

func (a *OrchestratorActor) Process(ctx *actor.Context) {
	for _, workbench := range a.workbenches {
		fixtures := workbench.GetFixturesContent()
		// workbench.SetLEDs(fixtures)
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
			if request.Task == enums.AssemblyTask4 && request.Step == enums.Step2 {
				workbench.SetLED(enums.Fixture1, "ASSEMBLING")
			}
			if request.Task == enums.AssemblyTask3 && request.Step == enums.Step2 {
				workbench.SetLED(enums.Fixture1, "ASSEMBLING")
			}

			if request.Task == enums.AssemblyTask4 && request.Step == enums.Step2 {
				workbench.SetLED(enums.Fixture1, "ASSEMBLING")
			}
			if request.Task == enums.AssemblyTask3 && request.Step == enums.Step2 {
				workbench.SetLED(enums.Fixture1, "ASSEMBLING")
			}

			workbench.SetFixtureOwner(request.Task, request.Caller, fixture.Fixture)
			component := enums.NoneComponent
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
		if len(fixtures) == 3 && fixtures[2].Component.Stage() == enums.Completed && workbench.Key() == enums.Workbench1 {
			logger.Get().Info("One chair has been assembled", "timestamp", time.Now())
			workbench.RemoveCompletedItem()
		}

		if workbench.CanRotate() {
			fixtures := workbench.RotateFixtures()
			//LEDs
			// a.workbench.SetLEDs(fixtures)

			if fixtures[1].Component.Stage() == enums.BaseAttached {
				workbench.SetLED(fixtures[1].Fixture, "ASSEMBLING")
				ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
					Source:      a.Orchestrator(),
					Destination: enums.AssemblyTask2,
					Step:        enums.Step1,
				})
			}

			if fixtures[2].Component.Stage() == enums.ScrewsAttached {
				workbench.SetLED(fixtures[2].Fixture, "ASSEMBLING")
				ctx.Send(ctx.PID(), &events.AssemblyTaskEvent{
					Source:      a.Orchestrator(),
					Destination: enums.AssemblyTask6,
					Step:        enums.Step1,
				})
			}
		}
	}
}
