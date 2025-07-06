package domain

type ProductGroupEvent interface {
}

type ProductGroupProductActualizedEvent struct {
    Promotion *Promotion
}
