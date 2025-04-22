package zeebe

import (
	"context"
	"log"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/entities"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/worker"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

func ValidateCredentialsWorker(client zbc.Client) worker.JobWorker{
	jobWorker := client.NewJobWorker().
		JobType("validate-credentials").
		Handler(func(jobClient worker.JobClient, job entities.Job) {
			vars, _ := job.GetVariablesAsMap()
			username := vars["username"].(string)
			password := vars["password"].(string)
			// log.Println("####On handler function", vars, username, password)
			// Validate credentials (example logic)
			isValid := len(username) > 3 && len(password) > 6

			// Complete the job with the result
			varJob, err := jobClient.NewCompleteJobCommand().
				JobKey(job.GetKey()).
				VariablesFromMap(map[string]interface{}{"isValid": isValid})
			if err != nil {
				log.Printf("###Failed to compelte job: %v", err.Error())
			}
			_, err =varJob.Send(context.Background())
			if err != nil {
				log.Printf("###Failed to complete job: %v", err)
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


func CreateUserWorker(client zbc.Client, userRepo *db.MongoUserRepository) worker.JobWorker {
    jobWorker := client.NewJobWorker().
        JobType("create-user").
        Handler(func(jobClient worker.JobClient, job entities.Job) {
            vars,_ := job.GetVariablesAsMap()
            username := vars["username"].(string)
            password := vars["password"].(string)
 
            err := userRepo.Create(domain.User{Username: username, Password: password})
            if err != nil {
                log.Printf("###failed to create user: %v", err)
                return 
            }

            // Complete the job
            _, err = jobClient.NewCompleteJobCommand().
                JobKey(job.GetKey()).
                Send(context.Background())
            if err != nil {
                log.Printf("###Failed to complete job: %v", err)
            }
        }).
        Open()
    return jobWorker 
}