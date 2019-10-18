package main

import (
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
	return baseName + time.Format("2006-01")
}

func (stateRepository StateRepository) GetWaterTemperature(systemID string) WaterTemperature {
	log.Printf("getting water temperature for %s device", systemID)
	tableName := getTableName(time.Now().UTC())
	query := stateRepository.createWaterTemperatureQuery(tableName, systemID)
	queryOutput, err := stateRepository.DynamoDB.Query(query)
	if err != nil {
		log.Printf("could not query dynamodb for water temperature: %s", err.Error())
	}
	if len(queryOutput.Items) == 0 {
		return WaterTemperature{}
	}
	temperature, err := strconv.ParseFloat(*queryOutput.Items[0]["WaterTemperature"].N, 64)
	if err != nil {
		log.Printf("could not parse water temperature to float: %s", err.Error())
	}
	timestamp, err := time.Parse("20060102150405", *queryOutput.Items[0]["Date"].S)
	if err != nil {
		log.Printf("could not parse water temperture timestamp to time struct: %s", err.Error())
	}
	ambientTemperature := WaterTemperature{
		Temperature: temperature,
		DeviceID:    systemID,
		Timestamp:   timestamp,
	}
	return ambientTemperature
}

func (stateRepository StateRepository) createWaterTemperatureQuery(tableName string, systemID string) *dynamodb.QueryInput {
	query := &dynamodb.QueryInput{}
	query.SetTableName(tableName)
	query.SetExpressionAttributeNames(map[string]*string{
		"#system_id": aws.String("SystemID"),
		"#date":      aws.String("Date"),
	})
	query.SetProjectionExpression("WaterTemperature, #date")
	query.SetExpressionAttributeValues(map[string]*dynamodb.AttributeValue{
		":s1": {
			S: aws.String(systemID),
		},
		":d1": {
			S: aws.String(time.Now().Add(time.Duration(-1) * time.Hour).Format("200601021504")),
		},
	})
	query.SetKeyConditionExpression("#system_id = :s1 AND #date > :d1")
	return query
}
