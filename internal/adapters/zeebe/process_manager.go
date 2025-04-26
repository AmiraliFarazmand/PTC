package zeebe

import (
	"fmt"

	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
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
	// defer func() {  // TODO: see if it needs recover function or not
	// 	if r := recover(); r != nil {
	// 		// Convert panic to error
	// 		return
	// 	}
	// }()

	

	result, err := StartSignUpProcessInstanceWithResult(z.client, username, password)
	if err != nil {
		return fmt.Errorf("failed to start signup process: %w", err)
	}

	if result.Error != "" {
		return fmt.Errorf(result.Error)
	}

	return nil
}

func (z *ZeebeProcessManagerImpl) StartLoginProcess(username, password string) (*domain.ProcessVariables, error) {
	result, err := StartLoginProcessInstanceWithResult(z.client, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to process login: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf(result.Error)
	}

	return result, nil
}
