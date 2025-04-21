package zeebe

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/pb"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

func MustStartProcessInstance(client zbc.Client, username,password string) *pb.CreateProcessInstanceResponse {
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

	log.Printf("###started process instance [%d] with {\"%s\": \"%s\"}", process.GetProcessInstanceKey(), "username", username)
	return process
}