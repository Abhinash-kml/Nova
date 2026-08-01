package economy

import "context"

type WalletRepository interface {
	CreateWalletOfNewPlayer(ctx context.Context, dto CreateNewPlayerWalletDTO) error
	GetWalletOfPlayer(ctx context.Context, dto GetWalletOfPlayerDTO) ([]WalletDTO, error)
	UpdateWalletOfPlayer(ctx context.Context, dto UpdateWalletOfPlayerDTO) ([]WalletDTO, error)
	DeleteWalletOfPlayer(ctx context.Context, dto DeleteWalletOfPlayerDTO) error
}
