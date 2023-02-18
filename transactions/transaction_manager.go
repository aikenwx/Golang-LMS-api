package transactions

import (
	"fmt"
	"learning-management-system/database"
	"learning-management-system/helpers"
	"learning-management-system/models"
	"learning-management-system/repositories"
	"strings"
)

type TransactionManager struct {
	studentRepo              *repositories.StudentRepo
	teacherRepo              *repositories.TeacherRepo
	registerRelationshipRepo *repositories.RegisterRelationshipRepo
}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{
		studentRepo:              repositories.NewStudentRepo(),
		teacherRepo:              repositories.NewTeacherRepo(),
		registerRelationshipRepo: repositories.NewRegisterRelationshipRepo(),
	}
}

func (transactionManager *TransactionManager) RegisterStudentsToTeacher(teacherEmail string, studentEmails []string, connection *database.Connection) error {
	db := connection.GetDb()
	return transactionManager.registerRelationshipRepo.CreateOneTeacherToManyStudentsRegisterRelationshipsIfNotExists(teacherEmail, studentEmails, db)
}

func (transactionManager *TransactionManager) SuspendStudent(studentEmail string, connection *database.Connection) error {
	db := connection.GetDb()
	studentToUpdate := models.Student{Email: studentEmail, IsSuspended: true}
	return transactionManager.studentRepo.UpdateStudent(&studentToUpdate, db)
}

func (transactionManager *TransactionManager) RetrieveCommonStudentEmails(teacherEmails []string, connection *database.Connection) ([]string, error) {
	db := connection.GetDb()

	relationships, err := transactionManager.registerRelationshipRepo.GetRelationshipsByTeacherEmails(teacherEmails, db)

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

func (transactionManager *TransactionManager) RetrieveStudentRecipients(teacherEmail string, mentionedStudentEmails []string, connection *database.Connection) ([]string, error) {
	db := connection.GetDb()

	relationships, err := transactionManager.registerRelationshipRepo.GetRelationshipsByTeacherEmail(teacherEmail, db)

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
	mentionedStudents, err := transactionManager.studentRepo.GetStudentsByEmails(mentionedStudentEmails, db)

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

func (transactionManager *TransactionManager) ClearDatabase(connection *database.Connection) error {
	db := connection.GetDb()
	tx := db.Begin()
	if err := transactionManager.registerRelationshipRepo.DeleteAllRegisterRelationships(tx); err != nil {
		tx.Rollback()
	}
	if err := transactionManager.studentRepo.DeleteAllStudents(tx); err != nil {
		tx.Rollback()
	}
	if err := transactionManager.teacherRepo.DeleteAllTeachers(tx); err != nil {
		tx.Rollback()
	}
	return tx.Commit().Error
}

func (transactionManager *TransactionManager) PopulateStudents(studentEmails []string, connection *database.Connection) error {
	db := connection.GetDb()
	err := transactionManager.studentRepo.CreateStudentsIfNotExist(studentEmails, db)
	return err
}

func (transactionManager *TransactionManager) PopulateTeachers(teacherEmails []string, connection *database.Connection) error {
	db := connection.GetDb()
	err := transactionManager.teacherRepo.CreateTeachersIfNotExist(teacherEmails, db)
	return err
}

func (transactionManager *TransactionManager) ValidateStudentsExists(studentEmails []string, connection *database.Connection) (userError error, dbError error) {
	db := connection.GetDb()

	if len(studentEmails) == 0 {
		return nil, nil
	}

	if len(studentEmails) == 1 {
		if student, dbError := transactionManager.studentRepo.GetStudentByEmail(studentEmails[0], db); dbError != nil {
			return nil, dbError
		} else if student == nil {
			return generateNonExistentStudentsError(studentEmails), nil
		}
	}

	existingStudents, err := transactionManager.studentRepo.GetStudentsByEmails(studentEmails, db)
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

func (transactionManager *TransactionManager) ValidateTeachersExists(teacherEmails []string, connection *database.Connection) (userError error, dbError error) {
	db := connection.GetDb()

	if len(teacherEmails) == 0 {
		return nil, nil
	}

	if len(teacherEmails) == 1 {
		if teacher, dbError := transactionManager.teacherRepo.GetTeacherByEmail(teacherEmails[0], db); dbError != nil {
			return nil, dbError
		} else if teacher == nil {
			return generateNonExistentTeachersError(teacherEmails), nil
		}
	}

	existingTeachers, err := transactionManager.teacherRepo.GetTeachersByEmails(teacherEmails, db)
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
