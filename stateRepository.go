package main

import (
	"fmt"
	"log"
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
	log.Printf("getting ambient temperature for %s device", systemID)
	tableName := getTableName(time.Now().UTC())
	condition := fmt.Sprintf("SystemID = %s", systemID)

	query := &dynamodb.QueryInput{
		ExpressionAttributeNames: map[string]*string{
			"#D": aws.String("Date"),
		},
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String(condition),
		ProjectionExpression:   aws.String("AmbientTemperature, #D"),
	}
	queryOutput, err := stateRepository.DynamoDB.Query(query)
	if err != nil {
		log.Printf("could not query dynamo for ambient temperature: %s", err.Error())
	}
	if len(queryOutput.Items) == 0 {
		return AmbientTemperature{}
	}
	temperature, err := strconv.ParseFloat(*queryOutput.Items[0]["AmbientTemperature"].N, 64)
	if err != nil {
		log.Printf("could not parse ambient temperature to float: %s", err.Error())
	}
	ambientTemperature := AmbientTemperature{
		Temperature: temperature,
		DeviceID:    systemID,
	}
	return ambientTemperature
}
