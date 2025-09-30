package main

import (
	"os"

	"github.com/anthdm/hollywood/actor"
	OrchestratorService "github.com/thankala/gregor_chair/orchestrator/services"
	"github.com/thankala/gregor_chair_common/configuration"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/services"
)

func main() {
	var robot1HttpClient interfaces.HttpClient
	var robot2HttpClient interfaces.HttpClient
	var robot3HttpClient interfaces.HttpClient
	var workbench1HttpClient interfaces.HttpClient
	var workbench2HttpClient interfaces.HttpClient

	redisStorer := services.NewRedisStore(getRedisOptions()...)
	if os.Getenv("ROBOT_1_URL") != "" {
		robot1HttpClient = services.NewHttpClient(os.Getenv("ROBOT_1_URL"))
	}

	if os.Getenv("ROBOT_2_URL") != "" {
		robot2HttpClient = services.NewHttpClient(os.Getenv("ROBOT_2_URL"))
	}

	if os.Getenv("ROBOT_3_URL") != "" {
		robot3HttpClient = services.NewHttpClient(os.Getenv("ROBOT_3_URL"))
	}

	if os.Getenv("WORKBENCH_1_URL") != "" {
		workbench1HttpClient = services.NewHttpClient(os.Getenv("WORKBENCH_1_URL"))
	}

	if os.Getenv("WORKBENCH_2_URL") != "" {
		workbench2HttpClient = services.NewHttpClient(os.Getenv("WORKBENCH_2_URL"))
	}

	robot1Controller := controllers.NewRobotController(
		redisStorer,
		robot1HttpClient,
		configuration.WithRobotKey(enums.Robot1),
		configuration.WithStorages(
			*configuration.NewStorageConfiguration(
				enums.StorageB1,
				enums.Position1,
				enums.Legs,
				configuration.NewLocation(float64(210), float64(90), float64(20), float64(0)),
			),
			*configuration.NewStorageConfiguration(
				enums.StorageB2,
				enums.Position1,
				enums.Base,
				configuration.NewLocation(float64(200), float64(180), float64(20), float64(0)),
			),
			*configuration.NewStorageConfiguration(
				enums.StorageB3,
				enums.Position2,
				enums.SeatPlate,
				configuration.NewLocation(float64(290), float64(0), float64(20), float64(0))),
		),
		configuration.WithWorkbenches(
			*configuration.NewWorkbenchConfiguration(
				enums.Workbench1,
				enums.Position1,
				enums.Fixture1,
				configuration.NewLocation(float64(260), float64(110), float64(90), float64(0)),
			),
			*configuration.NewWorkbenchConfiguration(
				enums.Workbench2,
				enums.Position2,
				enums.Fixture1,
				configuration.NewLocation(float64(270), float64(-90), float64(40), float64(0)),
			),
		),
		configuration.WithConveyorBelts(
			*configuration.NewConveyorBeltConfiguration(
				enums.ConveyorBelt1,
				enums.Position2,
				enums.Seat,
				false,
				configuration.NewLocation(float64(230), float64(-120), float64(40), float64(0)),
			),
		),
	)

	robot2Controller := controllers.NewRobotController(
		redisStorer,
		robot2HttpClient,
		configuration.WithRobotKey(enums.Robot2),
		configuration.WithStorages(
			*configuration.NewStorageConfiguration(
				enums.StorageB4,
				enums.Position1,
				enums.Castors,
				configuration.NewLocation(float64(170), float64(-160), float64(20), float64(0)),
			),
			*configuration.NewStorageConfiguration(
				enums.StorageB5,
				enums.Position1,
				enums.Lift,
				configuration.NewLocation(float64(290), float64(20), float64(20), float64(0)),
			),
		),
		configuration.WithWorkbenches(
			*configuration.NewWorkbenchConfiguration(
				enums.Workbench1,
				enums.Position1,
				enums.Fixture2,
				configuration.NewLocation(float64(260), float64(-110), float64(90), float64(0)),
			),
			*configuration.NewWorkbenchConfiguration(
				enums.Workbench2,
				enums.Position2,
				enums.Fixture1,
				configuration.NewLocation(float64(270), float64(80), float64(40), float64(0)),
			),
		),
	)

	robot3Controller := controllers.NewRobotController(
		redisStorer,
		robot3HttpClient,
		configuration.WithRobotKey(enums.Robot3),
		configuration.WithStorages(
			*configuration.NewStorageConfiguration(
				enums.StorageB6L,
				enums.Position1,
				enums.LeftArm,
				configuration.NewLocation(float64(210), float64(170), float64(20), float64(0)),
			),
			*configuration.NewStorageConfiguration(
				enums.StorageB6R,
				enums.Position1,
				enums.RightArm,
				configuration.NewLocation(float64(160), float64(180), float64(20), float64(0)),
			),
		),
		configuration.WithWorkbenches(
			*configuration.NewWorkbenchConfiguration(
				enums.Workbench1,
				enums.Position1,
				enums.Fixture3,
				configuration.NewLocation(float64(260), float64(0), float64(90), float64(0)),
			),
		),
		configuration.WithConveyorBelts(
			*configuration.NewConveyorBeltConfiguration(
				enums.ConveyorBelt2,
				enums.Position1,
				enums.Back,
				false,
				configuration.NewLocation(float64(220), float64(190), float64(20), float64(0)),
			),
			*configuration.NewConveyorBeltConfiguration(
				enums.ConveyorBelt3,
				enums.Position1,
				enums.NoneComponent,
				true,
				configuration.NewLocation(float64(220), float64(190), float64(20), float64(0)),
			),
		),
	)

	workbench1Controller := controllers.NewWorkbenchController(
		redisStorer,
		workbench1HttpClient,
		configuration.WithWorkbenchKey(enums.Workbench1),
		configuration.WithFixture(
			*configuration.NewFixtureConfiguration(
				enums.Fixture1,
				[]string{enums.Robot1.String()},
			),
			*configuration.NewFixtureConfiguration(
				enums.Fixture2,
				[]string{enums.Robot2.String()},
			),
			*configuration.NewFixtureConfiguration(
				enums.Fixture3,
				[]string{enums.Robot3.String()},
			),
		),
	)
	workbench2Controller := controllers.NewWorkbenchController(
		redisStorer,
		workbench2HttpClient,
		configuration.WithWorkbenchKey(enums.Workbench2),
		configuration.WithFixture(
			*configuration.NewFixtureConfiguration(
				enums.Fixture1,
				[]string{enums.Robot1.String(), enums.Robot2.String()},
			),
		),
	)

	e, _ := actor.NewEngine(actor.EngineConfig{})

	if os.Getenv("KAFKA_ADDR") != "" {
		e.Spawn(services.NewOrchestratorActor[OrchestratorService.OrchestratorActor](
			OrchestratorService.NewOrchestratorActor(
				[]controllers.WorkbenchController{*workbench1Controller, *workbench2Controller},
				[]controllers.RobotController{*robot1Controller, *robot2Controller, *robot3Controller}),
			services.NewConfluentKafkaServer(getKafkaOptions()...)),
			enums.Orchestrator.String())
	}
	if os.Getenv("TCP_ADDR") != "" {
		e.Spawn(services.NewOrchestratorActor[OrchestratorService.OrchestratorActor](
			OrchestratorService.NewOrchestratorActor(
				[]controllers.WorkbenchController{*workbench1Controller, *workbench2Controller},
				[]controllers.RobotController{*robot1Controller, *robot2Controller, *robot3Controller}),
			services.NewTCPServer(getTCPOptions()...)),
			enums.Orchestrator.String())
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
