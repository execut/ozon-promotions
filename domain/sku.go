package domain

import (
    "errors"
)

type SKU struct {
    sku int64
}

func NewSKU(sku int64) (SKU, error) {
    if sku < 1 {
        return SKU{}, errors.New("SKU must be greater than 0")
    }

    return SKU{sku: sku}, nil
}

func (s SKU) Value() int64 {
    return s.sku
}

func (s SKU) Equals(other SKU) bool {
    return s.sku == other.sku
}
