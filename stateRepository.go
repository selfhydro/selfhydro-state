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
	return baseName + time.Format("2006-01-02")
}

func (stateRepository StateRepository) GetAmbientTemperature(systemID string) AmbientTemperature {
	log.Printf("getting ambient temperature for %s device", systemID)
	tableName := getTableName(time.Now().UTC())
	query := &dynamodb.QueryInput{}
	query.SetTableName(tableName)
	query.SetExpressionAttributeNames(map[string]*string{
		"#D": aws.String("Date"),
	})
	query.SetProjectionExpression("AmbientTemperature, #D")
	query.SetExpressionAttributeValues(map[string]*dynamodb.AttributeValue{
		":s1": {
			S: aws.String(systemID),
		},
		":d1": {
			S: aws.String(time.Now().Add(time.Duration(-4) * time.Hour).Format("20060102150405")),
		},
	})
	query.SetKeyConditionExpression("SystemID = :s1 AND #D > :d1")
	log.Printf("querying dynamo with %s", query.GoString())
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
