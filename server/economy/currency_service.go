package economy

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
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
	ctx, span := tracer.Start(ctx, "economy.currency.service.get")
	defer span.End()

	currency, err := s.repo.Get(ctx, dto)
	if err != nil {
		return CurrencyDTO{}, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get currency")
	}

	return currency, nil
}

func (s *LocalCurrencyService) GetAll(ctx context.Context, dto GetAllCurrencyDTO) ([]CurrencyDTO, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.service.getall")
	defer span.End()

	currencies, err := s.repo.GetAll(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger)
	}

	return currencies, nil
}

func (s *LocalCurrencyService) Create(ctx context.Context, dto CreateCurrencyDTO) (CurrencyDTO, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.service.create")
	defer span.End()

	currency, err := s.repo.Create(ctx, dto)
	if err != nil {
		return CurrencyDTO{}, common.TranslatePostgresError(err, s.logger)
	}

	return currency, nil
}

func (s *LocalCurrencyService) Update(ctx context.Context, dto UpdateCurrencyDTO) (CurrencyDTO, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.service.update")
	defer span.End()

	currency, err := s.repo.Update(ctx, dto)
	if err != nil {
		return CurrencyDTO{}, common.TranslatePostgresError(err, s.logger)
	}

	return currency, nil
}

func (s *LocalCurrencyService) Delete(ctx context.Context, dto DeleetCurrencyDTO) (int, error) {
	ctx, span := tracer.Start(ctx, "economy.currency.service.delete")
	defer span.End()

	id, err := s.repo.Delete(ctx, dto)
	if err != nil {
		return 0, common.TranslatePostgresError(err, s.logger)
	}

	return id, nil
}
