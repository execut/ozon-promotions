package domain

func NewDiscount(priceValue float64) (Price, error) {
    return Price{
        value: priceValue,
    }, nil
}
