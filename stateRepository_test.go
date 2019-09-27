package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"gotest.tools/assert"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func (c *MockDynamoDB) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	time := time.Now().Format("2006-01-02")
	tableName := fmt.Sprintf("selfhydro-state-%s", time)
	items := []map[string]*dynamodb.AttributeValue{}
	state1 := map[string]*dynamodb.AttributeValue{
		"SystemID": &dynamodb.AttributeValue{
			S: aws.String("selfhydro"),
		},
		"AmbientTemperature": &dynamodb.AttributeValue{
			N: aws.String("12"),
		},
	}
	items = append(items, state1)
	if *input.KeyConditionExpression == "SystemID = selfhydro" && *input.TableName == tableName {
		return &dynamodb.QueryOutput{
			Items: items,
		}, nil
	} else {
		fmt.Printf("KeyConditionExpression: %s and TableName: %s", *input.KeyConditionExpression, *input.TableName)
	}
	return nil, nil
}

func Test_ShouldGetAmbientTemperature(t *testing.T) {
	stateRepository := &StateRepository{
		DynamoDB: &MockDynamoDB{},
	}
	ambientTempeture := stateRepository.GetAmbientTemperature("selfhydro")
	assert.Equal(t, float64(12), ambientTempeture.Temperature)
	assert.Equal(t, "selfhydro", ambientTempeture.DeviceID)
}
