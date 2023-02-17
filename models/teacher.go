package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"learning-management-system/helpers"
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

func createTeachersIfNotExist(emails []string, db *gorm.DB) (err error) {

	teachers := helpers.Map(emails, func(email string) *Teacher {
		return &Teacher{Email: email}
	})
	err = db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(teachers).Error
	return err
}

func deleteAllTeachers(db *gorm.DB) error {
	return db.Exec("DELETE FROM teachers").Error
}
