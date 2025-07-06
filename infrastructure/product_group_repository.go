package infrastructure

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "strconv"

    "github.com/diphantxm/ozon-api-client/ozon"
    "github.com/execut/ozon-promotions/domain"
    "github.com/stretchr/testify/assert/yaml"
)

type ProductGroupRepository struct {
    config  Config
    ozonAPI *ozon.Client
}

func NewProductGroupRepository(configPath string) *ProductGroupRepository {
    config := readConfig(configPath)
    opts := []ozon.ClientOption{
        ozon.WithAPIKey(config.ApiKey),
        ozon.WithClientId(strconv.FormatInt(config.CompanyID, 10)),
    }
    ozonAPI := ozon.NewClient(opts...)

    return &ProductGroupRepository{
        config:  config,
        ozonAPI: ozonAPI,
    }
}

func (r *ProductGroupRepository) LoadList(ctx context.Context) ([]*domain.ProductGroup, error) {
    result := make([]*domain.ProductGroup, 0, len(r.config.Groups))
    articleList := make([]string, 0, len(r.config.Groups))
    for _, group := range r.config.Groups {
        articleList = append(articleList, group.Articles...)
    }

    productListMap, err := r.productByArticleMap(ctx, articleList)
    if err != nil {
        return nil, err
    }

    for _, group := range r.config.Groups {
        var (
            promotionCriteria domain.PromotionCriteria
        )

        if group.MinimalPromotion {
            promotionCriteria = domain.PromotionCriteriaMinimalPrice
        } else {
            promotionCriteria = domain.PromotionCriteriaMaximalBoosting
        }

        productList := make([]domain.Product, 0, len(group.Articles))
        for _, article := range group.Articles {
            product, ok := productListMap[article]
            if !ok {
                return nil, fmt.Errorf("товар с артикулом %v не найден через API", article)
            }

            productList = append(productList, product)
        }

        inPromotions := group.InPromotions
        if inPromotions == 0 {
            inPromotions = 1
        }
        productForPromotionLimit, err := domain.NewProductForPromotionLimit(inPromotions)
        if err != nil {
            return nil, err
        }

        productGroup, err := domain.NewProductGroup(group.ID, productList, promotionCriteria, productForPromotionLimit)
        if err != nil {
            return nil, err
        }

        result = append(result, productGroup)
    }

    return result, nil
}

func (r *ProductGroupRepository) Load(ctx context.Context, productGroup *domain.ProductGroup) error {
    return nil
}

func (r *ProductGroupRepository) productByArticleMap(ctx context.Context, articleList []string) (map[string]domain.Product, error) {
    resp, err := r.ozonAPI.Products().ListProductsByIDs(ctx, &ozon.ListProductsByIDsParams{
        OfferId: articleList,
    })
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("не верный код ответа от API %v", resp.StatusCode)
    }

    result := make(map[string]domain.Product, len(articleList))
    for _, item := range resp.Items {
        productID, err := domain.NewProductID(item.Id)
        if err != nil {
            return nil, err
        }

        var skuValue int64
        for _, source := range item.Sources {
            if source.Source == "sds" {
                skuValue = source.SKU
                break
            }
        }
        if skuValue == 0 {
            return nil, fmt.Errorf("SKU типа sds не найден для продукта %v", item.Sources)
        }

        sku, err := domain.NewSKU(skuValue)
        if err != nil {
            return nil, err
        }

        minPriceValue, err := strconv.ParseFloat(item.MinPrice, 64)
        if err != nil {
            return nil, err
        }

        minPrice, err := domain.NewPrice(minPriceValue)
        if err != nil {
            return nil, err
        }

        maxPriceValue, err := strconv.ParseFloat(item.Price, 64)
        if err != nil {
            return nil, err
        }

        maxPrice, err := domain.NewPrice(maxPriceValue)
        if err != nil {
            return nil, err
        }

        product, err := domain.NewProduct(productID, sku, minPrice, maxPrice)
        if err != nil {
            return nil, err
        }

        result[item.OfferId] = product
    }

    return result, nil
}

func readConfig(configPath string) Config {
    yamlFile, err := os.ReadFile(configPath)

    if err != nil {
        panic(err)
    }

    var config Config

    err = yaml.Unmarshal(yamlFile, &config)
    if err != nil {
        panic(err)
    }

    return config
}
