package types

type RegisterStudentsToTeacherRequest struct {
	TeacherEmail  string   `json:"teacher" binding:"required"`
	StudentEmails []string `json:"students" binding:"required"`
}

type StudentSuspensionReceiveRequest struct {
	StudentEmail string `json:"student" binding:"required"`
}

type RetrieveCommonStudentsRequest struct {
	TeacherEmail string `form:"teacher" binding:"required"`
}

type RetrieveStudentRecipientsRequest struct {
	TeacherEmail        string `json:"teacher" binding:"required"`
	NotificationMessage string `json:"notification" binding:"required"`
}
