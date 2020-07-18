package transh

import (
	"github.com/fabric-app/models"
	"github.com/jinzhu/gorm"
)

type Sensor struct {
	ID       string `json:"id" gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Province string `json:"province"  gorm:"type:varchar(10)"`
	City     string `json:"city" gorm:"type:varchar(10)"`
	Country  string `json:"country" gorm:"type:varchar(10)"`
	Code     string `json:"code" gorm:"type:varchar(10)"`
	TypeID   string `json:"type_id" gorm:"type:varchar(10)"`
	Desc     string `json:"desc" gorm:"type:varchar(50)"`
	PicID    string `json:"pic_id" gorm:"type:varchar(10)"`
}

// get all sensors
func FindAllSensors() (*[]Sensor, error) {
	var sensors []Sensor
	err := models.db.Select("code,desc").Find(&sensors).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &sensors, err
}

// get each sensor's data number
func CountSensorsData(id Sensor) (int64, error) {
	var cnt int64
	err := models.db.Model(&Sensor{}).Where("Code = ?", id).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return cnt, err
}
