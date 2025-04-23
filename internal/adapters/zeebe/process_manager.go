package zeebe

import (
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

type ZeebeProcessManagerImpl struct {
	client zbc.Client
}

func NewZeebeProcessManager(client zbc.Client) *ZeebeProcessManagerImpl {
	return &ZeebeProcessManagerImpl{
		client: client,
	}
}

func (z *ZeebeProcessManagerImpl) StartSignupProcess(username, password string) error {
	defer func() {
		if r := recover(); r != nil {
			// Convert panic to error
			return
		}
	}()

	MustStartSignUpProcessInstance(z.client, username, password)
	return nil
}

func (z *ZeebeProcessManagerImpl) StartLoginProcess(username, password string) error {
	defer func() {
		if r := recover(); r != nil {
			// Convert panic to error
			return
		}
	}()

	MustStartLoginProcessInstance(z.client, username, password)
	return nil
}
