package main

import (
	"github.com/gin-gonic/gin"
	"learning-management-system/controllers"
	"learning-management-system/database"
	"learning-management-system/models"
	"log"
	"os/exec"
	"strconv"
)

func main() {

	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	// output has trailing \n
	// need to remove the \n
	// otherwise it will cause error for strconv.Atoi
	// log.Println(output[:len(output)-1])

	// 0 = root, 501 = non-root user
	i, err := strconv.Atoi(string(output[:len(output)-1]))

	if err != nil {
		log.Fatal(err)
	}

	if i == 0 {
		log.Println("Awesome! You are now running this program with root permissions!")
	} else {
		log.Fatal("This program must be run as root! (sudo)")
	}

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
