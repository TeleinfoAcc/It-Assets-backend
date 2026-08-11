package models

type AbcSiteRoom struct {
	Site_code       string `json:"site_code"`
	Room_code       string `json:"room_code"`
	Room_name       string `json:"room_name"`
	Room_total_seat uint   `json:"room_total_seat"`
	Room_row        uint   `json:"room_row"`
	Room_col        uint   `json:"room_col"`
	Is_active       uint   `json:"is_active"`
}

func (AbcSiteRoom) TableName() string {
	return "abcinv.abc_site_room"
}
