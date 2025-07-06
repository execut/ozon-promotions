package domain

type PromotionEvent interface{}

type PromotionProductPromotedEvent struct {
    Product  Product
    NewPrice Price
}

type PromotionProductPromotionPriceUpdatedEvent struct {
    Product  Product
    NewPrice Price
}

type PromotionProductStoppedEvent struct {
    Product Product
}
