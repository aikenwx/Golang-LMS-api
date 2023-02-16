package models

import (
	"gorm.io/gorm"
	"learning-management-system/helpers"

	"gorm.io/gorm/clause"
)

type RegisterRelationship struct {
	TeacherEmail      string  `gorm:"primaryKey"`
	Teacher           Teacher `gorm:"foreignKey:TeacherEmail"`
	StudentEmail      string  `gorm:"primaryKey"`
	RegisteredStudent Student `gorm:"foreignKey:StudentEmail"`
}

func createRegisterRelationshipsIfNotExists(teacherEmail string, studentEmails []string, db *gorm.DB) (err error) {
	registerRelationships := helpers.Map(studentEmails, func(studentEmail string) *RegisterRelationship {
		return &RegisterRelationship{TeacherEmail: teacherEmail, StudentEmail: studentEmail}
	})
	err = db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&registerRelationships).Error

	return err
}

func getAllStudentsRegisteredToTeacherEmails(teacherEmail string, db *gorm.DB) (studentEmails []string, err error) {
	var students []*Student

	query := db.Table("students").Select("students.email").Joins("left join register_relationships on students.email = register_relationships.student_email").
		Where(" register_relationships.teacher_email = ?", teacherEmail).Find(&students)
	err = query.Error

	studentEmails = helpers.Map(students, func(student *Student) string {
		return student.Email
	})
	return studentEmails, err
}

func getUnsuspendedStudentsRegisteredToTeacher(teacherEmail string, mentionedStudentEmails []string, db *gorm.DB) (studentEmails []string, err error) {
	var students []*Student

	err = db.Table("students").Select("students.email").Joins("left join register_relationships on students.email = register_relationships.student_email").
		Where("students.is_suspended = ?", false).Where(db.Where("register_relationships.teacher_email = ?", teacherEmail).Or("students.email in ?", mentionedStudentEmails)).
		Distinct("students.email").Find(&students).Error

	studentEmails = helpers.Map(students, func(student *Student) string {
		return student.Email
	})

	return studentEmails, err
}

func deleteAllRegisterRelationships(db *gorm.DB) error {
	return db.Exec("DELETE FROM register_relationships").Error
}