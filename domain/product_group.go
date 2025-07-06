package domain

import (
    "fmt"
    "testing"
)

func NewProductGroup(id string, productList []Product, promotionCriteria PromotionCriteria, productForPromotionLimit ProductForPromotionLimit) (*ProductGroup, error) {
    return &ProductGroup{
        id:                       id,
        productList:              productList,
        promotionCriteria:        promotionCriteria,
        productForPromotionLimit: productForPromotionLimit,
        eventList:                []ProductGroupEvent{},
    }, nil
}

func NewTestProductGroup(_ *testing.T, id string, productList []Product, promotionCriteria PromotionCriteria, productForPromotionLimit ProductForPromotionLimit) (*ProductGroup, error) {
    productGroup, _ := NewProductGroup(id, productList, promotionCriteria, productForPromotionLimit)
    return productGroup, nil
}

type ProductGroup struct {
    id                       string
    productList              []Product
    promotionCriteria        PromotionCriteria
    productForPromotionLimit ProductForPromotionLimit
    eventList                []ProductGroupEvent
}

func (p *ProductGroup) ActualizeProductListPromotion(promotion *Promotion) error {
    bestResults := make([]BestResult, 0, len(p.productList))
    for _, product := range p.productList {
        discount, err := promotion.CalculateDiscount(product, p.promotionCriteria)
        if err != nil {
            continue
        }

        bestResults = append(bestResults, BestResult{
            Product:  product,
            Discount: discount,
        })
    }

    bestProductMap := make(map[ProductID]Product, p.productForPromotionLimit.Limit())
    for _ = range p.productForPromotionLimit.Limit() {
        var (
            bestResult    *BestResult
            bestResultKey int
        )
        for key, result := range bestResults {
            if bestResult == nil || result.Discount.LessOrEqual(bestResult.Discount) {
                bestResult = &result
                bestResultKey = key
            }
        }

        if bestResult == nil {
            break
        }

        bestResults = append(bestResults[:bestResultKey], bestResults[bestResultKey+1:]...)
        bestProductMap[bestResult.Product.ID()] = bestResult.Product
    }

    for _, product := range p.productList {
        if _, ok := bestProductMap[product.ID()]; ok {
            err := promotion.Promote(product, p.promotionCriteria)
            if err != nil {
                return err
            }

            continue
        }

        if promotion.HasProduct(product) {
            err := promotion.StopPromoting(product)
            if err != nil {
                return err
            }
        }
    }

    p.addEvent(ProductGroupProductActualizedEvent{
        Promotion: promotion,
    })

    return nil
}

func (p *ProductGroup) EventList() []ProductGroupEvent {
    return p.eventList
}

func (p *ProductGroup) ApplyEvents() error {
    for _, event := range p.eventList {
        switch eventTyped := event.(type) {
        case ProductGroupProductActualizedEvent:
            err := eventTyped.Promotion.ApplyEvents()
            if err != nil {
                return err
            }
        default:
            return fmt.Errorf("unknown event type: %T", eventTyped)
        }
    }

    return nil
}

func (p *ProductGroup) addEvent(event ProductGroupEvent) {
    p.eventList = append(p.eventList, event)
}

type BestResult struct {
    Product  Product
    Discount Price
}
