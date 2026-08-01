package economy

import (
	"context"
)

type CurrencyRepository interface {
	Get(ctx context.Context, dto GetCurrencyDTO) (CurrencyDTO, error)
	GetAll(ctx context.Context, dto GetAllCurrencyDTO) ([]CurrencyDTO, error)
	Create(ctx context.Context, dto CreateCurrencyDTO) (CurrencyDTO, error)
	Update(ctx context.Context, dto UpdateCurrencyDTO) (CurrencyDTO, error)
	Delete(ctx context.Context, dto DeleetCurrencyDTO) (int, error)
}
