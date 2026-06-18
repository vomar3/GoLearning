package grpc

import (
	"context"
	"event-driven/internal/models"
	"event-driven/proto/option"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OptionServer struct {
	option.UnimplementedOptionServiceServer
	repo models.OptionRepository
}

func NewOptionServer(repo models.OptionRepository) *OptionServer {
	return &OptionServer{
		repo: repo,
	}
}

func (s *OptionServer) CreateOption(ctx context.Context, req *option.CreateOptionRequest) (*option.CreateOptionResponse, error) {
	poll_id, text := req.PollId, req.Text

	created, err := s.repo.Create(ctx, poll_id, text)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateOption: failed to create option: %v", err)
	}

	response := &option.CreateOptionResponse{
		Id: created,
	}

	return response, nil
}

func (s *OptionServer) ListOptions(ctx context.Context, req *option.ListOptionsRequest) (*option.ListOptionsResponse, error) {
	poll_id := req.PollId

	listed, err := s.repo.ListOptions(ctx, poll_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ListOptions: failed to list options: %v", err)
	}

	var listAllOptions []*option.Option

	for _, list := range listed {
		listAllOptions = append(listAllOptions, &option.Option{
			Id:        list.ID,
			PollId:    list.PollID,
			Text:      list.Text,
			VoteCount: list.VotesCount,
		})
	}

	response := &option.ListOptionsResponse{
		Options: listAllOptions,
	}

	return response, nil
}
