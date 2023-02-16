package types

type RetrieveRegisteredStudentsResponse struct {
	StudentEmails []string `json:"students"`
}

type RetrieveCommonStudentsResponse struct {
	StudentEmails []string `json:"recipients"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
