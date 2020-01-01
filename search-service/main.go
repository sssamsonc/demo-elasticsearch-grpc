package main

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"search-service/configs"
	"search-service/controllers"
	"search-service/databases"
	"search-service/grpc_services"
	"search-service/protos/esapi"
	"search-service/repositories"
	"search-service/utils"
	"time"
)

func main() {
	time.Local = time.UTC //make sure the timezone is UTC in this project

	//init logger
	logger := utils.NewLogger()
	defer logger.Sync()

	//init db
	esClient, err := databases.NewElasticSearchClient()
	if err != nil {
		utils.Logger.Error("failed to init elastic search connection::" + err.Error())
		return
	}
	//init repositories
	esApiRepository := &repositories.EsApiRepository{EsClient: esClient}
	//init controllers
	esApiController := &controllers.EsApiController{Repo: esApiRepository}

	//init the gRPC
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(utils.Logger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	reflection.Register(s)

	esapi.RegisterEsApiServer(s, &grpc_services.EsApiServer{Controller: esApiController})

	lis, err := net.Listen("tcp", ":"+configs.Get("GRPC_PORT"))
	if err != nil {
		utils.Logger.Error("failed to listen port" + configs.Get("GRPC_PORT") + "::" + err.Error())
		return
	}

	if err := s.Serve(lis); err != nil {
		utils.Logger.Error("failed to serve::" + err.Error())
		return
	}
}
