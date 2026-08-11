package models

type AbcAssetStatus struct {
	Asset_status      uint   `json:"asset_status" gorm:"primaryKey"`
	Asset_status_name string `json:"asset_status_name"`
	Is_active         uint   `json:"is_active"`
}

func (AbcAssetStatus) TableName() string {
	return "abcinv.abc_asset_status"
}
