package zeebe

import (
	"context"
	"encoding/json"
	"log"

	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
	"github.com/AmiraliFarazmand/PTC_Task/internal/ports"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/entities"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/worker"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

func CreatePurchaseWorker(client zbc.Client, purchaseService ports.PurchaseService) worker.JobWorker {
	return client.NewJobWorker().
		JobType("create-purchase-task").
		Handler(func(jobClient worker.JobClient, job entities.Job) {
			createPurchaseHandler(jobClient, job, purchaseService)
		}).
		Open()
}

func createPurchaseHandler(jobClient worker.JobClient, job entities.Job, purchaseService ports.PurchaseService) {
	var vars domain.PurchaseProcessVariables
	if err := json.Unmarshal([]byte(job.GetVariables()), &vars); err != nil {
		log.Printf("Failed to parse variables: %v", err)
		return
	}
	log.Println("too worker mire")
	// Create purchase using service
	purchaseID, err := purchaseService.CreatePurchase(vars.UserID, vars.Amount, vars.Address)
	if err != nil {
		vars.IsValid = false
		vars.Error = err.Error()
	} else {
		vars.IsValid = true
		vars.PurchaseID = purchaseID
	}

	varsJSON, err := json.Marshal(vars)
	if err != nil {
		log.Printf("Failed to marshal variables: %v", err)
		return
	}

	command, err := jobClient.NewCompleteJobCommand().
		JobKey(job.GetKey()).
		VariablesFromString(string(varsJSON))

	if err != nil {
		log.Printf("Failed to create complete job command: %v", err)
		return
	}

	if _, err := command.Send(context.Background()); err != nil {
		log.Printf("Failed to complete job: %v", err)
	}
}
