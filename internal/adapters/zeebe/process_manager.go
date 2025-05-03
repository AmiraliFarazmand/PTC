package zeebe

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

type ZeebeProcessManagerImpl struct {
	client zbc.Client  //TODO: in harekat roo purchase bezan
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

// TODO: instance saakhtane in mesle baghie beshe.
func (z *ZeebeProcessManagerImpl) StartPurchaseProcess(userID string, amount int, address string) (*domain.PurchaseProcessVariables, error) {
	variables := domain.PurchaseProcessVariables{
		UserID:  userID,
		Amount:  amount,
		Address: address,
		IsValid: true,
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	command, err := z.client.NewCreateInstanceCommand().
		BPMNProcessId("PurchaseProcess").
		LatestVersion().
		VariablesFromObject(variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create command: %w", err)
	}

	result, err := command.
		WithResult().
		FetchVariables("purchase_id", "isValid", "error").
		Send(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create and complete purchase instance: %w", err)
	}

	var resultVars domain.PurchaseProcessVariables
	if err := json.Unmarshal([]byte(result.GetVariables()), &resultVars); err != nil {
		return nil, fmt.Errorf("failed to parse result variables: %w", err)
	}
	return &resultVars, nil
}
