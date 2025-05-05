package model

type User struct {
	ID       int    `gorm:"primaryKey"`
	UserID   string `gorm:"type:not null"`
	Password string `gorm:"type:not null"`
}
