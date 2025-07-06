package infrastructure

import (
    "encoding/json"
    "strconv"

    "github.com/execut/ozon-reports-downloader/common"
)

type PromotionClient struct {
    companyID int64
    cookie    string
}

func NewPromotionClient(companyID int64, cookie string) *PromotionClient {
    return &PromotionClient{
        companyID: companyID,
        cookie:    cookie,
    }
}

func (c *PromotionClient) ProductCandidateList(actionID int64) ([]ProductCandidateListResponseProduct, error) {
    var offset int64 = 0
    var limit int64 = 100
    var productList []ProductCandidateListResponseProduct
    for {
        resp, err := common.DoGetRequest(struct{}{}, "https://seller.ozon.ru/api/site/sa-facade/v2/seller-actions/"+strconv.FormatInt(actionID, 10)+"/product/candidate?offset="+strconv.FormatInt(offset, 10)+"&limit="+strconv.FormatInt(limit, 10), c.cookie, c.companyID)
        if err != nil {
            return nil, err
        }

        respObject := &ProductCandidateListResponse{}
        err = json.Unmarshal(resp, respObject)
        if err != nil {
            return nil, err
        }

        if productList == nil {
            total, err := strconv.ParseInt(respObject.Total, 10, 64)
            if err != nil {
                return nil, err
            }

            productList = make([]ProductCandidateListResponseProduct, 0, total)
        }

        productList = append(productList, respObject.Products...)
        if len(respObject.Products) != int(limit) {
            break
        }
    }

    return productList, nil
}

func (c *PromotionClient) PromotionProductList(actionID int64) ([]PromotionProductListResponseProduct, error) {
    var offset int64 = 0
    var limit int64 = 100
    var productList []PromotionProductListResponseProduct
    for {
        resp, err := common.DoGetRequest(struct{}{}, "https://seller.ozon.ru/api/site/sa-facade/v2/seller-actions/"+strconv.FormatInt(actionID, 10)+"/product/active?offset="+strconv.FormatInt(offset, 10)+"&limit="+strconv.FormatInt(limit, 10), c.cookie, c.companyID)
        if err != nil {
            return nil, err
        }

        respObject := &PromotionProductListResponse{}
        err = json.Unmarshal(resp, respObject)
        if err != nil {
            return nil, err
        }

        if productList == nil {
            total, err := strconv.ParseInt(respObject.Total, 10, 64)
            if err != nil {
                return nil, err
            }

            productList = make([]PromotionProductListResponseProduct, 0, total)
        }

        productList = append(productList, respObject.Products...)
        if len(respObject.Products) != int(limit) {
            break
        }
    }

    return productList, nil
}

type PromotionProductListResponseProduct struct {
    Id                        string   `json:"id"`
    OfferId                   string   `json:"offerId"`
    Skus                      []string `json:"skus"`
    OzonSku                   string   `json:"ozonSku"`
    Name                      string   `json:"name"`
    ItemType                  string   `json:"itemType"`
    DescriptionCategoryId     string   `json:"descriptionCategoryId"`
    Price                     int      `json:"price"`
    ActionPrice               int      `json:"actionPrice"`
    BasePrice                 int      `json:"basePrice"`
    MaxDiscountPrice          int      `json:"maxDiscountPrice"`
    MinSellerPrice            int      `json:"minSellerPrice"`
    MarketplaceSellerPrice    int      `json:"marketplaceSellerPrice"`
    ActionPriceToAutoAdd      int      `json:"actionPriceToAutoAdd"`
    OzonStock                 string   `json:"ozonStock"`
    SellerStock               string   `json:"sellerStock"`
    TotalStock                string   `json:"totalStock"`
    MinActionQuantity         string   `json:"minActionQuantity"`
    Quantity                  string   `json:"quantity"`
    QuantityToAutoAdd         string   `json:"quantityToAutoAdd"`
    Currency                  string   `json:"currency"`
    Thumbnail                 string   `json:"thumbnail"`
    PriceReferenceForBoosting struct {
        BoostingX2Price  int `json:"boostingX2Price"`
        BoostingX3Price  int `json:"boostingX3Price"`
        BoostingX4Price  int `json:"boostingX4Price"`
        MaxBoostingPrice int `json:"maxBoostingPrice"`
    } `json:"priceReferenceForBoosting"`
    BoostingInSearch struct {
        BoostingPercent    float64 `json:"boostingPercent"`
        MaxBoostingPercent int     `json:"maxBoostingPercent"`
        BoostingScaleColor string  `json:"boostingScaleColor"`
    } `json:"boostingInSearch"`
    IsManuallyAdded bool `json:"isManuallyAdded"`
}

type PromotionProductListResponse struct {
    Products []PromotionProductListResponseProduct `json:"products"`
    Total    string                                `json:"total"`
}

