package zeebe

import (
    "context"
    "log"

    "github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
    "github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
    "github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
    "github.com/camunda-community-hub/zeebe-client-go/v8/pkg/worker"
    "github.com/camunda-community-hub/zeebe-client-go/v8/pkg/entities"
)
func CreateUserWorker(client zbc.Client, userRepo *db.MongoUserRepository) worker.JobWorker {
    jobWorker := client.NewJobWorker().
        JobType("create-user").
        Handler(func(jobClient worker.JobClient, job entities.Job) {
            vars,_ := job.GetVariablesAsMap()
            username := vars["username"].(string)
            password := vars["password"].(string)
            // isValid := vars["isValid"].(bool)

            // Create user in the database
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
        log.Println("###CreateUserWorker started")
    // defer jobWorker.Close()
    return jobWorker 
}