package commands

import (
    "context"

    "github.com/execut/ozon-promotions/application/contract"
    "github.com/execut/ozon-promotions/domain"
)

type PromotionActualizeForGroupHandler struct {
    repository          domain.ProductGroupRepository
    promotionRepository domain.PromotionRepository
}

func NewPromotionActualizeForGroupHandler(repository domain.ProductGroupRepository, promotionRepository domain.PromotionRepository) PromotionActualizeForGroupHandler {
    return PromotionActualizeForGroupHandler{repository: repository, promotionRepository: promotionRepository}
}

func (h *PromotionActualizeForGroupHandler) PromotionActualizeForGroup(ctx context.Context, cmd contract.PromotionActualizeForGroup) error {
    //productGroup, err := domain.NewProductGroup(cmd.GroupID)
    //if err != nil {
    //    return err
    //}
    //
    //err = h.repository.Load(ctx, productGroup)
    //if err != nil {
    //    return err
    //}
    //
    //promotionID, err := domain.NewPromotionID(cmd.PromotionID)
    //if err != nil {
    //    return err
    //}
    //
    //promotion, err := h.promotionRepository.Load(ctx, promotionID)
    //if err != nil {
    //    return err
    //}
    //
    //err = productGroup.ActualizeProductListPromotion(promotion)
    //if err != nil {
    //    return err
    //}
    //
    //err = h.repository.Save(ctx, productGroup)
    //if err != nil {
    //    return err
    //}

    return nil
}
