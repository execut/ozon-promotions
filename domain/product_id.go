package domain

func NewProductID(id int64) (ProductID, error) {
    return ProductID{id: id}, nil
}

type ProductID struct {
    id int64
}

func (p ProductID) ID() int64 {
    return p.id
}

func (p ProductID) Same(otherProductID ProductID) bool {
    return p.id == otherProductID.ID()
}
