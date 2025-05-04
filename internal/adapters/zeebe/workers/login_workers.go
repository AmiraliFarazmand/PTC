package workers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
	"github.com/AmiraliFarazmand/PTC_Task/internal/ports"
	"github.com/AmiraliFarazmand/PTC_Task/internal/utils"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/entities"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/worker"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
	"github.com/golang-jwt/jwt/v5"
)

func CheckLoginRequestWorker(client zbc.Client, userService ports.UserService) worker.JobWorker {
	jobWorker := client.NewJobWorker().
		JobType("check-login-request").
		Handler(func(jobClient worker.JobClient, job entities.Job) {
			checkLoginHandler(jobClient, job, userService)
		}).
		Open()
	return jobWorker
}

func CreateLoginTokenWorker(client zbc.Client) worker.JobWorker {
	return client.NewJobWorker().
		JobType("create-login-token").
		Handler(func(jobClient worker.JobClient, job entities.Job) {
			createTokenHandler(jobClient, job)
		}).
		Open()
}

func checkLoginHandler(jobClient worker.JobClient, job entities.Job, userService ports.UserService) {
	var vars domain.AuthProcessVariables
	if err := json.Unmarshal([]byte(job.GetVariables()), &vars); err != nil {
		log.Printf("Failed to parse variables: %v", err)
		return
	}

	// Check credentials
	user, err := userService.Login(vars.Username, vars.Password)
	if err != nil {
		vars.IsValid = false
		vars.Error = err.Error()
	} else {
		vars.Username = user.Username
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
}

func createTokenHandler(jobClient worker.JobClient, job entities.Job) {
	var vars domain.AuthProcessVariables
	if err := json.Unmarshal([]byte(job.GetVariables()), &vars); err != nil {
		log.Printf("Failed to parse variables: %v", err)
		return
	}

	// Generate JWT token	TODO: move it to utils
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": vars.Username,
		"exp": time.Now().Add(time.Hour * time.Duration(72)).Unix(),
	})

	secretKey, _ := utils.ReadEnv("SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		vars.Error = "Failed to generate token"
		vars.IsValid = false
	} else {
		vars.Token = tokenString
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
}
