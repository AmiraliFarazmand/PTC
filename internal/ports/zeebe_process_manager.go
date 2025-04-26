package ports

type ZeebeProcessManager interface {
	StartSignupProcess(username, password string) error
	StartLoginProcess(username, password string) error
}
