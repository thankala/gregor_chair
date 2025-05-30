package main

import (
	"os"

	"github.com/anthdm/hollywood/actor"
	AssemblyTaskServices "github.com/thankala/gregor_chair/orchestrator/services"
	"github.com/thankala/gregor_chair_common/configuration"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/services"
)

func main() {
	var httpClient interfaces.HttpClient

	redisStorer := services.NewRedisStore(getRedisOptions()...)

	if os.Getenv("WORKBENCH_URL") != "" {
		httpClient = services.NewHttpClient(os.Getenv("WORKBENCH_URL"))
	}

	workbench1Controller := controllers.NewWorkbenchController(
		redisStorer,
		httpClient,
		configuration.WithWorkbenchKey(enums.Workbench1),
		configuration.WithFixture(
			*configuration.NewFixtureConfiguration(enums.Fixture1, []string{enums.Robot1.String()}),
			*configuration.NewFixtureConfiguration(enums.Fixture2, []string{enums.Robot2.String()}),
			*configuration.NewFixtureConfiguration(enums.Fixture3, []string{enums.Robot3.String()}),
		),
		configuration.WithStateMapping(map[enums.Fixture]map[enums.Stage]string{
			enums.Fixture1: {
				// enums.Initial: "FREE",
				// enums.LegsAttached: "ASSEMBLING",
				enums.LegsAttached: "COMPLETED",
			},
			enums.Fixture2: {
				// enums.Initial: "FREE",
				// enums.BaseAttached:    "ASSEMBLING",
				// enums.CastorsAttached: "ASSEMBLING",
				// enums.LiftAttached: "PENDING",
				// enums.SeatAttached:   "ASSEMBLING",
				enums.SeatAttached: "COMPLETED",
			},
			enums.Fixture3: {
				// enums.Initial:        "FREE",
				// enums.ScrewsAttached: "ASSEMBLING",
				// enums.BackAttached:     "COMPLETED",
				// enums.LeftArmAttached:  "ASSEMBLING",
				// enums.RightArmAttached: "ASSEMBLING",
				enums.Completed: "COMPLETED",
			},
		}),
	)
	workbench2Controller := controllers.NewWorkbenchController(
		redisStorer,
		httpClient,
		configuration.WithWorkbenchKey(enums.Workbench2),
		configuration.WithFixture(
			*configuration.NewFixtureConfiguration(enums.Fixture1, []string{enums.Robot1.String(), enums.Robot2.String()}),
		),
		configuration.WithStateMapping(map[enums.Fixture]map[enums.Stage]string{
			enums.Fixture1: {
				enums.Initial:           "FREE",
				enums.InitialSeat:       "ASSEMBLING",
				enums.SeatPlateAttached: "PENDING",
				// enums.SeatScrewsAttached: "ASSEMBLING",
			},
		}),
	)
	e, _ := actor.NewEngine(actor.EngineConfig{})

	if os.Getenv("KAFKA_ADDR") != "" {
		e.Spawn(services.NewOrchestratorActor[AssemblyTaskServices.OrchestratorActor](
			AssemblyTaskServices.NewOrchestratorActor(*workbench1Controller, *workbench2Controller), services.NewConfluentKafkaServer(getKafkaOptions()...)),
			enums.Orchestrator.String(),
		)
	}
	if os.Getenv("TCP_ADDR") != "" {
		e.Spawn(services.NewOrchestratorActor[AssemblyTaskServices.OrchestratorActor](
			AssemblyTaskServices.NewOrchestratorActor(*workbench1Controller, *workbench2Controller), services.NewTCPServer(getTCPOptions()...)),
			enums.Orchestrator.String(),
		)
	}

	<-make(chan struct{})
}

func getRedisOptions() []configuration.RedisOptFunc {
	var redisOptions []configuration.RedisOptFunc
	// Conditionally add options based on environment variable presence
	if os.Getenv("REDIS_ADDR") != "" {
		redisOptions = append(redisOptions, configuration.WithStoreAddr(os.Getenv("REDIS_ADDR")))
	}
	if os.Getenv("REDIS_PASSWORD") != "" {
		redisOptions = append(redisOptions, configuration.WithStorePassword(os.Getenv("REDIS_PASSWORD")))
	}
	if os.Getenv("REDIS_DB") != "" {
		redisOptions = append(redisOptions, configuration.WithStoreDb(os.Getenv("REDIS_DB")))
	}
	return redisOptions
}

func getKafkaOptions() []configuration.KafkaOptionFunc {
	var kafkaOptions []configuration.KafkaOptionFunc

	kafkaOptions = append(kafkaOptions, configuration.WithKafkaTopic(enums.Orchestrator.String()))
	kafkaOptions = append(kafkaOptions, configuration.WithKafkaGroupId(enums.Orchestrator.String()))

	// Conditionally add options based on environment variable presence
	if os.Getenv("KAFKA_ADDR") != "" {
		kafkaOptions = append(kafkaOptions, configuration.WithKafkaListenAddresses(os.Getenv("KAFKA_ADDR")))
	}
	if os.Getenv("KAFKA_TOPIC") != "" {
		kafkaOptions = append(kafkaOptions, configuration.WithKafkaTopic(os.Getenv("KAFKA_TOPIC")))
	}
	if os.Getenv("KAFKA_GROUP_ID") != "" {
		kafkaOptions = append(kafkaOptions, configuration.WithKafkaGroupId(os.Getenv("KAFKA_GROUP_ID")))
	}

	return kafkaOptions
}

func getTCPOptions() []configuration.TcpOptFunc {
	var tcpOptions []configuration.TcpOptFunc
	// Conditionally add options based on environment variable presence
	if os.Getenv("TCP_ADDR") != "" {
		tcpOptions = append(tcpOptions, configuration.WithAddress(os.Getenv("TCP_ADDR")))
	}
	return tcpOptions
}
