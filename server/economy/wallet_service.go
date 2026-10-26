package economy

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"go.uber.org/zap"
)

type WalletService interface {
	CreateWalletOfNewPlayer(ctx context.Context, dto CreateNewPlayerWalletRequest) error
	GetWalletOfPlayer(ctx context.Context, dto GetWalletOfPlayerRequest) ([]WalletCurrency, error)
	UpdateWalletOfPlayer(ctx context.Context, dto UpdateWalletOfPlayerRequest) ([]WalletCurrency, error)
	DeleteWalletOfPlayer(ctx context.Context, dto DeleteWalletOfPlayerRequest) error
}

type LocalWalletService struct {
	repo   WalletRepository
	config *config.WalletConfig
	logger *zap.Logger
}

func NewLocalWalletService(repository WalletRepository, c *config.WalletConfig, l *zap.Logger) *LocalWalletService {
	return &LocalWalletService{
		repo:   repository,
		config: c,
		logger: l,
	}
}

func (s *LocalWalletService) CreateWalletOfNewPlayer(ctx context.Context, dto CreateNewPlayerWalletRequest) error {
	ctx, span := tracer.Start(ctx, "economy.wallet.CreateWalletOfNewPlayer")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.CreateWalletOfNewPlayer(ctx, dto),
		s.logger,
	).WithMessage("Failed to create wallet")
}

func (s *LocalWalletService) GetWalletOfPlayer(ctx context.Context, dto GetWalletOfPlayerRequest) ([]WalletCurrency, error) {
	ctx, span := tracer.Start(ctx, "economy.wallet.GetWalletOfPlayer")
	defer span.End()

	wallet, err := s.repo.GetWalletOfPlayer(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get wallet")
	}

	return wallet, nil
}

func (s *LocalWalletService) UpdateWalletOfPlayer(ctx context.Context, dto UpdateWalletOfPlayerRequest) ([]WalletCurrency, error) {
	ctx, span := tracer.Start(ctx, "economy.wallet.UpdateWalletOfPlayer")
	defer span.End()

	wallet, err := s.repo.UpdateWalletOfPlayer(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update wallet")
	}

	return wallet, nil
}

func (s *LocalWalletService) DeleteWalletOfPlayer(ctx context.Context, dto DeleteWalletOfPlayerRequest) error {
	ctx, span := tracer.Start(ctx, "economy.wallet.DeleteWalletOfPlayer")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.DeleteWalletOfPlayer(ctx, dto),
		s.logger,
	).WithMessage("Failed to delete wallet")
}
