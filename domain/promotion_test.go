package domain_test

import (
    "fmt"
    "testing"

    "github.com/execut/ozon-promotions/domain"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPromotion_CalculateDiscount(t *testing.T) {
    t.Parallel()
    t.Run("when_criteria_minimal_price_then_minimal_promotion_price", func(t *testing.T) {
        t.Parallel()
        promotion := NewPromotion(expectedPromotionID, nil, []domain.PromotionProductCandidate{
            NewPromotionProductCandidate(expectedProductID, NewPrice(170.0), NewPrice(200.0)),
        })

        discount, err := promotion.CalculateDiscount(expectedProduct, domain.PromotionCriteriaMinimalPrice)

        require.NoError(t, err)
        assert.True(t, discount.Same(NewPrice(30.0)))
    })
    t.Run("when_criteria_maximal_boosting_then_maximal_promotion_price", func(t *testing.T) {
        t.Parallel()
        promotion := NewPromotion(expectedPromotionID, nil, []domain.PromotionProductCandidate{
            NewPromotionProductCandidate(expectedProductID, NewPrice(170.0), NewPrice(110.0)),
        })

        discount, err := promotion.CalculateDiscount(expectedProduct, domain.PromotionCriteriaMaximalBoosting)

        require.NoError(t, err)
        assert.True(t, discount.Same(NewPrice(90.0)))
    })
    t.Run("when_criteria_maximal_boosting_and_max_busting_price_more_min_price_then_return_min_price", func(t *testing.T) {
        t.Parallel()
        promotion := NewPromotion(expectedPromotionID, nil, []domain.PromotionProductCandidate{
            NewPromotionProductCandidate(expectedProductID, NewPrice(170.0), NewPrice(90.0)),
        })

        discount, err := promotion.CalculateDiscount(expectedProduct, domain.PromotionCriteriaMaximalBoosting)

        require.NoError(t, err)
        assert.True(t, discount.Same(NewPrice(100.0)), "discount should be 110.0, but got %v", discount.Value())
    })
    t.Run("when_less_than_minimal_price_then_error", func(t *testing.T) {
        t.Parallel()
        promotion := NewPromotion(expectedPromotionID, nil, []domain.PromotionProductCandidate{
            NewPromotionProductCandidate(expectedProductID, NewPrice(89.0), NewPrice(90.0)),
        })

        _, err := promotion.CalculateDiscount(expectedProduct, domain.PromotionCriteriaMinimalPrice)

        require.ErrorIs(t, err, domain.ErrDiscountLessThanMinimalPrice)
    })
    t.Run("when_min_discount_price_less_than_min_price_then_err", func(t *testing.T) {
        t.Parallel()
        promotion := NewPromotion(expectedPromotionID, nil, []domain.PromotionProductCandidate{
            NewPromotionProductCandidate(expectedProductID, NewPrice(95.0), NewPrice(90.0)),
        })

        _, err := promotion.CalculateDiscount(expectedProduct, domain.PromotionCriteriaMaximalBoosting)

        require.ErrorIs(t, err, domain.ErrDiscountLessThanMinimalPrice)
    })
}

func TestPromotion_Live(t *testing.T) {
    t.Parallel()
    expectedProductIDForUpdate := NewProductID(expectedProductIDValue + 1)
    expectedProductForUpdate := NewProduct(expectedProductIDForUpdate, expectedSKU, NewPrice(150), NewPrice(200))
    expectedNewPrice := NewPrice(169.0)
    promotion := NewPromotion(expectedPromotionID, []domain.PromotionProduct{
        NewPromotionProduct(expectedProductIDForUpdate, NewPrice(200.0)),
    }, []domain.PromotionProductCandidate{
        NewPromotionProductCandidate(expectedProductID, NewPrice(170.0), NewPrice(200.0)),
        NewPromotionProductCandidate(expectedProductIDForUpdate, expectedNewPrice, NewPrice(200.0)),
    })
    t.Run("when_promote", func(t *testing.T) {
        err := promotion.Promote(expectedProduct, domain.PromotionCriteriaMinimalPrice)

        require.NoError(t, err)
    })

    t.Run("then_has_new_product", func(t *testing.T) {
        err := promotion.ApplyEvents()
        require.NoError(t, err)

        assert.True(t, promotion.HasProduct(expectedProduct))
    })

    t.Run("when_promote_already_existed_with_other_price", func(t *testing.T) {
        err := promotion.Promote(expectedProductForUpdate, domain.PromotionCriteriaMinimalPrice)

        require.NoError(t, err)
    })

    t.Run("then_update_price", func(t *testing.T) {
        err := promotion.ApplyEvents()

        require.NoError(t, err)
        product, err := promotion.Product(expectedProductForUpdate)

        require.NoError(t, err)
        assert.True(t, product.PromotionPrice().Same(expectedNewPrice), fmt.Sprintf("expected %v, got %v", expectedNewPrice, product.PromotionPrice().Value()))
    })

    t.Run("when_stop_promoting", func(t *testing.T) {
        err := promotion.StopPromoting(expectedProduct)

        require.NoError(t, err)
    })

    t.Run("then_product_not_exist", func(t *testing.T) {
        err := promotion.ApplyEvents()
        require.NoError(t, err)

        assert.False(t, promotion.HasProduct(expectedProduct))
    })

    t.Run("when_promote_already_stopped_then_err", func(t *testing.T) {
        err := promotion.StopPromoting(expectedProduct)

        require.ErrorIs(t, err, domain.ErrPromotionForProductAlreadyStopped)
    })
}
