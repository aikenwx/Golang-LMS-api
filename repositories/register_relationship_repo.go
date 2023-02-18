package repositories

import (
	"gorm.io/gorm"
	"learning-management-system/helpers"
	"learning-management-system/models"

	"gorm.io/gorm/clause"
)

type RegisterRelationshipRepo struct{}

func NewRegisterRelationshipRepo() *RegisterRelationshipRepo {
	return &RegisterRelationshipRepo{}
}

func (*RegisterRelationshipRepo) CreateOneTeacherToManyStudentsRegisterRelationshipsIfNotExists(teacherEmail string, studentEmails []string, db *gorm.DB) (err error) {
	registerRelationships := helpers.Map(studentEmails, func(studentEmail string) *models.RegisterRelationship {
		return &models.RegisterRelationship{TeacherEmail: teacherEmail, StudentEmail: studentEmail}
	})
	err = db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&registerRelationships).Error

	return err
}

func (*RegisterRelationshipRepo) GetRelationshipsByTeacherEmail(teacherEmail string, db *gorm.DB) (relationships []*models.RegisterRelationship, err error) {
	err = db.Table("register_relationships").Preload("Teacher").Preload("RegisteredStudent").Where("teacher_email = ?", teacherEmail).Find(&relationships).Error
	return relationships, err
}

func (*RegisterRelationshipRepo) GetRelationshipsByTeacherEmails(teacherEmails []string, db *gorm.DB) (relationships []*models.RegisterRelationship, err error) {
	err = db.Table("register_relationships").Preload("Teacher").Preload("RegisteredStudent").Where("teacher_email in ?", teacherEmails).Find(&relationships).Error
	return relationships, err
}

func (*RegisterRelationshipRepo) DeleteAllRegisterRelationships(db *gorm.DB) error {
	return db.Exec("DELETE FROM register_relationships").Error
}
