package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"learning-management-system/helpers"
)

type Student struct {
	Email       string `gorm:"primaryKey"`
	IsSuspended bool
}

func createStudentsIfNotExist(studentEmails []string, db *gorm.DB) (err error) {

	students := helpers.Map(studentEmails, func(studentEmail string) *Student {
		return &Student{Email: studentEmail, IsSuspended: false}
	})
	err = db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(students).Error
	return err
}

func updateOrCreateStudent(email string, isSuspended bool, db *gorm.DB) (err error) {
	student := &Student{Email: email, IsSuspended: isSuspended}

	updateQuery := db.Model(student).Where("email = ?", email).Updates(student)
	err = updateQuery.Error

	if updateQuery.RowsAffected == 0 {
		err = db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(student).Error
	}

	return err
}

func deleteAllStudents (db *gorm.DB) error {
	return db.Exec("DELETE FROM students").Error
}