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

func getAllTeachersExistingInList(teacherEmails []string, db *gorm.DB) ([]string, error) {
	var existingTeachers []Teacher

	err := db.Select("teachers.email").Where("teachers.email in ?", teacherEmails).Find(&existingTeachers).Error

	if err != nil {
		return nil, err
	}

	existingTeacherEmails := helpers.Map(existingTeachers, func(teacher Teacher) string {
		return teacher.Email
	})

	return existingTeacherEmails, nil
}

func verifyTeacherExist(email string, db *gorm.DB) (bool, error) {
	var existingTeacher Teacher

	err := db.Select("teachers.email").Where("teachers.email = ?", email).Find(&existingTeacher).Error

	if err != nil {
		return false, err
	}

	if existingTeacher.Email == "" {
		return false, nil
	}

	return true, nil
}
