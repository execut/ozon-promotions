package main

import (
    "context"
    "path/filepath"

    "github.com/execut/ozon-promotions/application"
    "github.com/execut/ozon-promotions/application/contract"
    "github.com/execut/ozon-promotions/infrastructure"
)

func main() {
    ctx := context.Background()
    configPath, err := filepath.Abs("config.yml")
    if err != nil {
        panic(err)
    }

    productGroupRepository := infrastructure.NewProductGroupRepository(configPath)
    promotionRepository := infrastructure.NewPromotionRepository(configPath)

    app := application.NewApplication(productGroupRepository, promotionRepository)
    var promotionID int64 = 1977747

    err = app.PromotionActualize(ctx, contract.PromotionActualize{
        PromotionID: promotionID,
    })

    if err != nil {
        panic(err)
    }
}
