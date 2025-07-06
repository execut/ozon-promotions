package application

import (
    "github.com/execut/ozon-promotions/application/commands"
    "github.com/execut/ozon-promotions/application/contract"
    "github.com/execut/ozon-promotions/application/processors"
    "github.com/execut/ozon-promotions/domain"
)

var _ contract.Application = (*App)(nil)

type App struct {
    appCommands
    appProcessors
}

type appCommands struct {
    commands.PromotionActualizeForGroupHandler
}

type appProcessors struct {
    processors.PromotionActualizeHandler
}

func NewApplication(productGroupRepository domain.ProductGroupRepository, promotionRepository domain.PromotionRepository) contract.Application {
    return &App{
        appCommands: appCommands{
            PromotionActualizeForGroupHandler: commands.NewPromotionActualizeForGroupHandler(productGroupRepository, promotionRepository),
        },
        appProcessors: appProcessors{
            processors.NewPromotionActualizeHandler(productGroupRepository, promotionRepository),
        },
    }
}
