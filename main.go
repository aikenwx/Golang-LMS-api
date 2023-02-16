package main

import (
	"github.com/gin-gonic/gin"
	"learning-management-system/controllers"
	"learning-management-system/database"
	"learning-management-system/models"
)

func main() {
	database.InitGlobalDbConnection()
	connection := database.GlobalConnection
	connection.GetDb().AutoMigrate(&models.Teacher{}, &models.Student{}, &models.RegisterRelationship{})

	router := setupRouter()
	_ = router.Run(":8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	userRepo := controllers.NewTeacherRepo()

	router.POST("/api/register", userRepo.RegisterStudentsToTeacher)
	router.GET("/api/commonstudents", userRepo.RetrieveCommonStudents)
	router.POST("/api/suspend", userRepo.SuspendStudent)
	router.POST("/api/retrievefornotifications", userRepo.RetrieveStudentRecipients)

	return router
}
