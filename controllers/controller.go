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

	studentRegistration := types.RegisterStudentsToTeacherRequest{}

	if contextErr := context.ShouldBindJSON(&studentRegistration); contextErr != nil {

		generateBadRequestErrorResponse(context, helpers.GenerateInvalidRequestsError())
		return
	}

	if validationErr := helpers.ValidateEmailAddresses(append(studentRegistration.StudentEmails, studentRegistration.TeacherEmail)); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if err := models.RegisterStudentsToTeacher(studentRegistration.TeacherEmail, studentRegistration.StudentEmails, teacherRepository.connection); err != nil {
		generateInternalServerErrorResponse(context, err)
		return
	}

	context.JSON(http.StatusNoContent, nil)
}

func (teacherRepository *TeacherRepository) SuspendStudent(context *gin.Context) {
	studentSuspension := &types.StudentSuspensionReceiveRequest{}

	if err := context.ShouldBindJSON(studentSuspension); err != nil {
		generateBadRequestErrorResponse(context, helpers.GenerateInvalidRequestsError())
		return
	}

	if err := helpers.ValidateEmailFormat(studentSuspension.StudentEmail); err != nil {
		generateBadRequestErrorResponse(context, err)
		return
	}

	if err := models.SuspendStudent(studentSuspension.StudentEmail, teacherRepository.connection); err != nil {
		generateInternalServerErrorResponse(context, err)
		return
	}

	context.JSON(http.StatusNoContent, nil)
}

func (teacherRepository *TeacherRepository) RetrieveCommonStudents(context *gin.Context) {
	registeredStudentRetrieval := &types.RetrieveCommonStudentsRequest{}

	if err := context.ShouldBindQuery(registeredStudentRetrieval); err != nil {
		generateBadRequestErrorResponse(context, helpers.GenerateInvalidRequestsError())
		return
	}

	if err := helpers.ValidateEmailFormat(registeredStudentRetrieval.TeacherEmail); err != nil {
		generateBadRequestErrorResponse(context, err)
		return
	}

	studentEmails, dbErr := models.RetrieveCommonStudents(registeredStudentRetrieval.TeacherEmail, teacherRepository.connection)
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

	if contextErr := context.ShouldBindJSON(retrieveStudentRecipientsRequest); contextErr != nil {
		generateBadRequestErrorResponse(context, helpers.GenerateInvalidRequestsError())
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
