package controllers

import (
	"github.com/gin-gonic/gin"
	"learning-management-system/database"
	"learning-management-system/helpers"
	"learning-management-system/models"
	"learning-management-system/types"
	"net/http"
)

type Repository struct {
	connection *database.Connection
}

func NewRepository(connection *database.Connection) *Repository {
	return &Repository{connection: connection}
}

func (repository *Repository) RegisterStudentsToTeacher(context *gin.Context) {

	registerStudentsToTeacherRequest := &types.RegisterStudentsToTeacherRequest{}

	if contextErr := helpers.BindRegisterStudentsToTeacherRequest(context, registerStudentsToTeacherRequest); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailAddresses(append(registerStudentsToTeacherRequest.StudentEmails, registerStudentsToTeacherRequest.TeacherEmail)); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if userError, dbError := models.ValidateTeachersExists([]string{registerStudentsToTeacherRequest.TeacherEmail}, repository.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	} else if userError != nil {
		generateBadRequestErrorResponse(context, userError)
		return
	}

	if userError, dbError := models.ValidateStudentsExists(registerStudentsToTeacherRequest.StudentEmails, repository.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	} else if userError != nil {
		generateBadRequestErrorResponse(context, userError)
		return
	}

	if err := models.RegisterStudentsToTeacher(registerStudentsToTeacherRequest.TeacherEmail, registerStudentsToTeacherRequest.StudentEmails, repository.connection); err != nil {
		generateInternalServerErrorResponse(context, err)
		return
	}

	context.JSON(http.StatusNoContent, nil)
}

func (repository *Repository) SuspendStudent(context *gin.Context) {
	studentSuspension := &types.StudentSuspensionRequest{}

	if contextErr := helpers.BindSuspendStudentRequest(context, studentSuspension); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailFormat(studentSuspension.StudentEmail); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if userError, dbError := models.ValidateStudentsExists([]string{studentSuspension.StudentEmail}, repository.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	} else if userError != nil {
		generateBadRequestErrorResponse(context, userError)
		return
	}

	if dbError := models.SuspendStudent(studentSuspension.StudentEmail, repository.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	}

	context.JSON(http.StatusNoContent, nil)
}

func (repository *Repository) RetrieveCommonStudents(context *gin.Context) {
	retrieveCommonStudentsRequest := &types.RetrieveCommonStudentsRequest{}

	if contextErr := helpers.BindRetrieveCommonStudentsRequest(context, retrieveCommonStudentsRequest); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	nonDuplicateTeacherEmails := helpers.RemoveDuplicatesInStringSlice(retrieveCommonStudentsRequest.TeacherEmails)

	if validationErr := helpers.ValidateEmailAddresses(nonDuplicateTeacherEmails); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if userError, dbError := models.ValidateTeachersExists(nonDuplicateTeacherEmails, repository.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	} else if userError != nil {
		generateBadRequestErrorResponse(context, userError)
		return
	}

	studentEmails, dbErr := models.RetrieveCommonStudents(nonDuplicateTeacherEmails, repository.connection)
	if dbErr != nil {
		generateInternalServerErrorResponse(context, dbErr)
		return
	}

	context.JSON(http.StatusOK, &types.RetrieveRegisteredStudentsResponse{
		StudentEmails: studentEmails,
	})
}

func (repository *Repository) RetrieveStudentRecipients(context *gin.Context) {
	retrieveStudentRecipientsRequest := &types.RetrieveStudentRecipientsRequest{}

	if contextErr := helpers.BindRetrieveStudentRecipientsRequest(context, retrieveStudentRecipientsRequest); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailFormat(retrieveStudentRecipientsRequest.TeacherEmail); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if userError, dbError := models.ValidateTeachersExists([]string{retrieveStudentRecipientsRequest.TeacherEmail}, repository.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	} else if userError != nil {
		generateBadRequestErrorResponse(context, userError)
		return
	}

	mentionedStudents := helpers.RemoveDuplicatesInStringSlice(helpers.FindValidEmailsInText(retrieveStudentRecipientsRequest.NotificationMessage))

	if userError, dbError := models.ValidateStudentsExists(mentionedStudents, repository.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	} else if userError != nil {
		generateBadRequestErrorResponse(context, userError)
		return
	}

	recipientEmails, dbErr := models.RetrieveStudentRecipients(retrieveStudentRecipientsRequest.TeacherEmail, mentionedStudents, repository.connection)
	if dbErr != nil {
		generateInternalServerErrorResponse(context, dbErr)
		return
	}

	context.JSON(http.StatusOK, &types.RetrieveCommonStudentsResponse{
		StudentEmails: recipientEmails,
	})
}

func (repository *Repository) ClearDatabase(context *gin.Context) {
	if err := models.ClearDatabase(repository.connection); err != nil {
		generateInternalServerErrorResponse(context, err)
		return
	}
	context.JSON(http.StatusNoContent, nil)
}

func (repository *Repository) PopulateStudents(context *gin.Context) {

	populateStudentsRequest := &types.PopulateStudentsRequest{}

	if contextErr := helpers.BindPopulateStudentsRequest(context, populateStudentsRequest); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailAddresses(populateStudentsRequest.StudentEmails); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if err := models.PopulateStudents(populateStudentsRequest.StudentEmails, repository.connection); err != nil {
		generateInternalServerErrorResponse(context, err)
		return
	}

	context.JSON(http.StatusNoContent, nil)
}

func (repository *Repository) PopulateTeachers(context *gin.Context) {
	populateTeachersRequest := &types.PopulateTeachersRequest{}

	if contextErr := helpers.BindPopulateTeachersRequest(context, populateTeachersRequest); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailAddresses(populateTeachersRequest.TeacherEmails); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if err := models.PopulateTeachers(populateTeachersRequest.TeacherEmails, repository.connection); err != nil {
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
