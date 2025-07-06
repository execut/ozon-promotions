package domain_test

import "github.com/execut/ozon-promotions/domain"

var (
    expectedMinPrice                            = NewPrice(100.0)
    expectedPrice                               = NewPrice(200.0)
    expectedProductIDValue                int64 = 1
    expectedPromotionIDValue              int64 = 2
    expectedPromotionID, _                      = domain.NewPromotionID(expectedPromotionIDValue)
    expectedProductID, _                        = domain.NewProductID(expectedProductIDValue)
    expectedSKU, _                              = domain.NewSKU(3)
    expectedProduct, _                          = domain.NewProduct(expectedProductID, expectedSKU, expectedMinPrice, expectedPrice)
    expectedProductGroupID                      = "test-product-group-id"
    expectedProductForPromotionLimitValue uint8 = 2
    expectedProductForPromotionLimit, _         = domain.NewProductForPromotionLimit(expectedProductForPromotionLimitValue)
    expectedPromotionProduct, _                 = domain.NewPromotionProduct(expectedProductID, expectedMinPrice)
    expectedPromotion                           = NewPromotion(expectedPromotionID, []domain.PromotionProduct{expectedPromotionProduct}, nil)
)

func NewProductID(productID int64) domain.ProductID {
    id, _ := domain.NewProductID(productID)

    return id
}

func NewProduct(productID domain.ProductID, sku domain.SKU, minPrice domain.Price, price domain.Price) domain.Product {
    product, _ := domain.NewProduct(productID, sku, minPrice, price)

    return product
}

func NewPrice(priceValue float64) domain.Price {
    price, _ := domain.NewPrice(priceValue)

    return price
}

func NewPromotion(id domain.PromotionID, productList []domain.PromotionProduct, candidateProductList []domain.PromotionProductCandidate) *domain.Promotion {
    promotion, _ := domain.NewPromotion(id, productList, candidateProductList)
    return promotion
}

func NewPromotionProductCandidate(id domain.ProductID, promotionMinPrice domain.Price, promotionMaxPrice domain.Price) domain.PromotionProductCandidate {
    candidate, _ := domain.NewPromotionProductCandidate(id, promotionMinPrice, promotionMaxPrice)

    return candidate
}

func NewPromotionProduct(id domain.ProductID, promotionPrice domain.Price) domain.PromotionProduct {
    product, _ := domain.NewPromotionProduct(id, promotionPrice)

    return product
}
