package main

import (
	"os"

	"github.com/anthdm/hollywood/actor"
	AssemblyTaskServices "github.com/thankala/gregor_chair/assembly_task_7/services"
	"github.com/thankala/gregor_chair_common/configuration"
	"github.com/thankala/gregor_chair_common/controllers"
	"github.com/thankala/gregor_chair_common/enums"
	"github.com/thankala/gregor_chair_common/interfaces"
	"github.com/thankala/gregor_chair_common/services"
)

func main() {
	var server interfaces.Server
	var httpClient interfaces.HttpClient

	if os.Getenv("KAFKA_ADDR") != "" {
		kafkaServer := services.NewConfluentKafkaServer(getKafkaOptions()...)
		server = kafkaServer
	}
	if os.Getenv("TCP_ADDR") != "" {
		tcpServer := services.NewTCPServer(getTCPOptions()...)
		server = tcpServer
	}
	if os.Getenv("ROBOT_3_URL") != "" {
		httpClient = services.NewHttpClient(os.Getenv("ROBOT_3_URL"))
	}

	redisStorer := services.NewRedisStore(getRedisOptions()...)
	robot3Controller := controllers.NewRobotController(
		redisStorer,
		httpClient,
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
	e, _ := actor.NewEngine(actor.EngineConfig{})

	e.Spawn(services.NewAssemblyTaskActor(
		AssemblyTaskServices.NewAssemblyTask7Actor(*robot3Controller), server),
		enums.AssemblyTask7.String(),
	)
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

	kafkaOptions = append(kafkaOptions, configuration.WithKafkaTopic(enums.AssemblyTask7.String()))
	kafkaOptions = append(kafkaOptions, configuration.WithKafkaGroupId(enums.AssemblyTask7.String()))

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
