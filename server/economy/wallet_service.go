package economy

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"go.uber.org/zap"
)

type WalletService interface {
	CreateWalletOfNewPlayer(ctx context.Context, dto CreateNewPlayerWalletDTO) error
	GetWalletOfPlayer(ctx context.Context, dto GetWalletOfPlayerDTO) ([]WalletDTO, error)
	UpdateWalletOfPlayer(ctx context.Context, dto UpdateWalletOfPlayerDTO) ([]WalletDTO, error)
	DeleteWalletOfPlayer(ctx context.Context, dto DeleteWalletOfPlayerDTO) error
}

type LocalWalletService struct {
	repo   WalletRepository
	config *config.Config
	logger *zap.Logger
}

func NewLocalWalletService(repository WalletRepository, c *config.Config, l *zap.Logger) *LocalWalletService {
	return &LocalWalletService{
		repo:   repository,
		config: c,
		logger: l,
	}
}

func (s *LocalWalletService) CreateWalletOfNewPlayer(ctx context.Context, dto CreateNewPlayerWalletDTO) error {
	ctx, span := tracer.Start(ctx, "economy.wallet.createwalletofnewplayer")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.CreateWalletOfNewPlayer(ctx, dto),
		s.logger,
	).WithMessage("Failed to create wallet")
}

func (s *LocalWalletService) GetWalletOfPlayer(ctx context.Context, dto GetWalletOfPlayerDTO) ([]WalletDTO, error) {
	ctx, span := tracer.Start(ctx, "economy.wallet.getwalletofplayer")
	defer span.End()

	wallet, err := s.repo.GetWalletOfPlayer(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to get wallet")
	}

	return wallet, nil
}

func (s *LocalWalletService) UpdateWalletOfPlayer(ctx context.Context, dto UpdateWalletOfPlayerDTO) ([]WalletDTO, error) {
	ctx, span := tracer.Start(ctx, "economy.wallet.updatewalletofplayer")
	defer span.End()

	wallet, err := s.repo.UpdateWalletOfPlayer(ctx, dto)
	if err != nil {
		return nil, common.TranslatePostgresError(err, s.logger).WithMessage("Failed to update wallet")
	}

	return wallet, nil
}

func (s *LocalWalletService) DeleteWalletOfPlayer(ctx context.Context, dto DeleteWalletOfPlayerDTO) error {
	ctx, span := tracer.Start(ctx, "economy.wallet.deletewalletofplayer")
	defer span.End()

	return common.TranslatePostgresError(
		s.repo.DeleteWalletOfPlayer(ctx, dto),
		s.logger,
	).WithMessage("Failed to delete wallet")
}
