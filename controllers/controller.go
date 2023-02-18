package controllers

import (
	"github.com/gin-gonic/gin"
	"learning-management-system/database"
	"learning-management-system/helpers"
	"learning-management-system/transaction_managers"
	"learning-management-system/types"
	"net/http"
)

type Controller struct {
	connection         *database.Connection
	transactionManager *transaction_managers.TransactionManager
}

func NewController(connection *database.Connection) *Controller {
	return &Controller{
		connection:         connection,
		transactionManager: transaction_managers.NewTransactionManager(),
	}
}

func (controller *Controller) RegisterStudentsToTeacher(context *gin.Context) {

	registerStudentsToTeacherRequest := &types.RegisterStudentsToTeacherRequest{}

	if contextErr := helpers.BindRegisterStudentsToTeacherRequest(context, registerStudentsToTeacherRequest); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailAddresses(append(registerStudentsToTeacherRequest.StudentEmails, registerStudentsToTeacherRequest.TeacherEmail)); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if userError, dbError := controller.transactionManager.ValidateTeachersExists([]string{registerStudentsToTeacherRequest.TeacherEmail}, controller.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	} else if userError != nil {
		generateBadRequestErrorResponse(context, userError)
		return
	}

	if userError, dbError := controller.transactionManager.ValidateStudentsExists(registerStudentsToTeacherRequest.StudentEmails, controller.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	} else if userError != nil {
		generateBadRequestErrorResponse(context, userError)
		return
	}

	if err := controller.transactionManager.RegisterStudentsToTeacher(registerStudentsToTeacherRequest.TeacherEmail, registerStudentsToTeacherRequest.StudentEmails, controller.connection); err != nil {
		generateInternalServerErrorResponse(context, err)
		return
	}

	context.JSON(http.StatusNoContent, nil)
}

func (controller *Controller) SuspendStudent(context *gin.Context) {
	studentSuspension := &types.StudentSuspensionRequest{}

	if contextErr := helpers.BindSuspendStudentRequest(context, studentSuspension); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailFormat(studentSuspension.StudentEmail); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if userError, dbError := controller.transactionManager.ValidateStudentsExists([]string{studentSuspension.StudentEmail}, controller.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	} else if userError != nil {
		generateBadRequestErrorResponse(context, userError)
		return
	}

	if dbError := controller.transactionManager.SuspendStudent(studentSuspension.StudentEmail, controller.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	}

	context.JSON(http.StatusNoContent, nil)
}

func (controller *Controller) RetrieveCommonStudents(context *gin.Context) {
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

	if userError, dbError := controller.transactionManager.ValidateTeachersExists(nonDuplicateTeacherEmails, controller.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	} else if userError != nil {
		generateBadRequestErrorResponse(context, userError)
		return
	}

	studentEmails, dbErr := controller.transactionManager.RetrieveCommonStudentEmails(nonDuplicateTeacherEmails, controller.connection)
	if dbErr != nil {
		generateInternalServerErrorResponse(context, dbErr)
		return
	}

	context.JSON(http.StatusOK, &types.RetrieveRegisteredStudentsResponse{
		StudentEmails: studentEmails,
	})
}

func (controller *Controller) RetrieveStudentRecipients(context *gin.Context) {
	retrieveStudentRecipientsRequest := &types.RetrieveStudentRecipientsRequest{}

	if contextErr := helpers.BindRetrieveStudentRecipientsRequest(context, retrieveStudentRecipientsRequest); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailFormat(retrieveStudentRecipientsRequest.TeacherEmail); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if userError, dbError := controller.transactionManager.ValidateTeachersExists([]string{retrieveStudentRecipientsRequest.TeacherEmail}, controller.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	} else if userError != nil {
		generateBadRequestErrorResponse(context, userError)
		return
	}

	mentionedStudents := helpers.RemoveDuplicatesInStringSlice(helpers.FindValidEmailsInText(retrieveStudentRecipientsRequest.NotificationMessage))

	if userError, dbError := controller.transactionManager.ValidateStudentsExists(mentionedStudents, controller.connection); dbError != nil {
		generateInternalServerErrorResponse(context, dbError)
		return
	} else if userError != nil {
		generateBadRequestErrorResponse(context, userError)
		return
	}

	recipientEmails, dbErr := controller.transactionManager.RetrieveStudentRecipients(retrieveStudentRecipientsRequest.TeacherEmail, mentionedStudents, controller.connection)
	if dbErr != nil {
		generateInternalServerErrorResponse(context, dbErr)
		return
	}

	context.JSON(http.StatusOK, &types.RetrieveCommonStudentsResponse{
		StudentEmails: recipientEmails,
	})
}

func (controller *Controller) ClearDatabase(context *gin.Context) {
	if err := controller.transactionManager.ClearDatabase(controller.connection); err != nil {
		generateInternalServerErrorResponse(context, err)
		return
	}
	context.JSON(http.StatusNoContent, nil)
}

func (controller *Controller) PopulateStudents(context *gin.Context) {

	populateStudentsRequest := &types.PopulateStudentsRequest{}

	if contextErr := helpers.BindPopulateStudentsRequest(context, populateStudentsRequest); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailAddresses(populateStudentsRequest.StudentEmails); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if err := controller.transactionManager.PopulateStudents(populateStudentsRequest.StudentEmails, controller.connection); err != nil {
		generateInternalServerErrorResponse(context, err)
		return
	}

	context.JSON(http.StatusNoContent, nil)
}

func (controller *Controller) PopulateTeachers(context *gin.Context) {
	populateTeachersRequest := &types.PopulateTeachersRequest{}

	if contextErr := helpers.BindPopulateTeachersRequest(context, populateTeachersRequest); contextErr != nil {
		generateBadRequestErrorResponse(context, contextErr)
		return
	}

	if validationErr := helpers.ValidateEmailAddresses(populateTeachersRequest.TeacherEmails); validationErr != nil {
		generateBadRequestErrorResponse(context, validationErr)
		return
	}

	if err := controller.transactionManager.PopulateTeachers(populateTeachersRequest.TeacherEmails, controller.connection); err != nil {
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
