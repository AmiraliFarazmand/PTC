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

	result, err := StartSignUpProcessInstanceWithResult(z.client, username, password)
	if err != nil {
		return fmt.Errorf("failed to start signup process: %w", err)
	}

	if result.Error != "" {
		return fmt.Errorf(result.Error)
	}

	return nil
}

func (z *ZeebeProcessManagerImpl) StartLoginProcess(username, password string) (*domain.AuthProcessVariables, error) {
	result, err := StartLoginProcessInstanceWithResult(z.client, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to process login: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf(result.Error)
	}

	return result, nil
}

func (z *ZeebeProcessManagerImpl) StartPurchaseProcess(userID string, amount int, address string) (*domain.PurchaseProcessVariables, error) {
	result, err := StartPurchaseProcessWithResult(z.client, userID, address, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to create purchase: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf(result.Error)
	}

	return result, nil

}
