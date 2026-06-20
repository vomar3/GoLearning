package grpc

import (
	"context"
	"errors"
	"event-driven/internal/models"
	"event-driven/internal/repository"
	"event-driven/proto/vote"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type VoteServer struct {
	vote.UnimplementedVoteServiceServer
	repo       models.VoteRepository
	optionRepo models.OptionRepository
	redis      *redis.Client
}

func NewVoteServer(repo models.VoteRepository, optionRepo models.OptionRepository, redis *redis.Client) *VoteServer {
	return &VoteServer{
		repo:       repo,
		optionRepo: optionRepo,
		redis:      redis,
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

	key := "leaderboard:poll:" + pollId
	if err := v.redis.ZIncrBy(ctx, key, 1, optionId).Err(); err != nil {
	}
	if err := v.redis.Publish(ctx, key, "updated").Err(); err != nil {
	}

	return &emptypb.Empty{}, nil
}

func (v *VoteServer) SubscribeLeaderboard(req *vote.SubscribeLeaderboardRequest, stream vote.VoteService_SubscribeLeaderboardServer) error {
	ctx := stream.Context()
	pollID := req.PollId
	topN := int(req.TopN)
	if topN <= 0 {
		topN = 10
	}

	// upload data from db
	options, err := v.optionRepo.ListOptions(ctx, pollID)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to load options: %v", err)
	}

	optionText := make(map[string]string)
	for _, opt := range options {
		optionText[opt.ID] = opt.Text
	}

	key := "leaderboard:poll:" + pollID

	// currently leaderboard
	getTop := func() ([]*vote.LeaderboardEntry, error) {
		res, err := v.redis.ZRevRangeWithScores(ctx, key, 0, int64(topN-1)).Result()
		if err != nil {
			return nil, err
		}

		entries := make([]*vote.LeaderboardEntry, 0, len(res))
		for rank, z := range res {
			optionID := z.Member.(string)
			votes := int64(z.Score)
			text := optionText[optionID]
			entries = append(entries, &vote.LeaderboardEntry{
				OptionId: optionID,
				Text:     text,
				Votes:    votes,
				Rank:     int32(rank + 1),
			})
		}

		return entries, nil
	}

	entries, err := getTop()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get initial leaderboard: %v", err)
	}

	if err := stream.Send(&vote.LeaderboardUpdate{Entries: entries}); err != nil {
		return err
	}

	pubsub := v.redis.Subscribe(ctx, key)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ch:
			entries, err := getTop()
			if err != nil {
				continue
			}

			if err := stream.Send(&vote.LeaderboardUpdate{Entries: entries}); err != nil {
				return err
			}
		}
	}
}
