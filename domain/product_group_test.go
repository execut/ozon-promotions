package domain_test

import (
    "testing"

    "github.com/execut/ozon-promotions/domain"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestProductGroup_ActualizeProductListPromotion(t *testing.T) {
    t.Parallel()
    productWithLessDiscountID := NewProductID(2)
    productWithLessDiscount := NewProduct(productWithLessDiscountID, expectedSKU, expectedMinPrice, expectedPrice)
    productWithLess2DiscountID := NewProductID(3)
    productWithLess2Discount := NewProduct(productWithLess2DiscountID, expectedSKU, expectedMinPrice, expectedPrice)
    productList := []domain.Product{
        expectedProduct,
        productWithLessDiscount,
        productWithLess2Discount,
    }
    group, _ := domain.NewTestProductGroup(t, expectedProductGroupID, productList, domain.PromotionCriteriaMaximalBoosting, expectedProductForPromotionLimit)
    promotion := NewPromotion(expectedPromotionID, nil, []domain.PromotionProductCandidate{
        NewPromotionProductCandidate(expectedProductID, NewPrice(198), NewPrice(200.0)),
        NewPromotionProductCandidate(productWithLessDiscountID, NewPrice(197), NewPrice(200.0)),
        NewPromotionProductCandidate(productWithLess2DiscountID, NewPrice(196), NewPrice(200.0)),
    })
    t.Run("when_has_other_product_and_add_with_less_price", func(t *testing.T) {
        err := group.ActualizeProductListPromotion(promotion)

        require.NoError(t, err)
    })

    t.Run("then_replace_old_product", func(t *testing.T) {
        err := group.ApplyEvents()

        require.NoError(t, err)
        assert.Len(t, promotion.ProductList(), 2)
    })
}
