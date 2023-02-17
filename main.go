package main

import (
	"github.com/gin-gonic/gin"
	"learning-management-system/controllers"
	"learning-management-system/database"
	"learning-management-system/models"
)

func main() {
	connection := setupDb()
	router := setupRouter(connection)
	_ = router.Run(":8080")
}

func setupRouter(connection *database.Connection) *gin.Engine {

	router := gin.Default()
	repository := controllers.NewRepository(connection)

	router.POST("/api/register", repository.RegisterStudentsToTeacher)
	router.GET("/api/commonstudents", repository.RetrieveCommonStudents)
	router.POST("/api/suspend", repository.SuspendStudent)
	router.POST("/api/retrievefornotifications", repository.RetrieveStudentRecipients)
	router.DELETE("/api/clear", repository.ClearDatabase)
	router.POST("/api/populateteachers", repository.PopulateTeachers)
	router.POST("/api/populatestudents", repository.PopulateStudents)

	return router
}

func setupDb() *database.Connection {
	connection := database.InitDefaultConnection()
	connection.GetDb().AutoMigrate(&models.Teacher{}, &models.Student{}, &models.RegisterRelationship{})
	return connection
}
