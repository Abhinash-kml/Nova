package economy

import "context"

type WalletRepository interface {
	CreateWalletOfNewPlayer(ctx context.Context, dto CreateNewPlayerWalletRequest) error
	GetWalletOfPlayer(ctx context.Context, dto GetWalletOfPlayerRequest) ([]WalletCurrency, error)
	UpdateWalletOfPlayer(ctx context.Context, dto UpdateWalletOfPlayerRequest) ([]WalletCurrency, error)
	DeleteWalletOfPlayer(ctx context.Context, dto DeleteWalletOfPlayerRequest) error
}
