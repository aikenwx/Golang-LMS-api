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

func updateStudent(email string, isSuspended bool, db *gorm.DB) (err error) {
	student := &Student{Email: email, IsSuspended: isSuspended}

	updateQuery := db.Model(student).Where("email = ?", email).Updates(student)
	err = updateQuery.Error

	return err
}

func deleteAllStudents(db *gorm.DB) error {
	return db.Exec("DELETE FROM students").Error
}

func getAllStudentsExistingInList(studentEmails []string, db *gorm.DB) (existingStudentEmails []string, err error) {
	var existingStudents []Student

	err = db.Select("students.email").Where("students.email in ?", studentEmails).Find(&existingStudents).Error

	if err != nil {
		return nil, err
	}

	existingStudentEmails = helpers.Map(existingStudents, func(student Student) string {
		return student.Email
	})

	return existingStudentEmails, nil
}

func verifyStudentExist(studentEmail string, db *gorm.DB) (studentExists bool, err error) {
	var existingStudent Student

	err = db.Select("students.email").Where("students.email = ?", studentEmail).Find(&existingStudent).Error

	if err != nil {
		return false, err
	}

	if existingStudent.Email == "" {
		return false, nil
	}

	return true, nil
}
