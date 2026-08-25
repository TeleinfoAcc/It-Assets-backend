package models

type AbcHistory struct {
	His_id            uint   `json:"his_id" gorm:"primaryKey"`
	Agent_id          uint   `json:"agent_id"`
	Asset_status_name string `json:"asset_status_name"`
	Com_desc1         string `json:"com_desc1"`
	It_asset_id       uint   `json:"it_asset_id"`
}

func (AbcHistory) TableName() string {
	return "abcinv.abc_history"
}
