package domain

func NewPromotionID(id int64) (PromotionID, error) {
    return PromotionID{id: id}, nil
}

type PromotionID struct {
    id int64
}

func (p PromotionID) ID() int64 {
    return p.id
}

func (p PromotionID) Same(otherPromotionID PromotionID) bool {
    return p.id == otherPromotionID.ID()
}
