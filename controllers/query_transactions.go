package controllers

import (
	"fmt"
	"learning-management-system/database"
	"learning-management-system/helpers"
	"learning-management-system/models"
	"learning-management-system/repositories"
	"strings"
)

var teacherRepo *repositories.TeacherRepo
var studentRepo *repositories.StudentRepo
var registerRelationshipRepo *repositories.RegisterRelationshipRepo

func InitRepositories() {
	teacherRepo = repositories.NewTeacherRepo()
	studentRepo = repositories.NewStudentRepo()
	registerRelationshipRepo = repositories.NewRegisterRelationshipRepo()
}

func RegisterStudentsToTeacher(teacherEmail string, studentEmails []string, connection *database.Connection) error {
	db := connection.GetDb()
	return registerRelationshipRepo.CreateOneTeacherToManyStudentsRegisterRelationshipsIfNotExists(teacherEmail, studentEmails, db)
}

func SuspendStudent(studentEmail string, connection *database.Connection) error {
	db := connection.GetDb()
	studentToUpdate := models.Student{Email: studentEmail, IsSuspended: true}
	return studentRepo.UpdateStudent(&studentToUpdate, db)
}

func RetrieveCommonStudentEmails(teacherEmails []string, connection *database.Connection) ([]string, error) {
	db := connection.GetDb()

	relationships, err := registerRelationshipRepo.GetRelationshipsByTeacherEmails(teacherEmails, db)

	if err != nil {
		return nil, err
	}

	// create a map of student emails to their count
	studentEmailCountMap := make(map[string]int)
	for _, relationship := range relationships {
		studentEmailCountMap[relationship.StudentEmail]++
	}

	// fetch all students with count equal to number of teachers
	commonStudents := make([]string, 0)
	for studentEmail, count := range studentEmailCountMap {
		if count == len(teacherEmails) {
			commonStudents = append(commonStudents, studentEmail)
		}
	}

	return commonStudents, nil
}

func RetrieveStudentRecipients(teacherEmail string, mentionedStudentEmails []string, connection *database.Connection) ([]string, error) {
	db := connection.GetDb()

	relationships, err := registerRelationshipRepo.GetRelationshipsByTeacherEmail(teacherEmail, db)

	if err != nil {
		return nil, err
	}
	// map relationships to students
	students := helpers.Map(relationships, func(relationship *models.RegisterRelationship) *models.Student {
		return &relationship.RegisteredStudent
	})

	if err != nil {
		return nil, err
	}

	// get mentioned students
	mentionedStudents, err := studentRepo.GetStudentsByEmails(mentionedStudentEmails, db)

	// append mentioned students to students
	students = append(students, mentionedStudents...)

	// filter out suspended students
	students = helpers.Filter(students, func(student *models.Student) bool {
		return !student.IsSuspended
	})

	// map students to emails
	studentEmails := helpers.Map(students, func(student *models.Student) string {
		return student.Email
	})

	// remove duplicates
	studentEmails = helpers.RemoveDuplicatesInStringSlice(studentEmails)

	return studentEmails, nil
}

func ClearDatabase(connection *database.Connection) error {
	db := connection.GetDb()
	tx := db.Begin()
	if err := registerRelationshipRepo.DeleteAllRegisterRelationships(tx); err != nil {
		tx.Rollback()
	}
	if err := studentRepo.DeleteAllStudents(tx); err != nil {
		tx.Rollback()
	}
	if err := teacherRepo.DeleteAllTeachers(tx); err != nil {
		tx.Rollback()
	}
	return tx.Commit().Error
}

func PopulateStudents(studentEmails []string, connection *database.Connection) error {
	db := connection.GetDb()
	err := studentRepo.CreateStudentsIfNotExist(studentEmails, db)
	return err
}

func PopulateTeachers(teacherEmails []string, connection *database.Connection) error {
	db := connection.GetDb()
	err := teacherRepo.CreateTeachersIfNotExist(teacherEmails, db)
	return err
}

func ValidateStudentsExists(studentEmails []string, connection *database.Connection) (userError error, dbError error) {
	db := connection.GetDb()

	if len(studentEmails) == 0 {
		return nil, nil
	}

	if len(studentEmails) == 1 {
		if student, dbError := studentRepo.GetStudentByEmail(studentEmails[0], db); dbError != nil {
			return nil, dbError
		} else if student == nil {
			return generateNonExistentStudentsError(studentEmails), nil
		}
	}

	existingStudents, err := studentRepo.GetStudentsByEmails(studentEmails, db)
	if err != nil {
		return nil, err
	}

	// map students to emails
	existingStudentEmails := helpers.Map(existingStudents, func(student *models.Student) string {
		return student.Email
	})

	nonExistentStudentEmails := helpers.RemoveAllStringsInSlice(studentEmails, existingStudentEmails)

	return generateNonExistentStudentsError(nonExistentStudentEmails), nil
}

func ValidateTeachersExists(teacherEmails []string, connection *database.Connection) (userError error, dbError error) {
	db := connection.GetDb()

	if len(teacherEmails) == 0 {
		return nil, nil
	}

	if len(teacherEmails) == 1 {
		if teacher, dbError := teacherRepo.GetTeacherByEmail(teacherEmails[0], db); dbError != nil {
			return nil, dbError
		} else if teacher == nil {
			return generateNonExistentTeachersError(teacherEmails), nil
		}
	}

	existingTeachers, err := teacherRepo.GetTeachersByEmails(teacherEmails, db)
	if err != nil {
		return nil, err
	}

	// map teachers to emails
	existingTeacherEmails := helpers.Map(existingTeachers, func(teacher *models.Teacher) string {
		return teacher.Email
	})

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
