package models

import "time"

type IamProject struct {
	Proj_id         uint      `json:"proj_id" gorm:"primaryKey"`
	Proj_name       string    `json:"proj_name"`
	Proj_desc       string    `json:"proj_desc"`
	Start_date      time.Time `json:"start_date"`
	End_date        time.Time `json:"end_date"`
	Is_active       uint      `json:"is_active"`
	Main_shif_id    uint      `json:"main_shif_id"`
	Proj_cost       string    `json:"proj_cost"`
	Proj_infra_desc string    `json:"proj_infra_desc"`
}

func (IamProject) TableName() string {
	return "qamon.iam_project"
}
