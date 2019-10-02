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
	}
	return &dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{},
	}, nil
}

func Test_GetAmbientTemperature(t *testing.T) {
	t.Run("ShouldGetAmbientTemperatureWhenThereIsAtLeastOneTempValue", func(t *testing.T) {
		stateRepository := &StateRepository{
			DynamoDB: &MockDynamoDB{},
		}
		ambientTempeture := stateRepository.GetAmbientTemperature("selfhydro")
		assert.Equal(t, float64(12), ambientTempeture.Temperature)
		assert.Equal(t, "selfhydro", ambientTempeture.DeviceID)
	})

	t.Run("ShouldReturnNilWhenThereIsNoAmbientTemperatureForDevice", func(t *testing.T) {
		stateRepository := &StateRepository{
			DynamoDB: &MockDynamoDB{},
		}
		ambientTempeture := stateRepository.GetAmbientTemperature("nothing")
		assert.Equal(t, float64(0), ambientTempeture.Temperature)
	})
}
