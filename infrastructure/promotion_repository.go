package infrastructure

import (
    "context"
    "errors"
    "fmt"
    "strconv"

    "github.com/execut/ozon-promotions/domain"
)

type PromotionRepository struct {
    client *PromotionClient
}

func NewPromotionRepository(configPath string) *PromotionRepository {
    config := readConfig(configPath)
    client := NewPromotionClient(config.CompanyID, config.Cookie)
    return &PromotionRepository{
        client: client,
    }
}

func (r *PromotionRepository) Load(ctx context.Context, id domain.PromotionID) (*domain.Promotion, error) {
    productListResponse, err := r.client.PromotionProductList(id.ID())
    if err != nil {
        return nil, err
    }

    productCandidateListResponse, err := r.client.ProductCandidateList(id.ID())
    if err != nil {
        return nil, err
    }

    productCandidateList := make([]domain.PromotionProductCandidate, 0, len(productCandidateListResponse)+len(productListResponse))

    productList := make([]domain.PromotionProduct, 0, len(productListResponse))
    for _, productResponse := range productListResponse {
        productIDValue, err := strconv.ParseInt(productResponse.Id, 10, 64)
        if err != nil {
            return nil, err
        }

        productID, err := domain.NewProductID(productIDValue)
        if err != nil {
            return nil, err
        }

        price, err := domain.NewPrice(float64(productResponse.ActionPrice))
        if err != nil {
            return nil, err
        }

        product, err := domain.NewPromotionProduct(productID, price)
        if err != nil {
            return nil, err
        }

        productList = append(productList, product)

        boostingX2PriceValue := float64(productResponse.PriceReferenceForBoosting.BoostingX2Price)
        boostingX2Price, err := domain.NewPrice(boostingX2PriceValue)
        if err != nil {
            return nil, err
        }

        boostingMaxPriceValue := float64(productResponse.PriceReferenceForBoosting.MaxBoostingPrice)
        boostingMaxPrice, err := domain.NewPrice(boostingMaxPriceValue)
        if err != nil {
            return nil, err
        }

        candidate, err := domain.NewPromotionProductCandidate(productID, boostingX2Price, boostingMaxPrice)
        if err != nil {
            return nil, err
        }

        productCandidateList = append(productCandidateList, candidate)
    }

    for _, productResponse := range productCandidateListResponse {
        productIDValue, err := strconv.ParseInt(productResponse.Id, 10, 64)
        if err != nil {
            panic(err)
        }

        productID, err := domain.NewProductID(productIDValue)
        if err != nil {
            return nil, err
        }

        priceValue := float64(productResponse.PriceReferenceForBoosting.BoostingX2Price)
        price, err := domain.NewPrice(priceValue)
        if err != nil {
            return nil, err
        }

        boostingMaxPriceValue := float64(productResponse.PriceReferenceForBoosting.MaxBoostingPrice)
        boostingMaxPrice, err := domain.NewPrice(boostingMaxPriceValue)
        if err != nil {
            return nil, err
        }

        product, err := domain.NewPromotionProductCandidate(productID, price, boostingMaxPrice)
        if err != nil {
            return nil, err
        }

        productCandidateList = append(productCandidateList, product)
    }

    promotion, err := domain.NewPromotion(id, productList, productCandidateList)
    if err != nil {
        return nil, err
    }

    return promotion, nil
}

func (r *PromotionRepository) Save(ctx context.Context, promotion *domain.Promotion) error {
    var requestProducts []ActivateProductListRequestProduct
    var deactivateProductList []string

    for _, event := range promotion.EventList() {
        switch e := event.(type) {
        case domain.PromotionProductPromotedEvent:
            fmt.Printf("PromotionProductPromotedEvent sku: %v, price: %v\n", e.Product.SKU().Value(), e.NewPrice.Value())
            requestProduct := ActivateProductListRequestProduct{
                Id:              strconv.FormatInt(e.Product.ID().ID(), 10),
                ActionPrice:     int(e.NewPrice.Value()),
                DiscountPercent: 0,
                Quantity:        0,
                Currency:        "RUB",
                Sku:             strconv.FormatInt(e.Product.SKU().Value(), 10),
            }
            requestProducts = append(requestProducts, requestProduct)
        case domain.PromotionProductStoppedEvent:
            fmt.Printf("PromotionProductStoppedEvent sku: %v\n", e.Product.SKU().Value())
            deactivateProductList = append(deactivateProductList, strconv.FormatInt(e.Product.ID().ID(), 10))
            break
        case domain.PromotionProductPromotionPriceUpdatedEvent:
            fmt.Printf("PromotionProductPromotionPriceUpdatedEvent sku: %v, price: %v\n", e.Product.SKU().Value(), e.NewPrice.Value())
            requestProduct := ActivateProductListRequestProduct{
                Id:              strconv.FormatInt(e.Product.ID().ID(), 10),
                ActionPrice:     int(e.NewPrice.Value()),
                DiscountPercent: 0,
                Quantity:        0,
                Currency:        "RUB",
                Sku:             strconv.FormatInt(e.Product.SKU().Value(), 10),
            }
            requestProducts = append(requestProducts, requestProduct)
            break
        }
    }

    if len(requestProducts) != 0 {
        request := ActivateProductListRequest{
            Products: requestProducts,
        }

        response, err := r.client.ActivateProductList(promotion.ID().ID(), request)
        if err != nil {
            return err
        }

        if len(response.ProductIds) != len(request.Products) {
            return errors.New(fmt.Sprintf("product activation count mismatch. Request products: %v, response products: %v", requestProducts, response.ProductIds))
        }
    }

    if len(deactivateProductList) != 0 {
        request := DeactivateProductListRequest{
            ProductIds: deactivateProductList,
        }

        response, err := r.client.DeactivateProductList(promotion.ID().ID(), request)
        if err != nil {
            return err
        }

        if len(response.ProductIds) != len(deactivateProductList) {
            return errors.New("product deactivation count mismatch")
        }
    }

    return nil
}
