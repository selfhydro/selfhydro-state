package state

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type PubSubMessage struct {
	Data        []byte            `json:"data"`
	Attributes  map[string]string `json:"attributes"`
	PublishTime string            `json:"publishTime"`
}

type SeflhydroState struct {
	AmbientTemperature          float64 `json:"ambientTemperature"`
	AmbientHumidity             float64 `json:"ambientHumidity"`
	WaterTemperature            float64 `json:"waterTemperature"`
	WaterElectricalConductivity float64 `json:"waterElectricalConductivity"`
	Time                        string  `json:"time"`
	deviceId                    string
}

type StateItem struct {
	AmbientTemperture float64
	Date              string
	SystemID          string
}

func TransferStateToDynamoDB(ctx context.Context, m PubSubMessage) error {
	session := createSession()
	selfhydroState := deseraliseState(m.Data)
	stateItem := createStateItem(selfhydroState)
	insertStateItem(session, stateItem)
	return nil
}

func createSession() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess, aws.NewConfig().WithRegion("ap-southeast-2"))
	return svc
}

func deseraliseState(data []byte) SeflhydroState {
	var state = SeflhydroState{}
	err := json.Unmarshal(data, &state)
	if err != nil {
		log.Fatalf("can't decode state from message: %v", err.Error())
	}
	return state
}

func createStateItem(state SeflhydroState) map[string]*dynamodb.AttributeValue {
	itemState := StateItem{
		SystemID:          state.deviceId,
		AmbientTemperture: state.AmbientTemperature,
		Date:              state.Time,
	}
	if itemState.SystemID == "" {
		itemState.SystemID = "selfhydro-default"
	}
	log.Printf("inserting state, with temperture %f for device %s", itemState.AmbientTemperture, itemState.SystemID)
	ao, err := dynamodbattribute.MarshalMap(itemState)
	if err != nil {
		fmt.Println("Got error marshalling new state:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return ao
}

func insertStateItem(dynamodbSession *dynamodb.DynamoDB, stateItem map[string]*dynamodb.AttributeValue) {
	tableName := getTableName(time.Now())
	input := &dynamodb.PutItemInput{
		Item:      stateItem,
		TableName: aws.String(tableName),
	}

	_, err := dynamodbSession.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func getTableName(time time.Time) string {
	baseName := "selfhydro-state-"
	return baseName + time.Format("2006-01-02")
}
