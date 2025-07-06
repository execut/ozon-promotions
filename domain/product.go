package domain

func NewProduct(id ProductID, sku SKU, minPrice Price, price Price) (Product, error) {
    return Product{id: id, sku: sku, minPrice: minPrice, price: price}, nil
}

type Product struct {
    id       ProductID
    sku      SKU
    minPrice Price
    price    Price
}

func (p Product) MinPrice() Price {
    return p.minPrice
}

func (p Product) SKU() SKU {
    return p.sku
}

func (p Product) Price() Price {
    return p.price
}

func (p Product) Same(otherProduct Product) bool {
    return otherProduct.ID().Same(otherProduct.ID())
}

func (p Product) DiscountIsPossible(discount Price) bool {
    newPrice, err := NewPrice(p.price.Value() - discount.Value())
    if err != nil {
        return false
    }

    return p.minPrice.LessOrEqual(newPrice)
}

func (p Product) ID() ProductID {
    return p.id
}
