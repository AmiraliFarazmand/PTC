package zeebe

import (
	"context"
	"log"

	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/entities"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/worker"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

func ValidateCredentialsWorker(client zbc.Client) {
	log.Println("####Started validate-credentials worker")
	jobWorker := client.NewJobWorker().
		JobType("validate-credentials").
		Handler(func(jobClient worker.JobClient, job entities.Job) {
			vars, _ := job.GetVariablesAsMap()
			log.Println("####On handler function", vars)
			username := vars["username"].(string)
			password := vars["password"].(string)
			// Validate credentials (example logic)
			isValid := len(username) > 3 && len(password) > 6

			// Complete the job with the result
			varJob, err := jobClient.NewCompleteJobCommand().
				JobKey(job.GetKey()).
				VariablesFromMap(map[string]interface{}{"isValid": isValid})
			if err != nil {
				log.Printf("###Failed to compelte job: %v", err.Error())
			}
			varJob.Send(context.Background())
			if err != nil {
				log.Printf("###Failed to complete job: %v", err)
			}
		}).
		Open()
	// defer jobWorker.Close()
	log.Println("####Ended validate-credentials worker")
	defer jobWorker.Close()
}
