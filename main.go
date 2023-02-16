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
	userRepo := controllers.NewTeacherRepo(connection)

	router.POST("/api/register", userRepo.RegisterStudentsToTeacher)
	router.GET("/api/commonstudents", userRepo.RetrieveCommonStudents)
	router.POST("/api/suspend", userRepo.SuspendStudent)
	router.POST("/api/retrievefornotifications", userRepo.RetrieveStudentRecipients)
	router.DELETE("/api/clear", userRepo.ClearDatabase)

	return router
}

func setupDb() *database.Connection {
	connection := database.InitDefaultConnection()
	connection.GetDb().AutoMigrate(&models.Teacher{}, &models.Student{}, &models.RegisterRelationship{})
	return connection
}
