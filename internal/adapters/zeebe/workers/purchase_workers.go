package workers

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

// type FilaanClient struct {
// 	zbc.Client
// 	purchaseService ports.PurchaseService

func CreatePurchaseWorker(client zbc.Client, purchaseService ports.PurchaseService) worker.JobWorker {
	return client.NewJobWorker().
		JobType("create-purchase-task").
		Handler(func(jobClient worker.JobClient, job entities.Job) {
			createPurchaseHandler(jobClient, job, purchaseService)
		}).
		Open()
}

func ProcessPaymentWorker(client zbc.Client) worker.JobWorker {
	jobWorker := client.NewJobWorker().
		JobType("start-payment-process").
		Handler(func(jobClient worker.JobClient, job entities.Job) {
			startPaymentHandler(client, jobClient, job)
		}).
		Open()
	return jobWorker
}

func CancelUnpaidPurchaseWorker(client zbc.Client, purchaseService ports.PurchaseService) worker.JobWorker {
	return client.NewJobWorker().
		JobType("cancel-if-unpaid-task").
		Handler(func(jobClient worker.JobClient, job entities.Job) {
			cancelUnpaidHandler(jobClient, job, purchaseService)
		}).
		Open()
}

func createPurchaseHandler(jobClient worker.JobClient, job entities.Job, purchaseService ports.PurchaseService) {
	var vars domain.PurchaseProcessVariables
	if err := json.Unmarshal([]byte(job.GetVariables()), &vars); err != nil {
		log.Printf("Failed to parse variables: %v", err)
		return
	}
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

func startPaymentHandler(zeebeClient zbc.Client, client worker.JobClient, job entities.Job) {
	// Get variables from the job
	variables, err := job.GetVariablesAsMap()
	if err != nil {
		log.Println("Failed to get variables:", err)
		return
	}

	purchaseID, ok := variables["purchase_id"].(string)
	if !ok {
		log.Println("purchase_id not found in variables")
		return
	}

	// Publish message to trigger CheckPaymentProcess
	ctx := context.Background()
	tempCommand, err := zeebeClient.NewPublishMessageCommand().
		MessageName("start-check-payment"). // Must match message name in BPMN
		CorrelationKey(purchaseID).         // Used to correlate the message
		VariablesFromMap(variables)         // Pass along all variables
	tempCommand.Send(ctx)

	if err != nil {
		log.Println("Failed to publish message:", err)
		return
	}

	// Complete the task
	_, err = client.NewCompleteJobCommand().
		JobKey(job.GetKey()).
		Send(ctx)

	if err != nil {
		log.Println("Failed to complete start-payment-process:", err)
		return
	}
}

func cancelUnpaidHandler(jobClient worker.JobClient, job entities.Job, purchaseService ports.PurchaseService) {
	var vars domain.PurchaseProcessVariables
	if err := json.Unmarshal([]byte(job.GetVariables()), &vars); err != nil {
		log.Printf("Failed to parse variables: %v", err)
		return
	}

	err := purchaseService.CancelUnpaidPurchase(vars.PurchaseID)
	if err != nil {	//is useless for now; its BPMN can be modified 
		vars.Error = err.Error()
		vars.IsValid = false
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
