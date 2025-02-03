package outage


type OutageConfig struct {
	Id                     int64    `gorm:"primary_key"`
	EnableOutage      bool   `gorm:"default:false"`
	MinorOutageName   string `gorm:"default:'Minor Outage'"`
	MinorOutageColor  string `gorm:"default:'yellow'"`
	MajorOutageName   string `gorm:"default:'Major Outage'"`
	MajorOutageColor  string `gorm:"default:'orange'"`
}