package progression

import (
	"context"
)

type ProgressionRepository interface {
	Update(ctx context.Context, dto UpdateDTO) (Progression, error)
}
