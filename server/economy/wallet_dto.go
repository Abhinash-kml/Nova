package economy

import "github.com/google/uuid"

type WalletUserID struct {
	UserID string `form:"userid" binding:"required,uuid"`
}

type WalletCurrencyID struct {
	CurrencyID string `form:"currencyid" binding:"required,numeric"`
}

type CreateNewPlayerWalletRequest struct {
	WalletUserID
}

type GetWalletOfPlayerRequest struct {
	WalletUserID
}

type WalletCurrency struct {
	CurrencyID   int `json:"currency_id"`
	BalanceMinor int `json:"balance_minor"`
}

type UpdateWalletOfPlayerRequest struct {
	UserID     uuid.UUID `json:"user_id" binding:"required"`
	CurrencyID int       `json:"currency_id" binding:"required"`
	Operation  string    `json:"operation" binding:"required,oneof=credit debit"`
	Amount     int       `json:"amount" binding:"required"`
}

type DeleteWalletOfPlayerRequest struct {
	WalletUserID
}
