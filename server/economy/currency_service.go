package economy

import (
	"context"

	"github.com/abhinash-kml/nova/server/config"
	"go.uber.org/zap"
)

type CurrencyService interface {
	Get(ctx context.Context, dto GetCurrencyDTO) (CurrencyDTO, error)
	GetAll(ctx context.Context, dto GetAllCurrencyDTO) ([]CurrencyDTO, error)
	Create(ctx context.Context, dto CreateCurrencyDTO) (CurrencyDTO, error)
	Update(ctx context.Context, dto UpdateCurrencyDTO) (CurrencyDTO, error)
	Delete(ctx context.Context, dto DeleetCurrencyDTO) (int, error)
}

type LocalCurrencyService struct {
	repo   CurrencyRepository
	config *config.Config
	logger *zap.Logger
}

func NewLocalCurrencyService(repo CurrencyRepository, c *config.Config, l *zap.Logger) *LocalCurrencyService {
	return &LocalCurrencyService{
		repo:   repo,
		config: c,
		logger: l,
	}
}

func (s *LocalCurrencyService) Get(ctx context.Context, dto GetCurrencyDTO) (CurrencyDTO, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.get")
	defer span.End()

	return s.repo.Get(ctx, dto)
}

func (s *LocalCurrencyService) GetAll(ctx context.Context, dto GetAllCurrencyDTO) ([]CurrencyDTO, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.getall")
	defer span.End()

	return s.repo.GetAll(ctx, dto)
}

func (s *LocalCurrencyService) Create(ctx context.Context, dto CreateCurrencyDTO) (CurrencyDTO, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.get")
	defer span.End()

	return s.repo.Create(ctx, dto)
}

func (s *LocalCurrencyService) Update(ctx context.Context, dto UpdateCurrencyDTO) (CurrencyDTO, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.get")
	defer span.End()

	return s.repo.Update(ctx, dto)
}

func (s *LocalCurrencyService) Delete(ctx context.Context, dto DeleetCurrencyDTO) (int, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.get")
	defer span.End()

	return s.repo.Delete(ctx, dto)
}
