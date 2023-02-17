package controllers

import (
	"learning-management-system/database"
	"learning-management-system/helpers"
	"learning-management-system/models"
	"learning-management-system/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeacherRepository struct {
	connection *database.Connection
}

func NewTeacherRepo(connection *database.Connection) *TeacherRepository {
	return &TeacherRepository{connection: connection}
}

func (teacherRepository *TeacherRepository) RegisterStudentsToTeacher(context *gin.Context) {

	registerStudentsToTeacherRequest := &types.RegisterStudentsToTeacherRequest{}

	if contextErr := helpers.BindRegisterStudentsToTeacherRequest(context, registerStudentsToTeacherRequest); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailAddresses(append(registerStudentsToTeacherRequest.StudentEmails, registerStudentsToTeacherRequest.TeacherEmail)); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if err := models.RegisterStudentsToTeacher(registerStudentsToTeacherRequest.TeacherEmail, registerStudentsToTeacherRequest.StudentEmails, teacherRepository.connection); err != nil {
		generateInternalServerErrorResponse(context, err)
		return
	}

	context.JSON(http.StatusNoContent, nil)
}

func (teacherRepository *TeacherRepository) SuspendStudent(context *gin.Context) {
	studentSuspension := &types.StudentSuspensionRequest{}

	if contextErr := helpers.BindSuspendStudentRequest(context, studentSuspension); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailFormat(studentSuspension.StudentEmail); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if dbError := models.SuspendStudent(studentSuspension.StudentEmail, teacherRepository.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	}

	context.JSON(http.StatusNoContent, nil)
}

func (teacherRepository *TeacherRepository) RetrieveCommonStudents(context *gin.Context) {
	retrieveCommonStudentsRequest := &types.RetrieveCommonStudentsRequest{}

	if contextErr := helpers.BindRetrieveCommonStudents(context, retrieveCommonStudentsRequest); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailAddresses(retrieveCommonStudentsRequest.TeacherEmails); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	nonDuplicateTeacherEmails := helpers.RemoveDuplicatesInStringSlice(retrieveCommonStudentsRequest.TeacherEmails)

	studentEmails, dbErr := models.RetrieveCommonStudents(nonDuplicateTeacherEmails, teacherRepository.connection)
	if dbErr != nil {
		generateInternalServerErrorResponse(context, dbErr)
		return
	}

	context.JSON(http.StatusOK, &types.RetrieveRegisteredStudentsResponse{
		StudentEmails: studentEmails,
	})
}

func (teacherRepository *TeacherRepository) RetrieveStudentRecipients(context *gin.Context) {
	retrieveStudentRecipientsRequest := &types.RetrieveStudentRecipientsRequest{}

	if contextErr := helpers.BindRetrieveStudentRecipients(context, retrieveStudentRecipientsRequest); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailFormat(retrieveStudentRecipientsRequest.TeacherEmail); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	mentionedStudents := helpers.FindValidEmailsInText(retrieveStudentRecipientsRequest.NotificationMessage)
	recipientEmails, dbErr := models.RetrieveStudentRecipients(retrieveStudentRecipientsRequest.TeacherEmail, mentionedStudents, teacherRepository.connection)
	if dbErr != nil {
		generateInternalServerErrorResponse(context, dbErr)
		return
	}

	context.JSON(http.StatusOK, &types.RetrieveCommonStudentsResponse{
		StudentEmails: recipientEmails,
	})
}

func (teacherRepository *TeacherRepository) ClearDatabase(context *gin.Context) {
	if err := models.ClearDatabase(teacherRepository.connection); err != nil {
		generateInternalServerErrorResponse(context, err)
		return
	}
	context.JSON(http.StatusNoContent, nil)
}

func generateBadRequestErrorResponse(context *gin.Context, err error) {
	context.AbortWithStatusJSON(400, types.ErrorResponse{Message: err.Error()})
}

func generateInternalServerErrorResponse(context *gin.Context, err error) {
	context.AbortWithStatusJSON(500, types.ErrorResponse{Message: err.Error()})
}
