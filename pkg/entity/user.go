package entity

// User 用户实体
type User struct {
	BaseModel
	CID          string `gorm:"column:cid;type:varchar(32);uniqueIndex;not null" json:"cid"`
	Password     string `gorm:"column:password;type:varchar(255);not null" json:"-"`
	Callsign     string `gorm:"column:callsign;type:varchar(32)" json:"callsign"`
	RealName     string `gorm:"column:real_name;type:varchar(100)" json:"real_name"`
	Email        string `gorm:"column:email;type:varchar(255);uniqueIndex" json:"email"`
	Rating       int    `gorm:"column:rating;default:1" json:"rating"`
	PilotRating  int    `gorm:"column:pilot_rating;default:0" json:"pilot_rating"`
	Status       int8   `gorm:"column:status;default:1;index" json:"status"` // 1: active, 0: inactive, -1: banned
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == 1
}

// IsBanned 检查用户是否被禁用
func (u *User) IsBanned() bool {
	return u.Status == -1
}
