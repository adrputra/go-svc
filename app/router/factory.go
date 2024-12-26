package router

import (
	"face-recognition-svc/app/client"
	"face-recognition-svc/app/config"
	"face-recognition-svc/app/controller"
	"face-recognition-svc/app/service"

	"github.com/aws/aws-sdk-go/service/s3"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

type ServiceFactory struct {
	user  service.InterfaceUserService
	role  service.InterfaceRoleService
	param service.InterfaceParamService
}

type ControllerFactory struct {
	user  controller.InterfaceUserController
	role  controller.InterfaceRoleController
	param controller.InterfaceParamController
}

type ClientFactory struct {
	user    client.InterfaceUserClient
	storage client.InterfaceStorageClient
	role    client.InterfaceRoleClient
	param   client.InterfaceParamClient
}

type Factory struct {
	Service    ServiceFactory
	Controller ControllerFactory
	Client     ClientFactory
}

var factory *Factory

func InitFactory(cfg *config.Config, db *gorm.DB, s3 *s3.S3, redis *redis.Client, mq *amqp.Channel) {
	client := ClientFactory{
		user:    client.NewUserClient(db, cfg),
		storage: client.NewStorageClient(s3, db),
		role:    client.NewRoleClient(db),
		param:   client.NewParamClient(db),
	}
	controller := ControllerFactory{
		user:  controller.NewUserController(client.user, client.role),
		role:  controller.NewRoleController(client.role),
		param: controller.NewParamController(redis, client.param),
	}
	service := ServiceFactory{
		user:  service.NewUserService(controller.user),
		role:  service.NewRoleService(controller.role),
		param: service.NewParamService(controller.param),
	}
	factory = &Factory{
		Service:    service,
		Controller: controller,
		Client:     client,
	}
}
