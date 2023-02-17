package system_test

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"learning-management-system/controllers"
	"learning-management-system/database"
	"learning-management-system/models"
	"net/http"
	"testing"
)

const TEST_USERNAME string = "root"
const TEST_PASSWORD string = ""
const TEST_HOST string = "localhost"
const TEST_NAME string = "test_lms_db"
const TEST_PORT string = "3306"
const TEST_DEBUG bool = false

const SERVER_PORT = "8090"

func TestCase1(t *testing.T) {
	setUpTestDb()
	testPost(`{"students": ["test1@gmail.com", "test2@gmail.com", "test3@gmail.com", "test4@gmail.com", "test5@gmail.com"]}`, "/api/populatestudents", 204, "", t)
	testPost(`{"teachers": ["test@gmail.com", "nani@gmail.com"]}`, "/api/populateteachers", 204, "", t)
	testPost(`{"teacher": "test@gmail.com", "students":["test1@gmail.com","test2@gmail.com"]}`, "/api/register", 204, "", t)
	testPost(`{}`, "/api/register", 400, `{"message":"The required field teacher is not supplied"}`, t)
	testGet("/api/commonstudents?teacher=test%40gmail.com", 200, `{"students":["test1@gmail.com","test2@gmail.com"]}`, t)
	testGet("/api/commonstudents?teacher=nani%40gmail.com", 200, `{"students":[]}`, t)
	testPost(`{"student":"test1@gmail.com"}`, "/api/suspend", 204, "", t)
	testPost(`{"student":"test3@gmail.com"}`, "/api/suspend", 204, "", t)
	testPost(`{"teacher":"test@gmail.com", "notification":"hello @test3@gmail.com and @test4@gmail.com and @test5@gmail.com"}`, "/api/retrievefornotifications", 200, `{"recipients":["test2@gmail.com","test4@gmail.com","test5@gmail.com"]}`, t)
	testPost(`{"teacher":"test@gmail.com", "notification":"hello @test2@gmail.com @test3@gmail.com and @test4@gmail.com and @test5@gmail.com"}`, "/api/retrievefornotifications", 200, `{"recipients":["test2@gmail.com","test4@gmail.com","test5@gmail.com"]}`, t)
	testPost(`{"teacher": "tes8gmail.com", "students":["test7gmail.com","test6gmail.com"]}`, "/api/register", 400, `{"message":"The email address test7gmail.com has an invalid format"}`, t)
	testPost(`{"teacher": [], "students":"test1@gmail.com"}`, "/api/register", 400, `{"message":"The field teacher must be a string"}`, t)
	testDelete(t)
	testGet("/api/commonstudents", 400, `{"message":"The required field teacher is not supplied"}`, t)
}

func TestCase2(t *testing.T) {
	setUpTestDb()
	testPost(`{"students": ["test1@gmail.com", "test2@gmail.com", "test3@gmail.com"]}`, "/api/populatestudents", 204, "", t)
	testPost(`{"teachers": ["test@gmail.com", "test7@gmail.com"]}`, "/api/populateteachers", 204, "", t)
	testGet("/api/commonstudents?teacher=test%40gmail.com&teacher=test7%40gmail.com", 200, `{"students":[]}`, t)
	testPost(`{"teacher": "test@gmail.com"}`, "/api/register", 400, `{"message":"The required field students is not supplied"}`, t)
	testPost(`{"teacher": "test@gmail.com", "students": "test2@gmail.com"}`, "/api/register", 400, `{"message":"The field students must be a []string"}`, t)
	testGet("/api/commonstudents?teacher=test%40gmail.com", 200, `{"students":[]}`, t)
	testPostWithContentTypeHeaderHeader(`{"message":"Content-Type header must be application/json"}`, "/api/register", 400, `{"message":"Content-Type header must be application/json"}`, "text/html; charset=UTF-8", t)
	testPost(`{"teacher": "test@gmail.com", "students":["test1@gmail.com","test2@gmail.com"]}`, "/api/register", 204, "", t)
	testPost(`{"teacher": "test7@gmail.com", "students":["test1@gmail.com","test3@gmail.com"]}`, "/api/register", 204, "", t)
	testPost(`{"teacher": "test7@gmail.com", "students":["test1@gmail.com","test3@gmail.com"]}`, "/api/register", 204, "", t)
	testGet("/api/commonstudents?teacher=test%40gmail.com&teacher=test7%40gmail.com", 200, `{"students":["test1@gmail.com"]}`, t)
	testGet("/api/commonstudents?teacher=test%40gmail.com&teacher=test7%40gmail.com&teacher=test7%40gmail.com", 200, `{"students":["test1@gmail.com"]}`, t)
	testDelete(t)
	testGet("/api/commonstudents?teacher=test%40gmail.com&teacher=test7%40gmail.com", 400, `{"message":"Teachers with emails test@gmail.com, test7@gmail.com do not exist in the database"}`, t)
	testGet("/api/commonstudents?teacher=testgmail.com&teacher=test3%40gmail.com", 400, `{"message":"The email address testgmail.com has an invalid format"}`, t)
}

