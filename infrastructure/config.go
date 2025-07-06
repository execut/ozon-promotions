package infrastructure

type Config struct {
    ApiKey    string `yaml:"apiKey"`
    CompanyID int64  `yaml:"companyID"`
    Cookie    string `yaml:"cookie"`
    Groups    []struct {
        ID               string   `yaml:"id"`
        Articles         []string `yaml:"articles"`
        InPromotions     uint8    `yaml:"in-promotions"`
        MinimalPromotion bool     `yaml:"minimal-promotion"`
    } `yaml:"groups"`
}
