package tests

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"learning-management-system/controllers"
	"learning-management-system/database"
	"learning-management-system/models"
	"learning-management-system/transaction_managers"
	"net/http"
	"reflect"
	"testing"
)

const TEST_USERNAME string = "root"
const TEST_PASSWORD string = ""
const TEST_HOST string = "localhost"
const TEST_NAME string = "test_lms_db"
const TEST_PORT string = "3306"
const TEST_DEBUG bool = false

const SERVER_PORT = "8090"

func SetUpTestDb() (*database.Connection, *controllers.Controller) {

	connection := database.NewConnection(&database.Credentials{
		Username: TEST_USERNAME,
		Password: TEST_PASSWORD,
		Host:     TEST_HOST,
		Name:     TEST_NAME,
		Port:     TEST_PORT,
		Debug:    TEST_DEBUG,
	})

	connection.GetDb().AutoMigrate(&models.Teacher{}, &models.Student{}, &models.RegisterRelationship{})
	router := gin.Default()
	controller := controllers.NewController(connection)

	transactionManager := transaction_managers.NewTransactionManager()
	transactionManager.ClearDatabase(connection)

	router.POST("/api/register", controller.RegisterStudentsToTeacher)
	router.GET("/api/commonstudents", controller.RetrieveCommonStudents)
	router.POST("/api/suspend", controller.SuspendStudent)
	router.POST("/api/retrievefornotifications", controller.RetrieveStudentRecipients)
	router.DELETE("/api/clear", controller.ClearDatabase)
	router.POST("/api/populateteachers", controller.PopulateTeachers)
	router.POST("/api/populatestudents", controller.PopulateStudents)
	go router.Run(":" + SERVER_PORT)

	return connection, controller
}

func testPost(jsonString string, relativePath string, expectedStatusCode int, expectedBody string, t *testing.T) {
	var jsonData = []byte(jsonString)
	responseBody := bytes.NewBuffer(jsonData)
	resp, err := http.Post("http://"+TEST_HOST+":"+SERVER_PORT+relativePath, "application/json", responseBody)

	if err != nil {
		t.Errorf(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatusCode {
		t.Errorf("wrong response code")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf(err.Error())
	}

	if string(body) != expectedBody {
		t.Errorf("wrong response body expected: " + expectedBody + " got: " + string(body))
	}

}

func testPostWithContentTypeHeaderHeader(jsonString string, relativePath string, expectedStatusCode int, expectedBody string, contentType string, t *testing.T) {
	var jsonData = []byte(jsonString)
	responseBody := bytes.NewBuffer(jsonData)
	resp, err := http.Post("http://"+TEST_HOST+":"+SERVER_PORT+relativePath, contentType, responseBody)

	if err != nil {
		t.Errorf(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatusCode {
		t.Errorf("wrong response code")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf(err.Error())
	}

	if string(body) != expectedBody {
		t.Errorf("wrong response body expected: " + expectedBody + " got: " + string(body))
	}

}

func testGet(relativePathWithParams string, expectedStatusCode int, expectedBody string, t *testing.T) {
	resp, err := http.Get("http://" + TEST_HOST + ":" + SERVER_PORT + relativePathWithParams)

	if err != nil {
		t.Errorf(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatusCode {
		t.Errorf("wrong response code")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf(err.Error())
	}

	if string(body) != expectedBody {

		t.Errorf("wrong response body expected: " + expectedBody + " got: " + string(body))

	}
}

func testDelete(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://"+TEST_HOST+":"+SERVER_PORT+"/api/clear", nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		t.Errorf("wrong response code")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf(err.Error())
	}

	if string(body) != "" {
		t.Errorf("wrong response body expected empty string, but got: " + string(body))
	}
}

func assertEquals[T any](t *testing.T, expected T, actual T) {
	// we use reflect.DeepEqual to compare the two values
	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Printf("Expected: %v, Actual: %v", expected, actual))
	}
}
