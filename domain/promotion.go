package domain

import "errors"

var (
    ErrDiscountLessThanMinimalPrice      = errors.New("promotion less than minimal price")
    ErrProductAlreadyPromoted            = errors.New("product already promoted")
    ErrPromotionForProductAlreadyStopped = errors.New("promotion for product already stopped")
    ErrProductNotExists                  = errors.New("product does not exist")
    ErrPromotionCandidateNotFound        = errors.New("promotion candidate not found")
)

func NewPromotion(id PromotionID, productList []PromotionProduct, candidateProductList []PromotionProductCandidate) (*Promotion, error) {
    return &Promotion{
        id:                   id,
        productList:          productList,
        candidateProductList: candidateProductList,
    }, nil
}

type Promotion struct {
    id                   PromotionID
    productList          []PromotionProduct
    candidateProductList []PromotionProductCandidate
    eventList            []PromotionEvent
}

func (p *Promotion) ID() PromotionID {
    return p.id
}

func (p *Promotion) EventList() []PromotionEvent {
    return p.eventList
}

func (p *Promotion) ProductList() []PromotionProduct {
    return p.productList
}

func (p *Promotion) Product(product Product) (PromotionProduct, error) {
    for _, pr := range p.productList {
        if pr.ID().Same(product.ID()) {
            return pr, nil
        }
    }

    return PromotionProduct{}, ErrProductNotExists
}

func (p *Promotion) CalculateDiscount(product Product, criteria PromotionCriteria) (Price, error) {
    productCandidate := p.findCandidate(product)
    if productCandidate == nil {
        return Price{}, ErrPromotionCandidateNotFound
    }

    var targetPrice Price
    switch criteria {
    case PromotionCriteriaMinimalPrice:
        targetPrice = productCandidate.PromotionMinPrice()
    case PromotionCriteriaMaximalBoosting:
        targetPrice = productCandidate.PromotionMaxPrice()
        if targetPrice.LessOrEqual(product.MinPrice()) {
            if productCandidate.PromotionMinPrice().LessOrEqual(product.MinPrice()) {
                return Price{}, ErrDiscountLessThanMinimalPrice
            }

            targetPrice = product.MinPrice()
        }
    default:
        return Price{}, errors.New("unsupported criteria for CalculateDiscount")
    }

    discount, err := product.Price().Diff(targetPrice)
    if err != nil {
        return Price{}, err
    }

    if !product.DiscountIsPossible(discount) {
        return Price{}, ErrDiscountLessThanMinimalPrice
    }

    return discount, nil
}

func (p *Promotion) Same(promotion Promotion) bool {
    return p.id == promotion.id
}

func (p *Promotion) Promote(product Product, criteria PromotionCriteria) error {
    discount, err := p.CalculateDiscount(product, criteria)
    if err != nil {
        return err
    }

    newPrice, err := product.Price().Reduce(discount)
    if err != nil {
        return err
    }

    if !p.HasProduct(product) {
        p.addEvent(PromotionProductPromotedEvent{
            Product:  product,
            NewPrice: newPrice,
        })

        return nil
    }

    promotionProduct, err := p.Product(product)
    if err != nil {
        return err
    }

    if promotionProduct.PromotionPrice().Same(newPrice) {
        return nil
    }

    p.addEvent(PromotionProductPromotionPriceUpdatedEvent{
        Product:  product,
        NewPrice: newPrice,
    })

    return nil
}

func (p *Promotion) addEvent(event PromotionEvent) {
    p.eventList = append(p.eventList, event)
}

func (p *Promotion) StopPromoting(product Product) error {
    if !p.HasProduct(product) {
        return ErrPromotionForProductAlreadyStopped
    }

    p.addEvent(PromotionProductStoppedEvent{
        Product: product,
    })

    return nil
}

func (p *Promotion) ApplyEvents() error {
    for _, event := range p.eventList {
        switch typedEvent := event.(type) {
        case PromotionProductPromotedEvent:
            product := typedEvent.Product
            promotionProduct, err := NewPromotionProduct(product.ID(), typedEvent.NewPrice)
            if err != nil {
                return err
            }

            p.productList = append(p.productList, promotionProduct)
        case PromotionProductPromotionPriceUpdatedEvent:
            for key, promotionProduct := range p.productList {
                if promotionProduct.ID().Same(typedEvent.Product.ID()) {
                    product := typedEvent.Product
                    currentPromotionProduct, err := NewPromotionProduct(product.ID(), typedEvent.NewPrice)
                    if err != nil {
                        return err
                    }

                    p.productList[key] = currentPromotionProduct
                }
            }
        case PromotionProductStoppedEvent:
            newProductList := make([]PromotionProduct, 0, len(p.productList))
            for _, promotionProduct := range p.productList {
                if !promotionProduct.ID().Same(typedEvent.Product.ID()) {
                    newProductList = append(newProductList, promotionProduct)
                }
            }

            p.productList = newProductList
        default:
            return errors.New("unsupported event type")
        }
    }
    p.eventList = nil
    return nil
}

func (p *Promotion) findCandidate(product Product) *PromotionProductCandidate {
    for _, promotionProduct := range p.candidateProductList {
        if promotionProduct.ID().Same(product.ID()) {
            return &promotionProduct
        }
    }

    return nil
}

func (p *Promotion) HasProduct(product Product) bool {
    for _, promotionProduct := range p.productList {
        if promotionProduct.ID().Same(product.ID()) {
            return true
        }
    }

    return false
}
