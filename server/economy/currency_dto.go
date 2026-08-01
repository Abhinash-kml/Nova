package economy

type CurrencyID struct {
	ID string `uri:"id" binding:"required,numeric"`
}
type GetCurrencyDTO struct {
	CurrencyID
}
type GetAllCurrencyDTO struct {
	Limit  int `form:"limit,default=10"`
	Cursor int `form:"cursor,default=0"`
}

type CreateCurrencyDTO struct {
	Name              string `json:"name" binding:"required"`
	Code              string `json:"code" binding:"required"`
	Kind              string `json:"kind" binding:"required"`
	IsPurchasable     bool   `json:"is_purchasable" binding:"required"`
	IsTradable        bool   `json:"is_tradable" binding:"required"`
	ExipiresAfterDays int    `json:"expires_after_days" binding:"required"`
	MaxBalanceMinor   int    `json:"max_balance_minor" binding:"required"`
	SortOrder         int    `json:"sort_order" binding:"required"`
}

type UpdateData struct {
	Name              string `json:"name"`
	Code              string `json:"code"`
	Kind              string `json:"kind"`
	IsPurchasable     *bool  `json:"is_purchasable"`
	IsTradable        *bool  `json:"is_tradable"`
	ExipiresAfterDays *int   `json:"expires_after_days"`
	MaxBalanceMinor   *int   `json:"max_balance_minor"`
	SortOrder         *int   `json:"sort_order"`
}

type UpdateCurrencyDTO struct {
	CurrencyID
	UpdateData
}

type DeleetCurrencyDTO struct {
	CurrencyID
}

type CurrencyDTO struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Code              string `json:"code"`
	Kind              string `json:"kind"`
	IsPurchasable     bool   `json:"is_purchasable"`
	IsTradable        bool   `json:"is_tradable"`
	ExipiresAfterDays int    `json:"expires_after_days"`
	MaxBalanceMinor   int    `json:"max_balance_minor"`
	SortOrder         int    `json:"sort_order"`
}
