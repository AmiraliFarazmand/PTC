package zeebe

import (
	"context"
	"log"

	"github.com/AmiraliFarazmand/PTC_Task/internal/ports"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/entities"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/worker"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

func CheckLoginRequestWorker(client zbc.Client, userService ports.UserService) worker.JobWorker {
	jobWorker := client.NewJobWorker().
		JobType("check-login-request").
		Handler(func(jobClient worker.JobClient, job entities.Job) {
			vars, _ := job.GetVariablesAsMap()
			username := vars["username"].(string)
			password := vars["password"].(string)
		
			// Check credentials (implement your own logic)
            user, err := userService.Login(username, password)
			isValid := true
			if err != nil {
				isValid = false
			}
			log.Printf("###login worker: %+v %+v %v", user, err, isValid)

			varJob, err := jobClient.NewCompleteJobCommand().
				JobKey(job.GetKey()).
				VariablesFromMap(map[string]interface{}{"loginValid": isValid})
			if err != nil {
				log.Printf("###Failed to compelte job on login worker: %v", err.Error())
			}
			_, err = varJob.Send(context.Background())
			if err != nil {
				log.Printf("###Failed to complete check-login-request job: %v", err)
			}
		}).
		Open()
	return jobWorker
}

func CreateLoginTokenWorker(client zbc.Client) worker.JobWorker {
	jobWroker := client.NewJobWorker().
		JobType("create-login-token").
		Handler(func(jobClient worker.JobClient, job entities.Job) {
			vars, _ := job.GetVariablesAsMap()
			username := vars["username"].(string)

			// Generate token (implement your own logic)
			token := GenerateTokenForUser(username)

			varJob, err := jobClient.NewCompleteJobCommand().
				JobKey(job.GetKey()).
				VariablesFromMap(map[string]interface{}{"token": token})
			if err != nil {
				log.Printf("###Failed to compelte job on login worker: %v", err.Error())
			}
			_, err = varJob.Send(context.Background())
			if err != nil {
				log.Printf("###Failed to complete check-login-request job: %v", err)
			}
		}).
		Open()
	return jobWroker
}
//  TODO: remove this part
// Dummy token generator
func GenerateTokenForUser(username string) string {
	return "token-for-" + username
}
