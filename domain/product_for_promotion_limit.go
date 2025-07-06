package domain

import "errors"

func NewProductForPromotionLimit(limit uint8) (ProductForPromotionLimit, error) {
    if limit == 0 {
        return ProductForPromotionLimit{}, errors.New("limit cannot be zero")
    }
    return ProductForPromotionLimit{
        limit: limit,
    }, nil
}

type ProductForPromotionLimit struct {
    limit uint8
}

func (p ProductForPromotionLimit) Limit() uint8 {
    return p.limit
}
