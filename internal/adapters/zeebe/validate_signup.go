package zeebe

import (
    "context"
    "log"

    "github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

func DeploySignupProcess(client zbc.Client) {
    deployResouceRespond, err := client.NewDeployResourceCommand().
        AddResourceFile("bpmn/signup.bpmn").
        Send(context.Background())
    if err != nil {
        log.Fatalf("####Failed to deploy process: %v", err)
    }
    log.Printf("###Signup process deployed successfully %+v\n",deployResouceRespond)
}