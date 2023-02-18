package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"learning-management-system/helpers"
	"learning-management-system/models"
)

type StudentRepo struct{}

func NewStudentRepo() *StudentRepo {
	return &StudentRepo{}
}

func (*StudentRepo) CreateStudentsIfNotExist(studentEmails []string, db *gorm.DB) (err error) {

	students := helpers.Map(studentEmails, func(studentEmail string) *models.Student {
		return &models.Student{Email: studentEmail, IsSuspended: false}
	})
	err = db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(students).Error
	return err
}

func (*StudentRepo) UpdateStudent(studentToUpdate *models.Student, db *gorm.DB) (err error) {

	updateQuery := db.Table("students").Where("email = ?", studentToUpdate.Email).Updates(studentToUpdate)
	err = updateQuery.Error

	return err
}

func (*StudentRepo) DeleteAllStudents(db *gorm.DB) error {
	return db.Exec("DELETE FROM students").Error
}

func (*StudentRepo) GetStudentsByEmails(studentEmails []string, db *gorm.DB) (students []*models.Student, err error) {
	err = db.Table("students").Where("students.email in ?", studentEmails).Find(&students).Error
	return students, err
}

func (*StudentRepo) GetStudentByEmail(studentEmail string, db *gorm.DB) (student *models.Student, err error) {
	student = &models.Student{}
	err = db.Table("students").Where("students.email = ?", studentEmail).Find(&student).Error
	if student.Email == "" {
		return nil, err
	}
	return student, err
}
