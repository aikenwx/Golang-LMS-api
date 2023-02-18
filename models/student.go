package models

type Student struct {
	Email       string `gorm:"primaryKey"`
	IsSuspended bool
}
