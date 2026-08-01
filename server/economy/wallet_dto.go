package economy

import "github.com/google/uuid"

type WalletUserID struct {
	UserID string `form:"userid" binding:"required,uuid"`
}

type WalletCurrencyID struct {
	CurrencyID string `form:"currencyid" binding:"required,numeric"`
}

type CreateNewPlayerWalletDTO struct {
	WalletUserID
}

type GetWalletOfPlayerDTO struct {
	WalletUserID
}

type WalletDTO struct {
	CurrencyID   int `json:"currency_id"`
	BalanceMinor int `json:"balance_minor"`
}

type UpdateWalletOfPlayerDTO struct {
	UserID     uuid.UUID `json:"user_id" binding:"required"`
	CurrencyID int       `json:"currency_id" binding:"required"`
	Operation  string    `json:"operation" binding:"required,oneof=credit debit"`
	Amount     int       `json:"amount" binding:"required"`
}

type DeleteWalletOfPlayerDTO struct {
	WalletUserID
}
