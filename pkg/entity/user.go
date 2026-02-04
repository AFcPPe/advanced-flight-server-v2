package entity

import "time"

// User 用户实体
type User struct {
	ID                      int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username                string     `gorm:"column:username;type:varchar(12);primaryKey;not null" json:"username"`
	Password                string     `gorm:"column:password;type:varchar(255);not null" json:"-"`
	Email                   string     `gorm:"column:email;type:varchar(255);not null" json:"email"`
	Level                   string     `gorm:"column:level;type:varchar(10);default:'1';not null" json:"level"`
	TempLevel               *int       `gorm:"column:temp_level" json:"temp_level"`
	IsPassExam              *int       `gorm:"column:is_pass_exam" json:"is_pass_exam"`
	Avatar                  *string    `gorm:"column:avatar;type:varchar(255)" json:"avatar"`
	Introduce               *string    `gorm:"column:introduce;type:varchar(255)" json:"introduce"`
	RealName                *string    `gorm:"column:real_name;type:varchar(255)" json:"real_name"`
	Authorized              *string    `gorm:"column:authorized;type:varchar(255)" json:"authorized"`
	AuthoInfo               *string    `gorm:"column:authoinfo;type:varchar(255)" json:"authoinfo"`
	AuthorizeType           *string    `gorm:"column:authorizetype;type:varchar(255)" json:"authorizetype"`
	SubAuthoInfo            *string    `gorm:"column:subauthoinfo;type:varchar(255)" json:"subauthoinfo"`
	PilotTime               *int64     `gorm:"column:pilottime" json:"pilottime"`
	AtcTime                 *int64     `gorm:"column:atctime" json:"atctime"`
	LastFlight              *string    `gorm:"column:lastflight;type:varchar(255)" json:"lastflight"`
	LastAtc                 *string    `gorm:"column:lastatc;type:varchar(255)" json:"lastatc"`
	Markup                  *string    `gorm:"column:markup;type:varchar(255)" json:"markup"`
	Timeout                 *time.Time `gorm:"column:timeout" json:"timeout"`
	Testable                *bool      `gorm:"column:testable" json:"testable"`
	IsBind                  *int       `gorm:"column:isbind" json:"isbind"`
	IsVerified              *int       `gorm:"column:isVerified" json:"is_verified"`
	AllowedNotify           *bool      `gorm:"column:allowed_notify" json:"allowed_notify"`
	InControllersFlightTime *int64     `gorm:"column:in_controllers_flight_time" json:"in_controllers_flight_time"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}
