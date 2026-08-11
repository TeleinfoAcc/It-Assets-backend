package models

type AbcAssetRent struct {
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
	Asset_type     string `json:"asset_type"`
	Asset_project  string `json:"asset_project"`
	Com_hdd        string `json:"com_hdd"`
	Com_wifi_mac   string `json:"com_wifi_mac"`
	Com_lan_mac    string `json:"com_lan_mac"`
	Com_adapt_sn   string `json:"com_adapt_sn"`
	Com_mouse_sn   string `json:"com_mouse_sn"`
	Com_ssd        string `json:"com_ssd"`
	Com_usb_sn     string `json:"com_usb_sn"`
}

func (AbcAssetRent) TableName() string {
	return "abcinv.abc_asset_rent"
}
