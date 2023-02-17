package models

import (
	"fmt"
	"learning-management-system/database"
	"learning-management-system/helpers"
	"strings"
)

func RegisterStudentsToTeacher(teacherEmail string, studentEmails []string, connection *database.Connection) error {

	db := connection.GetDb()
	return createRegisterRelationshipsIfNotExists(teacherEmail, studentEmails, db)
}

func SuspendStudent(studentEmail string, connection *database.Connection) error {
	db := connection.GetDb()
	return updateStudent(studentEmail, true, db)
}

func RetrieveCommonStudents(teacherEmail []string, connection *database.Connection) ([]string, error) {
	db := connection.GetDb()
	return getAllStudentsRegisteredToTeacherEmails(teacherEmail, db)
}

func RetrieveStudentRecipients(teacherEmail string, mentionedStudentEmails []string, connection *database.Connection) ([]string, error) {
	db := connection.GetDb()
	return getUnsuspendedStudentsRegisteredToTeacher(teacherEmail, mentionedStudentEmails, db)
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

func PopulateStudents(studentEmails []string, connection *database.Connection) error {
	db := connection.GetDb()
	err := createStudentsIfNotExist(studentEmails, db)
	return err
}

func PopulateTeachers(teacherEmails []string, connection *database.Connection) error {
	db := connection.GetDb()
	err := createTeachersIfNotExist(teacherEmails, db)
	return err
}

func ValidateStudentsExists(studentEmails []string, connection *database.Connection) (userError error, dbError error) {
	db := connection.GetDb()

	if len(studentEmails) == 0 {
		return nil, nil
	}

	if len(studentEmails) == 1 {
		if studentExists, dbError := verifyStudentExist(studentEmails[0], db); dbError != nil {
			return nil, dbError
		} else if !studentExists {
			return generateNonExistentStudentsError(studentEmails), nil
		}
	}
	existingStudentEmails, err := getAllStudentsExistingInList(studentEmails, db)
	if err != nil {
		return nil, err
	}

	nonExistentStudentEmails := helpers.RemoveAllStringsInSlice(studentEmails, existingStudentEmails)

	return generateNonExistentStudentsError(nonExistentStudentEmails), nil
}

func ValidateTeachersExists(teacherEmails []string, connection *database.Connection) (userError error, dbError error) {
	db := connection.GetDb()

	if len(teacherEmails) == 0 {
		return nil, nil
	}

	if len(teacherEmails) == 1 {

		if teacherExists, dbError := verifyTeacherExist(teacherEmails[0], db); dbError != nil {
			return nil, dbError
		} else if !teacherExists {
			return generateNonExistentTeachersError(teacherEmails), nil
		}
	}

	existingTeacherEmails, err := getAllTeachersExistingInList(teacherEmails, db)
	if err != nil {
		return nil, err
	}

	nonExistentTeacherEmails := helpers.RemoveAllStringsInSlice(teacherEmails, existingTeacherEmails)

	return generateNonExistentTeachersError(nonExistentTeacherEmails), nil
}

func generateNonExistentStudentsError(nonExistentTeacherEmails []string) error {
	if len(nonExistentTeacherEmails) == 0 {
		return nil
	}

	if len(nonExistentTeacherEmails) == 1 {
		return fmt.Errorf("Student with email %s does not exist in the database", nonExistentTeacherEmails[0])
	}

	return fmt.Errorf("Students with emails %s do not exist in the database", strings.Join(nonExistentTeacherEmails, ", "))
}

func generateNonExistentTeachersError(nonExistentTeacherEmails []string) error {
	if len(nonExistentTeacherEmails) == 0 {
		return nil
	}

	if len(nonExistentTeacherEmails) == 1 {
		return fmt.Errorf("Teacher with email %s does not exist in the database", nonExistentTeacherEmails[0])
	}

	return fmt.Errorf("Teachers with emails %s do not exist in the database", strings.Join(nonExistentTeacherEmails, ", "))
}
