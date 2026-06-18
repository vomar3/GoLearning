package grpc

import (
	"context"
	"errors"
	"event-driven/internal/models"
	"event-driven/internal/repository"
	"event-driven/proto/vote"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type VoteServer struct {
	vote.UnimplementedVoteServiceServer
	repo models.VoteRepository
}

func NewVoteServer(repo models.VoteRepository) *VoteServer {
	return &VoteServer{
		repo: repo,
	}
}

func (v *VoteServer) CastVote(ctx context.Context, req *vote.CastVoteRequest) (*emptypb.Empty, error) {
	if req == nil || req.PollId == "" || req.OptionId == "" || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "CastVote: poll_id, option_id, user_id are required")
	}

	pollId, optionId, userId := req.PollId, req.OptionId, req.UserId
	err := v.repo.Cast(ctx, pollId, optionId, userId)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "CastVote: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "CastVote: failed to cast data: %v", err)
	}

	return &emptypb.Empty{}, nil
}
