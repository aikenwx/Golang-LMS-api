package tests

import (
	"learning-management-system/models"
	"learning-management-system/repositories"
	"testing"
)

func TestCreateAndRetrieveStudent(t *testing.T) {
	connection, _ := SetUpTestDb()
	db := connection.GetDb()
	studentRepo := repositories.NewStudentRepo()
	studentRepo.CreateStudentsIfNotExist([]string{"test1@gmail.com"}, db)
	student, err := studentRepo.GetStudentByEmail("test1@gmail.com", db)
	if err != nil {
		t.Errorf("Error throwned: %s", err.Error())
	}
	assertEquals(t, &models.Student{Email: "test1@gmail.com"}, student)
}

func TestCreateAndRetrieveTeacher(t *testing.T) {
	connection, _ := SetUpTestDb()
	db := connection.GetDb()
	teacherRepo := repositories.NewTeacherRepo()
	teacherRepo.CreateTeachersIfNotExist([]string{"test1@gmail.com"}, db)
	teacher, err := teacherRepo.GetTeacherByEmail("test1@gmail.com", db)
	if err != nil {
		t.Errorf("Error throwned: %s", err.Error())
	}
	assertEquals(t, &models.Teacher{Email: "test1@gmail.com"}, teacher)
}

func TestUpdateStudent(t *testing.T) {
	connection, _ := SetUpTestDb()
	db := connection.GetDb()
	studentRepo := repositories.NewStudentRepo()
	studentRepo.CreateStudentsIfNotExist([]string{"test1@gmail.com"}, db)
	student := &models.Student{Email: "test1@gmail.com", IsSuspended: true}
	studentRepo.UpdateStudent(student, db)
	updatedStudent, err := studentRepo.GetStudentByEmail("test1@gmail.com", db)
	if err != nil {
		t.Errorf("Error throwned: %s", err.Error())
	}
	expectedStudent := &models.Student{Email: "test1@gmail.com", IsSuspended: true}
	assertEquals(t, expectedStudent, updatedStudent)
}

func TestCreateAndRetrieveMultipleStudents(t *testing.T) {
	connection, _ := SetUpTestDb()
	db := connection.GetDb()
	studentRepo := repositories.NewStudentRepo()
	studentRepo.CreateStudentsIfNotExist([]string{"test1@gmail.com", "test2@gmail.com"}, db)
	students, err := studentRepo.GetStudentsByEmails([]string{"test1@gmail.com", "test2@gmail.com"}, db)
	if err != nil {
		t.Errorf("Error throwned: %s", err.Error())
	}

	expectedStudents := []*models.Student{
		{Email: "test1@gmail.com", IsSuspended: false}, {Email: "test2@gmail.com", IsSuspended: false},
	}

	assertEquals(t, expectedStudents, students)
}

func TestCreateAndRetrieveMultipleTeachers(t *testing.T) {
	connection, _ := SetUpTestDb()
	db := connection.GetDb()
	teacherRepo := repositories.NewTeacherRepo()
	teacherRepo.CreateTeachersIfNotExist([]string{"test1@gmail.com", "test2@gmail.com"}, db)
	teachers, err := teacherRepo.GetTeachersByEmails([]string{"test1@gmail.com", "test2@gmail.com"}, db)
	if err != nil {
		t.Errorf("Error throwned: %s", err.Error())
	}

	expectedTeachers := []*models.Teacher{
		{Email: "test1@gmail.com"}, {Email: "test2@gmail.com"},
	}

	assertEquals(t, expectedTeachers, teachers)
}

func TestDeleteAllStudents(t *testing.T) {
	connection, _ := SetUpTestDb()
	db := connection.GetDb()
	studentRepo := repositories.NewStudentRepo()
	studentRepo.CreateStudentsIfNotExist([]string{"test1@gmail.com", "test2@gmail.com"}, db)
	studentRepo.DeleteAllStudents(db)
	students, err := studentRepo.GetStudentsByEmails([]string{"test1@gmail.com", "test2@gmail.com"}, db)
	if err != nil {
		t.Errorf("Error throwned: %s", err.Error())
	}

	assertEquals(t, 0, len(students))
}

func TestDeleteAllTeachers(t *testing.T) {
	connection, _ := SetUpTestDb()
	db := connection.GetDb()
	teacherRepo := repositories.NewTeacherRepo()
	teacherRepo.CreateTeachersIfNotExist([]string{"test1@gmail.com", "test2@gmail.com"}, db)
	teacherRepo.DeleteAllTeachers(db)
	teachers, err := teacherRepo.GetTeachersByEmails([]string{"test1@gmail.com", "test2@gmail.com"}, db)
	if err != nil {
		t.Errorf("Error throwned: %s", err.Error())
	}

	assertEquals(t, 0, len(teachers))
}

