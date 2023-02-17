package models

import "learning-management-system/database"

func RegisterStudentsToTeacher(teacherEmail string, studentEmails []string, connection *database.Connection) error {

	db := connection.GetDb()
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

func SuspendStudent(studentEmail string, connection *database.Connection) error {
	db := connection.GetDb()
	return updateOrCreateStudent(studentEmail, true, db)
}

func RetrieveCommonStudents(teacherEmail []string, connection *database.Connection) ([]string, error) {
	db := connection.GetDb()
	tx := db.Begin()
	if err := createTeachersIfNotExist(teacherEmail, tx); err != nil {
		tx.Rollback()
	}
	studentEmails, err := getAllStudentsRegisteredToTeacherEmails(teacherEmail, tx)
	if err != nil {
		tx.Rollback()
	}
	return studentEmails, tx.Commit().Error
}

func RetrieveStudentRecipients(teacherEmail string, mentionedStudentEmails []string, connection *database.Connection) ([]string, error) {

	db := connection.GetDb()
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

func ClearDatabase(connection *database.Connection) error {
	db := connection.GetDb()
	tx := db.Begin()
	if err := deleteAllRegisterRelationships(tx); err != nil {
		tx.Rollback()
	}
	if err := deleteAllStudents(tx); err != nil {
		tx.Rollback()
	}
	if err := deleteAllTeachers(tx); err != nil {
		tx.Rollback()
	}
	return tx.Commit().Error
}