func TestCase3(t *testing.T) {
	setUpTestDb()
	testPost(`{"student":"test3@gmail.com"}`, "/api/suspend", 400, `{"message":"Student with email test3@gmail.com does not exist in the database"}`, t)
	testPost(`{"teacher": "test@gmail.com", "students":["test1@gmail.com","test2@gmail.com"]}`, "/api/register", 400, `{"message":"Teacher with email test@gmail.com does not exist in the database"}`, t)
	testPost(`{"teachers": ["test@gmail.com", "nani@gmail.com"]}`, "/api/populateteachers", 204, "", t)
	testPost(`{"teacher": "test@gmail.com", "students":["test1@gmail.com","test2@gmail.com"]}`, "/api/register", 400, `{"message":"Students with emails test1@gmail.com, test2@gmail.com do not exist in the database"}`, t)
	testPost(`{"students": ["test1@gmail.com", "test2@gmail.com", "test3@gmail.com", "test4@gmail.com", "test5@gmail.com"]}`, "/api/populatestudents", 204, "", t)
	testPost(`{"students": ["test1@gmail.com", "test2@gmail.com", "test3@gmail.com", "test4@gmail.com", "test5@gmail.com"]}`, "/api/populatestudents", 204, "", t)
	testPost(`{"teachers": ["test@gmail.com", "nani@gmail.com"]}`, "/api/populateteachers", 204, "", t)
	testPost(`{}`, "/api/register", 400, `{"message":"The required field teacher is not supplied"}`, t)
	testGet("/api/commonstudents?teacher=test%40gmail.com", 200, `{"students":[]}`, t)
	testPost(`{"student":"test1@gmail.com"}`, "/api/suspend", 204, "", t)
	testPost(`{"student":"test3@gmail.com"}`, "/api/suspend", 204, "", t)
	testPost(`{"teacher":"test@gmail.com", "notification":"hello @test3@gmail.com and @test4@gmail.com and @test5@gmail.com"}`, "/api/retrievefornotifications", 200, `{"recipients":["test4@gmail.com","test5@gmail.com"]}`, t)
	testPost(`{"teacher":"test@gmail.com", "notification":"hello @test2@gmail.com @test3@gmail.com and @test4@gmail.com and @test5@gmail.com"}`, "/api/retrievefornotifications", 200, `{"recipients":["test2@gmail.com","test4@gmail.com","test5@gmail.com"]}`, t)
	testPost(`{"teacher": "tes8gmail.com", "students":["test7gmail.com","test6gmail.com"]}`, "/api/register", 400, `{"message":"The email address test7gmail.com has an invalid format"}`, t)
	testPost(`{"teacher": [], "students":"test1@gmail.com"}`, "/api/register", 400, `{"message":"The field teacher must be a string"}`, t)
	testDelete(t)
	testGet("/api/commonstudents", 400, `{"message":"The required field teacher is not supplied"}`, t)
}

func setUpTestDb() *database.Connection {

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
	repo := controllers.NewRepository(connection)

	models.ClearDatabase(connection)

	router.POST("/api/register", repo.RegisterStudentsToTeacher)
	router.GET("/api/commonstudents", repo.RetrieveCommonStudents)
	router.POST("/api/suspend", repo.SuspendStudent)
	router.POST("/api/retrievefornotifications", repo.RetrieveStudentRecipients)
	router.DELETE("/api/clear", repo.ClearDatabase)
	router.POST("/api/populateteachers", repo.PopulateTeachers)
	router.POST("/api/populatestudents", repo.PopulateStudents)
	go router.Run(":" + SERVER_PORT)

	return connection
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
