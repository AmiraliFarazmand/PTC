package zeebe

import (
    "context"
    "log"

    "github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

func DeploySignupProcess(client zbc.Client) {
    _, err := client.NewDeployResourceCommand().
        AddResourceFile("bpmn/signup.bpmn").
        Send(context.Background())
    if err != nil {
        log.Fatalf("Failed to deploy process: %v", err)
    }
    log.Println("Signup process deployed successfully")
}