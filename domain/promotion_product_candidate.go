package domain

func NewPromotionProductCandidate(id ProductID, promotionMinPrice Price, promotionMaxPrice Price) (PromotionProductCandidate, error) {
    return PromotionProductCandidate{
        id:                id,
        promotionMinPrice: promotionMinPrice,
        promotionMaxPrice: promotionMaxPrice,
    }, nil
}

type PromotionProductCandidate struct {
    id                ProductID
    promotionMinPrice Price
    promotionMaxPrice Price
}

func (p PromotionProductCandidate) PromotionMaxPrice() Price {
    return p.promotionMaxPrice
}

func (p PromotionProductCandidate) PromotionMinPrice() Price {
    return p.promotionMinPrice
}

func (p PromotionProductCandidate) ID() ProductID {
    return p.id
}
