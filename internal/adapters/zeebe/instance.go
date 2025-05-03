package zeebe

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

func StartSignUpProcessInstanceWithResult(client zbc.Client, username, password string) (*domain.ProcessVariables, error) {
	variables := domain.ProcessVariables{
		Username: username,
		Password: password,
		IsValid:  true,
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	// Create instance with variables and wait for result
	command, err := client.NewCreateInstanceCommand().
		BPMNProcessId("SignupProcess").
		LatestVersion().
		VariablesFromObject(variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create command: %w", err)
	}

	result, err := command.
		WithResult().
		FetchVariables("username", "isValid", "error").
		Send(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create and complete signup instance: %w", err)
	}

	// Parse result variables
	var resultVars domain.ProcessVariables
	if err := json.Unmarshal([]byte(result.GetVariables()), &resultVars); err != nil {
		return nil, fmt.Errorf("failed to parse result variables: %w", err)
	}

	return &resultVars, nil
}

func StartLoginProcessInstanceWithResult(client zbc.Client, username, password string) (*domain.ProcessVariables, error) {
	variables := domain.ProcessVariables{
		Username:   username,
		Password:   password,
		LoginValid: true,
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	command, err := client.NewCreateInstanceCommand().
		BPMNProcessId("LoginProcess").
		LatestVersion().
		VariablesFromObject(variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create command: %w", err)
	}

	result, err := command.
		WithResult().
		FetchVariables("username", "loginValid", "token", "error").
		Send(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create and complete login instance: %w", err)
	}

	var resultVars domain.ProcessVariables
	if err := json.Unmarshal([]byte(result.GetVariables()), &resultVars); err != nil {
		return nil, fmt.Errorf("failed to parse login result variables: %w", err)
	}

	return &resultVars, nil
}

func StartPurchaseProcessWithResult(client zbc.Client, userID, address string, amount int) (*domain.PurchaseProcessVariables, error) {

	variables := domain.PurchaseProcessVariables{
		IsValid: true,
		UserID:  userID,
		Amount:  amount,
		Address: address,
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	command, err := client.NewCreateInstanceCommand().
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
