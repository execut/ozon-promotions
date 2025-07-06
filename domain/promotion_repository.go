package domain

import (
    "context"
)

type PromotionRepository interface {
    Load(ctx context.Context, id PromotionID) (*Promotion, error)
    Save(ctx context.Context, promotion *Promotion) error
}
