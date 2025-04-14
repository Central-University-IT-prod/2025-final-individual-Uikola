package entity

type User struct {
	ID         int64  `gorm:"primaryKey"`
	Username   string `gorm:"not null"`
	PlatformID string `gorm:"not null"`
	Role       Role   `gorm:"not null"`
}

type Role string

const (
	Client     Role = "client"
	Advertiser Role = "advertiser"
)
