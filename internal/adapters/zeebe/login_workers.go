package zeebe

import (
	"context"
	"encoding/json"
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
			var vars ProcessVariables
			if err := json.Unmarshal([]byte(job.GetVariables()), &vars); err != nil {
				log.Printf("Failed to parse variables: %v", err)
				return
			}

			// Check credentials
			user, err := userService.Login(vars.Username, vars.Password)
			if err != nil {
				vars.LoginValid = false
				vars.Error = err.Error()
			} else {
				vars.LoginValid = true
				vars.Username = user.Username // In case username casing was normalized
			}

			varsJSON, err := json.Marshal(vars)
			if err != nil {
				log.Printf("Failed to marshal variables: %v", err)
				return
			}
			log.Printf("###LOGIN VARIABLES: \n%+v\n%+v", varsJSON,vars)

			tempCommand, err := jobClient.NewCompleteJobCommand().
				JobKey(job.GetKey()).
				VariablesFromString(string(varsJSON))
				if err != nil {
					log.Printf("Failed to create command: %v", err)
				}
				tempCommand.Send(context.Background())

			if err != nil {
				log.Printf("Failed to complete job: %v", err)
			}
		}).
		Open()
	return jobWorker
}

func CreateLoginTokenWorker(client zbc.Client) worker.JobWorker {
	return client.NewJobWorker().
		JobType("create-login-token").
		Handler(func(jobClient worker.JobClient, job entities.Job) {
			var vars ProcessVariables
			if err := json.Unmarshal([]byte(job.GetVariables()), &vars); err != nil {
				log.Printf("Failed to parse variables: %v", err)
				return
			}

			// Generate token
			vars.Token = GenerateTokenForUser(vars.Username)

			varsJSON, err := json.Marshal(vars)
			if err != nil {
				log.Printf("Failed to marshal variables: %v", err)
				return
			}

			tempCommand, err := jobClient.NewCompleteJobCommand().
				JobKey(job.GetKey()).
				VariablesFromString(string(varsJSON))
				if err != nil {
					log.Printf("Failed to create command: %v", err)
				}

				tempCommand.Send(context.Background())

			if err != nil {
				log.Printf("Failed to complete job: %v", err)
			}
		}).
		Open()
}

// TODO: remove this part and implement proper token generation
func GenerateTokenForUser(username string) string {
	return "token-for-" + username
}
