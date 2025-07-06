package contract

import "context"

type Application interface {
    Commands
    Processors
}

type Commands interface {
    PromotionActualizeForGroup(ctx context.Context, cmd PromotionActualizeForGroup) error
}

type PromotionActualizeForGroup struct {
    PromotionID int64
    GroupID     string
}

type Processors interface {
    PromotionActualize(ctx context.Context, cmd PromotionActualize) error
}

type PromotionActualize struct {
    PromotionID int64
}
