package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"learning-management-system/helpers"
	"learning-management-system/models"
)

type TeacherRepo struct{}

func NewTeacherRepo() *TeacherRepo {
	return &TeacherRepo{}
}

func (*TeacherRepo) CreateTeachersIfNotExist(emails []string, db *gorm.DB) (err error) {

	teachers := helpers.Map(emails, func(email string) *models.Teacher {
		return &models.Teacher{Email: email}
	})
	err = db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(teachers).Error
	return err
}

func (*TeacherRepo) DeleteAllTeachers(db *gorm.DB) error {
	return db.Exec("DELETE FROM teachers").Error
}

func (*TeacherRepo) GetTeachersByEmails(teacherEmails []string, db *gorm.DB) (teachers []*models.Teacher, err error) {
	err = db.Table("teachers").Where("teachers.email in ?", teacherEmails).Find(&teachers).Error
	return teachers, err
}

func (*TeacherRepo) GetTeacherByEmail(teacherEmail string, db *gorm.DB) (teacher *models.Teacher, err error) {
	teacher = &models.Teacher{}
	err = db.Table("teachers").Where("teachers.email = ?", teacherEmail).Find(&teacher).Error
	if teacher.Email == "" {
		return nil, err
	}
	return teacher, err
}
