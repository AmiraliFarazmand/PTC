package zeebe

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/core/app"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/entities"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/worker"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

func ValidateCredentialsWorker(client zbc.Client, userRepo *db.MongoUserRepository) worker.JobWorker {
	jobWorker := client.NewJobWorker().
		JobType("validate-credentials").
		Handler(func(jobClient worker.JobClient, job entities.Job) {
			// Parse incoming variables
			var vars ProcessVariables
			if err := json.Unmarshal([]byte(job.GetVariables()), &vars); err != nil {
				log.Printf("Failed to parse variables: %v", err)
				return
			}

			// Validate credentials
			isUsernameUnique, err := userRepo.IsUsernameUnique(vars.Username)
			vars.IsValid = len(vars.Username) > 3 && len(vars.Password) > 6
			if err != nil || !isUsernameUnique {
				vars.IsValid = false
				vars.Error = "Username already taken or validation failed"
			}

			// Complete the job with updated variables
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
		Concurrency(1).
		MaxJobsActive(10).
		RequestTimeout(1 * time.Second).
		PollInterval(1 * time.Second).
		Name("validate-credential").
		Open()
	return jobWorker
}

func CreateUserWorker(client zbc.Client, userService app.UserServiceImpl) worker.JobWorker {
	return client.NewJobWorker().
		JobType("create-user").
		Handler(func(jobClient worker.JobClient, job entities.Job) {
			var vars ProcessVariables
			if err := json.Unmarshal([]byte(job.GetVariables()), &vars); err != nil {
				log.Printf("Failed to parse variables: %v", err)
				return
			}

			err := userService.Signup(vars.Username, vars.Password)
			if err != nil {
				vars.Error = err.Error()
				vars.IsValid = false
				log.Printf("Failed to create user: %v", err)
			} else {
				vars.IsValid = true
			}

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
