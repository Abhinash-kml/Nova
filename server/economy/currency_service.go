package economy

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"go.uber.org/zap"
)

type CurrencyService interface {
	Get(ctx context.Context, req GetCurrencyRequest) (Currency, error)
	GetAll(ctx context.Context, req GetAllCurrencyRequest) ([]Currency, error)
	Create(ctx context.Context, req CreateCurrencyRequest) (Currency, error)
	Update(ctx context.Context, req UpdateCurrencyRequest) (Currency, error)
	Delete(ctx context.Context, req DeleteCurrencyRequest) (int, error)
}

type LocalCurrencyService struct {
	repo   CurrencyRepository
	config *config.CurrencyConfig
	logger *zap.Logger
}

func NewLocalCurrencyService(repo CurrencyRepository, c *config.CurrencyConfig, l *zap.Logger) *LocalCurrencyService {
	return &LocalCurrencyService{
		repo:   repo,
		config: c,
		logger: l,
	}
}

func (s *LocalCurrencyService) Get(ctx context.Context, req GetCurrencyRequest) (Currency, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.service.Get")
	defer span.End()

	currency, err := s.repo.Get(ctx, req)
	if err != nil {
		return Currency{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get currency")
	}

	return currency, nil
}

func (s *LocalCurrencyService) GetAll(ctx context.Context, req GetAllCurrencyRequest) ([]Currency, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.service.GetAll")
	defer span.End()

	currencies, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger)
	}

	return currencies, nil
}

func (s *LocalCurrencyService) Create(ctx context.Context, req CreateCurrencyRequest) (Currency, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.service.Create")
	defer span.End()

	currency, err := s.repo.Create(ctx, req)
	if err != nil {
		return Currency{}, common.TranslatePostgresError(err, s.logger)
	}

	return currency, nil
}

func (s *LocalCurrencyService) Update(ctx context.Context, req UpdateCurrencyRequest) (Currency, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.service.Update")
	defer span.End()

	currency, err := s.repo.Update(ctx, req)
	if err != nil {
		return Currency{}, common.TranslatePostgresError(err, s.logger)
	}

	return currency, nil
}

func (s *LocalCurrencyService) Delete(ctx context.Context, req DeleteCurrencyRequest) (int, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.service.Delete")
	defer span.End()

	id, err := s.repo.Delete(ctx, req)
	if err != nil {
		return 0, common.TranslatePostgresError(err, s.logger)
	}

	return id, nil
}