func TestCreateAndRetrieveRegisterRelation(t *testing.T) {
	connection, _ := SetUpTestDb()
	db := connection.GetDb()
	studentRepo := repositories.NewStudentRepo()
	teacherRepo := repositories.NewTeacherRepo()
	relationshipRepo := repositories.NewRegisterRelationshipRepo()
	studentRepo.CreateStudentsIfNotExist([]string{"test1@gmail.com", "test2@gmail.com"}, db)
	teacherRepo.CreateTeachersIfNotExist([]string{"test3@gmail.com"}, db)

	relationshipRepo.CreateOneTeacherToManyStudentsRegisterRelationshipsIfNotExists("test3@gmail.com", []string{"test1@gmail.com", "test2@gmail.com"}, db)

	relationships, err := relationshipRepo.GetRelationshipsByTeacherEmail("test3@gmail.com", db)

	if err != nil {
		t.Errorf("Error throwned: %s", err.Error())
	}

	expectedRelationships := []*models.RegisterRelationship{
		{TeacherEmail: "test3@gmail.com", StudentEmail: "test1@gmail.com", RegisteredStudent: models.Student{Email: "test1@gmail.com", IsSuspended: false}, Teacher: models.Teacher{Email: "test3@gmail.com"}},
		{TeacherEmail: "test3@gmail.com", StudentEmail: "test2@gmail.com", RegisteredStudent: models.Student{Email: "test2@gmail.com", IsSuspended: false}, Teacher: models.Teacher{Email: "test3@gmail.com"}},
	}

	assertEquals(t, expectedRelationships, relationships)
}

func TestGetRegisterRelationshipsByTeacherEmails(t *testing.T) {
	connection, _ := SetUpTestDb()
	db := connection.GetDb()
	studentRepo := repositories.NewStudentRepo()
	teacherRepo := repositories.NewTeacherRepo()
	relationshipRepo := repositories.NewRegisterRelationshipRepo()
	studentRepo.CreateStudentsIfNotExist([]string{"test1@gmail.com", "test2@gmail.com", "test3@gmail.com"}, db)
	teacherRepo.CreateTeachersIfNotExist([]string{"test4@gmail.com", "test5@gmail.com"}, db)

	relationshipRepo.CreateOneTeacherToManyStudentsRegisterRelationshipsIfNotExists("test4@gmail.com", []string{"test1@gmail.com", "test2@gmail.com"}, db)

	relationshipRepo.CreateOneTeacherToManyStudentsRegisterRelationshipsIfNotExists("test5@gmail.com", []string{"test1@gmail.com", "test3@gmail.com"}, db)

	relationships, err := relationshipRepo.GetRelationshipsByTeacherEmails([]string{"test4@gmail.com", "test5@gmail.com"}, db)

	if err != nil {
		t.Errorf("Error throwned: %s", err.Error())
	}

	expectedRelationships := []*models.RegisterRelationship{
		{TeacherEmail: "test4@gmail.com", StudentEmail: "test1@gmail.com", RegisteredStudent: models.Student{Email: "test1@gmail.com", IsSuspended: false}, Teacher: models.Teacher{Email: "test4@gmail.com"}},
		{TeacherEmail: "test5@gmail.com", StudentEmail: "test1@gmail.com", RegisteredStudent: models.Student{Email: "test1@gmail.com", IsSuspended: false}, Teacher: models.Teacher{Email: "test5@gmail.com"}},
		{TeacherEmail: "test4@gmail.com", StudentEmail: "test2@gmail.com", RegisteredStudent: models.Student{Email: "test2@gmail.com", IsSuspended: false}, Teacher: models.Teacher{Email: "test4@gmail.com"}},
		{TeacherEmail: "test5@gmail.com", StudentEmail: "test3@gmail.com", RegisteredStudent: models.Student{Email: "test3@gmail.com", IsSuspended: false}, Teacher: models.Teacher{Email: "test5@gmail.com"}},
	}
	assertEquals(t, expectedRelationships, relationships)
}

func TestDeleteAllRegisterRelationships(t *testing.T) {
	connection, _ := SetUpTestDb()
	db := connection.GetDb()
	studentRepo := repositories.NewStudentRepo()
	teacherRepo := repositories.NewTeacherRepo()
	relationshipRepo := repositories.NewRegisterRelationshipRepo()
	studentRepo.CreateStudentsIfNotExist([]string{"test1@gmail.com", "test2@gmail.com", "test3@gmail.com"}, db)
	teacherRepo.CreateTeachersIfNotExist([]string{"test4@gmail.com", "test5@gmail.com"}, db)

	relationshipRepo.CreateOneTeacherToManyStudentsRegisterRelationshipsIfNotExists("test4@gmail.com", []string{"test1@gmail.com", "test2@gmail.com"}, db)

	relationshipRepo.CreateOneTeacherToManyStudentsRegisterRelationshipsIfNotExists("test5@gmail.com", []string{"test1@gmail.com", "test3@gmail.com"}, db)

	if err := relationshipRepo.DeleteAllRegisterRelationships(db); err != nil {
		t.Errorf("Error throwned: %s", err.Error())
	}

	relationships, err := relationshipRepo.GetRelationshipsByTeacherEmails([]string{"test4@gmail.com", "test5@gmail.com"}, db)

	if err != nil {
		t.Errorf("Error throwned: %s", err.Error())
	}

	assertEquals(t, 0, len(relationships))
}
