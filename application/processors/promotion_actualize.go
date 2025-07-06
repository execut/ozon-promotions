package processors

import (
    "context"

    "github.com/execut/ozon-promotions/application/contract"
    "github.com/execut/ozon-promotions/domain"
)

type PromotionActualizeHandler struct {
    repository          domain.ProductGroupRepository
    promotionRepository domain.PromotionRepository
}

func NewPromotionActualizeHandler(repository domain.ProductGroupRepository, promotionRepository domain.PromotionRepository) PromotionActualizeHandler {
    return PromotionActualizeHandler{repository: repository, promotionRepository: promotionRepository}
}

func (h *PromotionActualizeHandler) PromotionActualize(ctx context.Context, cmd contract.PromotionActualize) error {
    promotionID, err := domain.NewPromotionID(cmd.PromotionID)
    if err != nil {
        return err
    }

    productGroupList, err := h.repository.LoadList(ctx)
    if err != nil {
        return err
    }

    var promotion *domain.Promotion
    promotion, err = h.promotionRepository.Load(ctx, promotionID)
    if err != nil {
        return err
    }

    for _, productGroup := range productGroupList {
        err = productGroup.ActualizeProductListPromotion(promotion)
        if err != nil {
            return err
        }
    }

    err = h.promotionRepository.Save(ctx, promotion)
    if err != nil {
        return err
    }

    return nil
}
