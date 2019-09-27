package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type StateRepository struct {
	DynamoDB dynamodbiface.DynamoDBAPI
}

func NewStateRepository() *StateRepository {
	stateRepository := &StateRepository{}
	stateRepository.createSession()
	return stateRepository
}

func (stateRepostiroy *StateRepository) createSession() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess, aws.NewConfig().WithRegion("ap-southeast-2"))
	stateRepostiroy.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
}

func getTableName(time time.Time) string {
	baseName := "selfhydro-state-"
	return baseName + time.Format("2006-01-02")
}

func (stateRepository StateRepository) GetAmbientTemperature(systemID string) AmbientTemperature {
	tableName := getTableName(time.Now().UTC())
	condition := fmt.Sprintf("SystemID = %s", systemID)
	query := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String(condition),
		ProjectionExpression:   aws.String("AmbientTemperature, Date"),
	}
	queryOutput, _ := stateRepository.DynamoDB.Query(query)
	temperature, _ := strconv.ParseFloat(*queryOutput.Items[0]["AmbientTemperature"].N, 64)
	ambientTemperature := AmbientTemperature{
		Temperature: temperature,
		DeviceID:    systemID,
	}
	return ambientTemperature
}
