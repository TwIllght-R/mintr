package events

const (
	TopicEmailVerification = "email.verification"
)

type EmailVerificationEvent struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}