type ProductCandidateListResponseProduct struct {
    Id                         string        `json:"id"`
    Price                      int           `json:"price"`
    ActionPrice                int           `json:"actionPrice"`
    MinDiscount                float64       `json:"minDiscount"`
    BasePrice                  int           `json:"basePrice"`
    Name                       string        `json:"name"`
    OfferID                    string        `json:"offerID"`
    ItemType                   string        `json:"itemType"`
    OzonStock                  string        `json:"ozonStock"`
    SellerStock                string        `json:"sellerStock"`
    CommissionPercentFBS       int           `json:"commissionPercentFBS"`
    CommissionPercentFBO       int           `json:"commissionPercentFBO"`
    Skus                       []string      `json:"skus"`
    MaxDiscountPrice           int           `json:"maxDiscountPrice"`
    AddMode                    string        `json:"addMode"`
    MinSellerPrice             int           `json:"minSellerPrice"`
    AutoAdd                    bool          `json:"autoAdd"`
    Thumbnail                  string        `json:"thumbnail"`
    MinActionQuantity          string        `json:"minActionQuantity"`
    Quantity                   string        `json:"quantity"`
    TotalStock                 string        `json:"totalStock"`
    IsFboDisabled              bool          `json:"isFboDisabled"`
    IsFbsDisabled              bool          `json:"isFbsDisabled"`
    IsActive                   bool          `json:"isActive"`
    MaxAdditionalDiscountPrice int           `json:"maxAdditionalDiscountPrice"`
    FboSku                     string        `json:"fboSku"`
    FbsSku                     string        `json:"fbsSku"`
    Currency                   string        `json:"currency"`
    FairDiscountFailed         bool          `json:"fairDiscountFailed"`
    RecommendedQuantity        string        `json:"recommendedQuantity"`
    RecommendedActionPrice     int           `json:"recommendedActionPrice"`
    MetazonCategoryId          string        `json:"metazonCategoryId"`
    FairDiscountFailedReasons  []interface{} `json:"fairDiscountFailedReasons"`
    PriceIndexEdlp             int           `json:"priceIndexEdlp"`
    OzonSku                    string        `json:"ozonSku"`
    Ads                        int           `json:"ads"`
    Idc                        int           `json:"idc"`
    RemainingActionStock       string        `json:"remainingActionStock"`
    IsActionStockSold          bool          `json:"isActionStockSold"`
    DescriptionCategoryId      string        `json:"descriptionCategoryId"`
    MarketplaceSellerPrice     int           `json:"marketplaceSellerPrice"`
    AllowStockEdit             bool          `json:"allowStockEdit"`
    Quant                      string        `json:"quant"`
    QuantQuantity              string        `json:"quantQuantity"`
    PriceReferenceForBoosting  struct {
        BoostingX2Price  int `json:"boostingX2Price"`
        BoostingX3Price  int `json:"boostingX3Price"`
        BoostingX4Price  int `json:"boostingX4Price"`
        MaxBoostingPrice int `json:"maxBoostingPrice"`
    } `json:"priceReferenceForBoosting"`
    BoostingInSearch struct {
        BoostingPercent    int    `json:"boostingPercent"`
        MaxBoostingPercent int    `json:"maxBoostingPercent"`
        BoostingScaleColor string `json:"boostingScaleColor"`
    } `json:"boostingInSearch"`
    SlowQuants            string `json:"slowQuants"`
    ActionDiscountPercent int    `json:"actionDiscountPercent"`
}

type ProductCandidateListResponse struct {
    Products     []ProductCandidateListResponseProduct `json:"products"`
    Total        string                                `json:"total"`
    SearchLastId string                                `json:"searchLastId"`
    LastId       string                                `json:"lastId"`
}

func (c *PromotionClient) ActivateProductList(actionID int64, request ActivateProductListRequest) (ActivateProductListResponse, error) {
    resp, err := common.DoRequest(request, "https://seller.ozon.ru/api/site/sa-facade/v2/seller-actions/"+strconv.FormatInt(actionID, 10)+"/product/activate", c.cookie, c.companyID)
    if err != nil {
        return ActivateProductListResponse{}, err
    }

    var respObject ActivateProductListResponse
    err = json.Unmarshal(resp, &respObject)
    if err != nil {
        return ActivateProductListResponse{}, err
    }

    return respObject, nil
}

type ActivateProductListRequestProduct struct {
    Id              string `json:"id"`
    ActionPrice     int    `json:"actionPrice"`
    DiscountPercent int    `json:"discountPercent"`
    Quantity        int    `json:"quantity"`
    Currency        string `json:"currency"`
    Sku             string `json:"sku"`
}

type ActivateProductListRequest struct {
    Products []ActivateProductListRequestProduct `json:"products"`
}

type ActivateProductListResponse struct {
    ProductIds []string      `json:"productIds"`
    Rejected   []interface{} `json:"rejected"`
}

func (c *PromotionClient) DeactivateProductList(actionID int64, request DeactivateProductListRequest) (DeactivateProductListResponse, error) {
    resp, err := common.DoRequest(request, "https://seller.ozon.ru/api/site/sa-facade/v2/seller-actions/"+strconv.FormatInt(actionID, 10)+"/product/deactivate", c.cookie, c.companyID)
    if err != nil {
        return DeactivateProductListResponse{}, err
    }

    var respObject DeactivateProductListResponse
    err = json.Unmarshal(resp, &respObject)
    if err != nil {
        return DeactivateProductListResponse{}, err
    }

    return respObject, nil
}

type DeactivateProductListRequest struct {
    ProductIds []string `json:"product_ids"`
    Skus       []string `json:"skus"`
}

type DeactivateProductListResponse struct {
    ProductIds []string      `json:"productIds"`
    Rejected   []interface{} `json:"rejected"`
}
