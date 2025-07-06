package domain

func NewPrice(priceValue float64) (Price, error) {
    return Price{
        value: priceValue,
    }, nil
}

type Price struct {
    value float64
}

func (p Price) Value() float64 {
    return p.value
}

func (p Price) LessOrEqual(price Price) bool {
    return p.value <= price.value
}

func (p Price) Diff(price Price) (Price, error) {
    return NewDiscount(p.value - price.value)
}

func (p Price) Reduce(price Price) (Price, error) {
    return NewPrice(p.value - price.value)
}

func (p Price) Same(price Price) bool {
    return p.value == price.Value()
}
