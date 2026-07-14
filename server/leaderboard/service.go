package leaderboard

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	// General operations
	GetAll(ctx context.Context, cursor int, limit int) ([]Leaderboard, error)
	Get(ctx context.Context, id uuid.UUID) (Leaderboard, error)
	Create(ctx context.Context, dto CreateDTO) (Leaderboard, error)
	Modify(ctx context.Context, dto ModifyDTO) (Leaderboard, error)
	Delete(ctx context.Context, dto DeleteDTO) (Leaderboard, error)

	// Event handlers & listeners
	ListenForEvents(ctx context.Context)
	HandleLeaderboardCreate(ctx context.Context, payload LeaderboardCreatedEventPayload)
	HandleLeaderboardDelete(ctx context.Context, payload LeaderboardDeletedEventPayload)

	// Score operations
	GetScore(ctx context.Context, dto GetScoreDTO) (ScoreDTO, error)
	UpdateScore(ctx context.Context, dto UpdateScoreDTO) error
	DeleteScore(ctx context.Context, dto DeleteScoreDTO) error
}

type LocalService struct {
	logger     *zap.Logger
	metaRepo   MetadataRepository
	scoreRepo  ScoreRepository
	eventsChan chan Event
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewLocalService(ctx context.Context, l *zap.Logger, mr MetadataRepository, sr ScoreRepository) (*LocalService, chan Event) {
	channel := make(chan Event, 100)
	ctx, cancelFunc := context.WithCancel(ctx)
	service := &LocalService{
		logger:     l,
		metaRepo:   mr,
		scoreRepo:  sr,
		eventsChan: channel,
		ctx:        ctx,
		cancel:     cancelFunc,
	}

	return service, channel
}

// General operations
func (s *LocalService) GetAll(ctx context.Context, cursor int, limit int) ([]Leaderboard, error) {
	ctx, span := tracer.Start(ctx, "leaderboard.service.getall")
	defer span.End()

	return s.metaRepo.GetAll(ctx, cursor, limit)
}

func (s *LocalService) Get(ctx context.Context, id uuid.UUID) (Leaderboard, error) {
	ctx, span := tracer.Start(ctx, "leaderboard.service.get")
	defer span.End()

	return s.metaRepo.Get(ctx, id)
}

func (s *LocalService) Create(ctx context.Context, dto CreateDTO) (Leaderboard, error) {
	ctx, span := tracer.Start(ctx, "leaderboard.service.create")
	defer span.End()

	return s.metaRepo.Create(ctx, dto)
}

func (s *LocalService) Modify(ctx context.Context, dto ModifyDTO) (Leaderboard, error) {
	ctx, span := tracer.Start(ctx, "leaderboard.service.modify")
	defer span.End()

	return s.metaRepo.Modify(ctx, dto)
}

func (s *LocalService) Delete(ctx context.Context, dto DeleteDTO) (Leaderboard, error) {
	ctx, span := tracer.Start(ctx, "leaderboard.service.delete")
	defer span.End()

	return s.metaRepo.Delete(ctx, dto)
}

// Events related
func (s *LocalService) ListenForEvents(ctx context.Context) {
	for {
		select {
		case event := <-s.eventsChan:
			switch event.Type() {
			case "created":
				s.HandleLeaderboardCreate(ctx, event.Payload().(LeaderboardCreatedEventPayload))
			case "deleted":
				s.HandleLeaderboardDelete(ctx, event.Payload().(LeaderboardDeletedEventPayload))
			}
		case <-s.ctx.Done():
			close(s.eventsChan)
		}
	}
}

func (s *LocalService) HandleLeaderboardCreate(ctx context.Context, payload LeaderboardCreatedEventPayload) {

}

func (s *LocalService) HandleLeaderboardDelete(ctx context.Context, payload LeaderboardDeletedEventPayload) {

}

// Score operations
func (s *LocalService) GetScore(ctx context.Context, dto GetScoreDTO) (ScoreDTO, error) {
	ctx, span := tracer.Start(ctx, "leaderboard.service.getscore")
	defer span.End()

	return s.scoreRepo.GetScore(ctx, dto)
}

func (s *LocalService) UpdateScore(ctx context.Context, dto UpdateScoreDTO) error {
	ctx, span := tracer.Start(ctx, "leaderboard.service.updatescore")
	defer span.End()

	return s.scoreRepo.UpdateScore(ctx, dto)
}

func (s *LocalService) DeleteScore(ctx context.Context, dto DeleteScoreDTO) error {
	ctx, span := tracer.Start(ctx, "leaderboard.service.deletescore")
	defer span.End()

	return s.scoreRepo.DeleteScore(ctx, dto)
}
