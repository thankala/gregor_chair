package main

import (
	"os"

	"github.com/anthdm/hollywood/actor"
	AssemblyTaskServices "github.com/thankala/gregor_chair/assembly_task_3/services"
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
	if os.Getenv("ROBOT_URL") != "" {
		httpClient = services.NewHttpClient(os.Getenv("ROBOT_URL"))
	}

	redisStorer := services.NewRedisStore(getRedisOptions()...)

	robot2Controller := controllers.NewRobotController(
		redisStorer,
		httpClient,
		configuration.WithRobotKey(enums.Robot2.String()),
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
	e, _ := actor.NewEngine(actor.EngineConfig{})
	e.Spawn(services.NewAssemblyTaskActor[AssemblyTaskServices.AssemblyTask3Actor](
		AssemblyTaskServices.NewAssemblyTask3Actor(*robot2Controller), server),
		enums.AssemblyTask3.String(),
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
	kafkaOptions = append(kafkaOptions, configuration.WithKafkaTopic(enums.AssemblyTask3.String()))
	kafkaOptions = append(kafkaOptions, configuration.WithKafkaGroupId(enums.AssemblyTask3.String()))

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
