package grpc_services

import (
	"context"
	"search-service/controllers"
	pb "search-service/protos/esapi"
	"search-service/utils"
)

type EsApiServer struct {
	pb.UnimplementedEsApiServer
	Controller *controllers.EsApiController
}

func (s *EsApiServer) Search(ctx context.Context, in *pb.SearchRequest) (*pb.PageDTO, error) {
	r, err := s.Controller.Search(ctx, in)
	if err != nil {
		utils.Logger.Error("EsApiServer Search::" + err.Error())
		return &pb.PageDTO{}, nil
	}
	return r, nil
}

func (s *EsApiServer) Upsert(ctx context.Context, in *pb.News) (*pb.CommonResponse, error) {
	r, err := s.Controller.Upsert(ctx, in)
	if err != nil {
		utils.Logger.Error("EsApiServer Upsert::" + err.Error())
		return &pb.CommonResponse{}, nil
	}
	return r, nil
}

func (s *EsApiServer) Delete(ctx context.Context, in *pb.NewsId) (*pb.CommonResponse, error) {
	r, err := s.Controller.Delete(ctx, in)
	if err != nil {
		utils.Logger.Error("EsApiServer Delete::" + err.Error())
		return &pb.CommonResponse{}, nil
	}
	return r, nil
}

func (s *EsApiServer) Test(ctx context.Context, in *pb.Empty) (*pb.CommonResponse, error) {
	//log.Fatal("Testing!!!")

	return &pb.CommonResponse{
		IsSuccess: false,
	}, nil
}
