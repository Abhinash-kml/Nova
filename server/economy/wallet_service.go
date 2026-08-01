package economy

import (
	"context"

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

	return s.repo.CreateWalletOfNewPlayer(ctx, dto)
}

func (s *LocalWalletService) GetWalletOfPlayer(ctx context.Context, dto GetWalletOfPlayerDTO) ([]WalletDTO, error) {
	ctx, span := tracer.Start(ctx, "economy.wallet.getwalletofplayer")
	defer span.End()

	return s.repo.GetWalletOfPlayer(ctx, dto)
}

func (s *LocalWalletService) UpdateWalletOfPlayer(ctx context.Context, dto UpdateWalletOfPlayerDTO) ([]WalletDTO, error) {
	ctx, span := tracer.Start(ctx, "economy.wallet.updatewalletofplayer")
	defer span.End()

	return s.repo.UpdateWalletOfPlayer(ctx, dto)
}

func (s *LocalWalletService) DeleteWalletOfPlayer(ctx context.Context, dto DeleteWalletOfPlayerDTO) error {
	ctx, span := tracer.Start(ctx, "economy.wallet.deletewalletofplayer")
	defer span.End()

	return s.repo.DeleteWalletOfPlayer(ctx, dto)
}
