package controllers

import (
	"context"
	"errors"
	pb "search-service/protos/esapi"
	"search-service/repositories"
)

type EsApiController struct {
	Repo *repositories.EsApiRepository
}

func (controller *EsApiController) Search(ctx context.Context, in *pb.SearchRequest) (*pb.PageDTO, error) {
	if in == nil {
		return nil, errors.New("searchRequest is nil")
	}

	current := int32(1)
	size := int32(20)
	if in.GetCurrent() > 0 {
		current = in.GetCurrent()
	}
	if in.GetSize() > 0 {
		size = in.GetSize()
	}

	result, err := controller.Repo.SearchNews(ctx, in.GetExcludeCategoryIds(), in.GetKeyword(), current, size)
	if err != nil {
		return nil, errors.New("error occurred in Search::" + err.Error())
	}

	var pages = int(result.EsResult.TotalHits())/int(size) + 1

	if int(result.EsResult.TotalHits())%int(size) == 0 {
		pages = int(result.EsResult.TotalHits()) / int(size)
	}

	return &pb.PageDTO{
		Current: current,
		Size:    size,
		Total:   int32(result.EsResult.TotalHits()),
		Pages:   int32(pages),
		Records: result.News,
	}, nil
}

func (controller *EsApiController) Upsert(ctx context.Context, in *pb.News) (*pb.CommonResponse, error) {
	if in == nil {
		return nil, errors.New("news is nil")
	}

	_, err := controller.Repo.UpsertNews(ctx, in)
	if err != nil {
		return nil, errors.New("error occurred in Upsert::" + err.Error())
	}

	return &pb.CommonResponse{
		IsSuccess: true,
	}, nil
}

func (controller *EsApiController) Delete(ctx context.Context, in *pb.NewsId) (*pb.CommonResponse, error) {
	if in == nil {
		return nil, errors.New("newsId is nil")
	}

	_, err := controller.Repo.DeleteNews(ctx, in.GetId())
	if err != nil {
		return nil, errors.New("error occurred in Delete::" + err.Error())
	}

	return &pb.CommonResponse{
		IsSuccess: true,
	}, nil
}
