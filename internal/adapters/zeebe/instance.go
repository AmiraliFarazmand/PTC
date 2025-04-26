package zeebe

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/pb"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

type ProcessVariables struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	IsValid    bool   `json:"isValid"`
	LoginValid bool   `json:"loginValid"` 
	Token      string `json:"token"`
	Error      string `json:"error"`
}


func StartSignUpProcessInstanceWithResult(client zbc.Client, username, password string) (*ProcessVariables, error) {
	variables := ProcessVariables{
		Username: username,
		Password: password,
		IsValid: true,
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
	var resultVars ProcessVariables
	if err := json.Unmarshal([]byte(result.GetVariables()), &resultVars); err != nil {
		return nil, fmt.Errorf("failed to parse result variables: %w", err)
	}

	return &resultVars, nil
}

func StartLoginProcessInstanceWithResult(client zbc.Client, username, password string) (*ProcessVariables, error) {
	variables := ProcessVariables{
		Username: username,
		Password: password,
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
		FetchVariables("username", "loginValid", "token").
		Send(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create and complete login instance: %w", err)
	}

	var resultVars ProcessVariables
	if err := json.Unmarshal([]byte(result.GetVariables()), &resultVars); err != nil {
		return nil, fmt.Errorf("failed to parse login result variables: %w", err)
	}

	return &resultVars, nil
}

// Kept for backward compatibility
func MustStartLoginProcessInstance(client zbc.Client, username, password string) *pb.CreateProcessInstanceResponse {
	variables := ProcessVariables{
		Username: username,
		Password: password,
	}

	command, err := client.NewCreateInstanceCommand().
		BPMNProcessId("LoginProcess").
		LatestVersion().
		VariablesFromObject(variables)
	if err != nil {
		panic(fmt.Errorf("failed to create login instance: %+v, %+v", err, command))
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	process, err := command.Send(ctx)
	if err != nil {
		panic(err)
	}

	return process
}
