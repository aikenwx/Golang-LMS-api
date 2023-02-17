package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"learning-management-system/types"
	"reflect"
)

func BindRegisterStudentsToTeacherRequest(context *gin.Context, registerStudentsToTeacherRequest *types.RegisterStudentsToTeacherRequest) error {

	if contentTypeErr := validateContentTypeIsApplicationJson(context); contentTypeErr != nil {
		return contentTypeErr
	}

	if ginErr := context.ShouldBindJSON(registerStudentsToTeacherRequest); ginErr != nil {
		return validateGinBindings(registerStudentsToTeacherRequest, "json", ginErr)
	}

	return nil
}

func BindSuspendStudentRequest(context *gin.Context, studentSuspensionRequest *types.StudentSuspensionRequest) error {

	if contentTypeErr := validateContentTypeIsApplicationJson(context); contentTypeErr != nil {
		return contentTypeErr
	}

	if ginErr := context.ShouldBindJSON(studentSuspensionRequest); ginErr != nil {
		return validateGinBindings(studentSuspensionRequest, "json", ginErr)
	}

	return nil
}

func BindRetrieveCommonStudents(context *gin.Context, retrieveCommonStudentsRequest *types.RetrieveCommonStudentsRequest) error {

	if ginErr := context.ShouldBindQuery(retrieveCommonStudentsRequest); ginErr != nil {
		return validateGinBindings(retrieveCommonStudentsRequest, "form", ginErr)
	}

	return nil
}

func BindRetrieveStudentRecipients(context *gin.Context, retrieveStudentRecipientsRequest *types.RetrieveStudentRecipientsRequest) error {

	if contentTypeErr := validateContentTypeIsApplicationJson(context); contentTypeErr != nil {
		return contentTypeErr
	}

	if ginErr := context.ShouldBindJSON(retrieveStudentRecipientsRequest); ginErr != nil {
		return validateGinBindings(retrieveStudentRecipientsRequest, "json", ginErr)
	}

	return nil
}

func validateContentTypeIsApplicationJson(context *gin.Context) error {
	headerContentType := context.GetHeader("Content-Type")

	if headerContentType == "" {
		return fmt.Errorf("Content-Type header of application/json must be provided")
	}

	if headerContentType != "application/json" {
		return fmt.Errorf("Content-Type header must be application/json")
	}
	return nil
}

func validateGinBindings[T any](requestStruct *T, binding string, errs ...error) error {
	var out []string
	for _, err := range errs {
		switch typedError := any(err).(type) {
		// We parse the following possible errors into more readable formats for users
		case validator.ValidationErrors:
			for _, e := range typedError {
				out = append(out, parseFieldError(e, requestStruct, binding))
			}
		case *json.UnmarshalTypeError:
			out = append(out, parseMarshallingError(*typedError))
		default:
			out = append(out, err.Error())
		}
	}

	// We only return the first error for simplification
	if len(out) > 0 {
		return fmt.Errorf(out[0])
	}

	return nil
}

func parseFieldError[T any](fieldError validator.FieldError, requestStruct *T, binding string) string {
	field, _ := reflect.TypeOf(requestStruct).Elem().FieldByName(fieldError.Field())
	fieldTagName, _ := field.Tag.Lookup(binding)

	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("The required field %s is not supplied", fieldTagName)
	default:
		return fmt.Errorf("%v", fieldError).Error()
	}
}

func parseMarshallingError(e json.UnmarshalTypeError) string {
	return fmt.Sprintf("The field %s must be a %s", e.Field, e.Type.String())
}
