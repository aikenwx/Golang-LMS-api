package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Teacher struct {
	Email string `gorm:"primaryKey"`
}

type TeacherManager struct {
	email string
}

func createTeacherIfNotExist(email string, db *gorm.DB) (err error) {
	teacher := &Teacher{Email: email}
	err = db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(teacher).Error

	return err
}

func deleteAllTeachers(db *gorm.DB) error {
	return db.Exec("DELETE FROM teachers").Error
}
