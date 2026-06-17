package grpc

import (
	"context"
	"event-driven/internal/models"
	poll "event-driven/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PollServer struct {
	poll.UnimplementedPollServiceServer
	repo models.PollRepository
}

func NewPollServer(repo models.PollRepository) *PollServer {
	return &PollServer{
		repo: repo,
	}
}

// gRPC method CreatePoll
func (s *PollServer) CreatePoll(ctx context.Context, req *poll.CreatePollRequest) (*poll.CreatePollResponse, error) {
	// proto struct
	title, description := req.Title, req.Description

	createdPoll, err := s.repo.Create(ctx, title, description)
	if err != nil {
		// gRPC status Internal Server Error
		return nil, status.Errorf(codes.Internal, "CreatePoll: failed to create poll: %v", err)
	}

	response := &poll.CreatePollResponse{
		Id: createdPoll,
	}

	return response, nil
}

func (s *PollServer) GetPoll(ctx context.Context, req *poll.GetPollRequest) (*poll.GetPollResponse, error) {
	id := req.Id

	gotPoll, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetPoll: failed to get poll: %v", err)
	}

	response := &poll.GetPollResponse{
		Poll: &poll.Poll{
			Id:          gotPoll.ID,
			Title:       gotPoll.Title,
			Description: *gotPoll.Description,
			IsActive:    gotPoll.IsActive,
			CreatedAt:   gotPoll.CreatedAt.String(),
		},
	}

	return response, nil
}

func (s *PollServer) UpdatePoll(ctx context.Context, req *poll.UpdatePollRequest) (*poll.UpdatePollResponse, error) {
	id, title, description := req.Id, req.Title, req.Description
	isActive := req.IsActive

	updatedPoll, err := s.repo.Update(ctx, id, title, description, isActive)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UpdatePoll: failed to update poll: %v", err)
	}

	response := &poll.UpdatePollResponse{
		Poll: &poll.Poll{
			Id:          updatedPoll.ID,
			Title:       updatedPoll.Title,
			Description: *updatedPoll.Description,
			IsActive:    updatedPoll.IsActive,
			CreatedAt:   updatedPoll.CreatedAt.String(),
		},
	}

	return response, nil
}

func (s *PollServer) DeletePoll(ctx context.Context, req *poll.DeletePollRequest) (*emptypb.Empty, error) {
	id := req.Id

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeletePoll: failed to delete poll: %v", err)
	}

	return nil, nil
}

func (s *PollServer) ListPoll(ctx context.Context, req *poll.ListPollRequest) (*poll.ListPollResponse, error) {
	limit, offset := req.Limit, req.Offset
	onlyActive := req.OnlyActive

	listedPolls, err := s.repo.List(ctx, int(limit), int(offset), onlyActive)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ListPoll: failed to list polls: %v", err)
	}

	var protoPolls []*poll.Poll
	for _, p := range listedPolls {
		protoPolls = append(protoPolls, &poll.Poll{
			Id:          p.ID,
			Title:       p.Title,
			Description: *p.Description,
			IsActive:    p.IsActive,
			CreatedAt:   p.CreatedAt.String(),
		})
	}

	response := &poll.ListPollResponse{
		Polls: protoPolls,
	}

	return response, nil
}
