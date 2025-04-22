package zeebe

import (
	"context"
	"errors"
	"time"

	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/pb"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

func DeploySignupProcess(client zbc.Client) *pb.ProcessMetadata {
	command := client.NewDeployResourceCommand().
		AddResourceFile("bpmn/signup.bpmn")

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	resource, err := command.Send(ctx)
	if err != nil {
		panic(err)
	}

	if len(resource.GetDeployments()) < 0 {
		panic(errors.New("failed to deploy send-email model; nothing was deployed"))
	}

	deployment := resource.GetDeployments()[0]
	process := deployment.GetProcess()
	if process == nil {
		panic(errors.New("failed to deploy send-email process; the deployment was successful, but no process was returned"))
	}

	// log.Printf("###deployed BPMN process [%s] with key [%d]", process.GetBpmnProcessId(), process.GetProcessDefinitionKey())
	// log.Printf("###Signup process deployed successfully %+v\n", process)
	return process
}
