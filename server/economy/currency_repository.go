package economy

import (
	"context"
)

type CurrencyRepository interface {
	Get(ctx context.Context, req GetCurrencyRequest) (Currency, error)
	GetAll(ctx context.Context, req GetAllCurrencyRequest) ([]Currency, error)
	Create(ctx context.Context, req CreateCurrencyRequest) (Currency, error)
	Update(ctx context.Context, req UpdateCurrencyRequest) (Currency, error)
	Delete(ctx context.Context, req DeleteCurrencyRequest) (int, error)
}
