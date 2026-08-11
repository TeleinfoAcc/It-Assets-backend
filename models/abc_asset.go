package models

type AbcAsset struct {
	It_asset_id    uint   `json:"it_asset_id" gorm:"primaryKey"`
	Serialnumber   string `json:"serialnumber"`
	Emp_id         string `json:"emp_id"`
	Com_name       string `json:"com_name"`
	Com_local_ip   string `json:"com_local_ip"`
	Com_join_ip    string `json:"com_join_ip"`
	Curr_room_code string `json:"curr_room_code"`
	Com_brand      string `json:"com_brand"`
	Com_model      string `json:"com_model"`
	Com_type       string `json:"com_type"`
	Com_desc1      string `json:"com_desc1"`
	Com_desc2      string `json:"com_desc2"`
	Com_desc3      string `json:"com_desc3"`
	Gl_asset_code  string `json:"gl_asset_code"`
	Loc_type       string `json:"loc_type"`
	Create_date    string `json:"create_date"`
	Mdf_date       string `json:"mdf_date"`
	Asset_status   uint   `json:"asset_status"`
	Loc_seat       string `json:"loc_seat"`
	Location       string `json:"location"`
	Com_status     string `json:"com_status"`
	Cap_date       string `json:"cap_date"`
	Mdf_agent_id   uint   `json:"mdf_agent_id"`
	Iss_date       string `json:"iss_date"`
	Return_date    string `json:"return_date"`
}

func (AbcAsset) TableName() string {
	return "abcinv.abc_asset"
}
