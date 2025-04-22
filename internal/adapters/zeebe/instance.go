package zeebe

import (
	"context"
	"fmt"
	"time"

	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/pb"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

func MustStartSignUpProcessInstance(client zbc.Client, username,password string) *pb.CreateProcessInstanceResponse {
	command, err := client.NewCreateInstanceCommand().
		BPMNProcessId("SignupProcess").
		LatestVersion().
		VariablesFromMap(map[string]interface{}{"username": username, "password": password})
	if err != nil {
		panic(fmt.Errorf("###failed to create instance %+v, %+v",err,command))
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	process, err := command.Send(ctx)
	if err != nil {
		panic(err)
	}

	return process
}


func MustStartLoginProcessInstance(client zbc.Client, username, password string) *pb.CreateProcessInstanceResponse {
    command, err := client.NewCreateInstanceCommand().
        BPMNProcessId("LoginProcess").
        LatestVersion().
        VariablesFromMap(map[string]interface{}{"username": username, "password": password})
    if err != nil {
        panic(fmt.Errorf("###failed to create login instance: %+v, %+v", err, command))
    }

    ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancelFn()

    process, err := command.Send(ctx)
    if err != nil {
        panic(err)
    }

    // log.Printf("###Started login process instance [%d] for username: %s", process.GetProcessInstanceKey(), username)
    return process
}