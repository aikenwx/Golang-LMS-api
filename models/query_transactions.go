package models

import "learning-management-system/database"

func RegisterStudentsToTeacher(teacherEmail string, studentEmails []string) error {

	db := database.GlobalConnection.GetDb()
	tx := db.Begin()
	if err := createTeacherIfNotExist(teacherEmail, tx); err != nil {
		tx.Rollback()
	}
	if err := createStudentsIfNotExist(studentEmails, tx); err != nil {
		tx.Rollback()
	}
	if err := createRegisterRelationshipsIfNotExists(teacherEmail, studentEmails, tx); err != nil {
		tx.Rollback()
	}

	return tx.Commit().Error
}

func SuspendStudent(studentEmail string) error {
	db := database.GlobalConnection.GetDb()
	return updateOrCreateStudent(studentEmail, true, db)
}

func RetrieveCommonStudents(teacherEmail string) ([]string, error) {
	db := database.GlobalConnection.GetDb()
	tx := db.Begin()
	if err := createTeacherIfNotExist(teacherEmail, tx); err != nil {
		tx.Rollback()
	}
	studentEmails, err := getAllStudentsRegisteredToTeacherEmails(teacherEmail, tx)
	if err != nil {
		tx.Rollback()
	}
	return studentEmails, tx.Commit().Error
}

func RetrieveStudentRecipients(teacherEmail string, mentionedStudentEmails []string) ([]string, error) {

	db := database.GlobalConnection.GetDb()
	tx := db.Begin()
	if err := createTeacherIfNotExist(teacherEmail, tx); err != nil {
		tx.Rollback()
	}
	if err := createStudentsIfNotExist(mentionedStudentEmails, tx); err != nil {
		tx.Rollback()
	}

	recipientEmails, err := getUnsuspendedStudentsRegisteredToTeacher(teacherEmail, mentionedStudentEmails, tx)
	if err != nil {
		tx.Rollback()
	}

	return recipientEmails, tx.Commit().Error

}
