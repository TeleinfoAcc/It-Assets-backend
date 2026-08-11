package models

type AbcSite struct {
	Site_address string `json:"site_address"`
	Site_code    string `json:"site_code"`
	Site_name    string `json:"site_name"`
	Is_active    uint   `json:"is_active"`
}

func (AbcSite) TableName() string {
	return "abcinv.abc_site"
}
