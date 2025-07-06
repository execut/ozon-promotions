package domain

import "context"

type ProductGroupRepository interface {
    LoadList(ctx context.Context) ([]*ProductGroup, error)
    Load(ctx context.Context, productGroup *ProductGroup) error
}
