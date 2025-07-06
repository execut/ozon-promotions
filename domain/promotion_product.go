package domain

func NewPromotionProduct(id ProductID, promotionPrice Price) (PromotionProduct, error) {
    return PromotionProduct{
        id:             id,
        promotionPrice: promotionPrice,
    }, nil
}

type PromotionProduct struct {
    id             ProductID
    promotionPrice Price
}

func (p PromotionProduct) PromotionPrice() Price {
    return p.promotionPrice
}

func (p PromotionProduct) ID() ProductID {
    return p.id
}

func (p PromotionProduct) Same(another PromotionProduct) bool {
    return p.id.Same(another.id) && p.promotionPrice.Same(another.promotionPrice)
}

func (p PromotionProduct) Replace(another PromotionProduct) (PromotionProduct, error) {
    return NewPromotionProduct(p.id, another.promotionPrice)
}
