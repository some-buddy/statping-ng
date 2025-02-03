package intermediate_status


type IntermediateStatusConfig struct {
	Id                     int64    `gorm:"primary_key"`
	EnableIntermediate      bool   `gorm:"default:false"`
	StatusMinorOutageName   string `gorm:"default:'Minor Outage'"`
	StatusMinorOutageColor  string `gorm:"default:'yellow'"`
	StatusMajorOutageName   string `gorm:"default:'Major Outage'"`
	StatusMajorOutageColor  string `gorm:"default:'orange'"`
}