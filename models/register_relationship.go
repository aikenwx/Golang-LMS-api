package models

type RegisterRelationship struct {
	TeacherEmail      string  `gorm:"primaryKey"`
	Teacher           Teacher `gorm:"foreignKey:TeacherEmail"`
	StudentEmail      string  `gorm:"primaryKey"`
	RegisteredStudent Student `gorm:"foreignKey:StudentEmail"`
}
