package dynamo_db_tools

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/Akavall/GoGamesProject/zombie_dice"
)

func PutGameStateInDynamoDB(game_state zombie_dice.GameState) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsWest2RegionID),
	}))

	// svc should probably be created once,
	// and passed to a function
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(game_state)
	if err != nil {
		return err
	}

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("GameStates"),
		Item:      av,
	})

	if err != nil {
		return err
	}

	return nil
}

func GetGameStateFromDynamoDB(game_state_uuid string) (zombie_dice.GameState, error) {
	// probably should be passing svc here
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsWest2RegionID),
	}))

	svc := dynamodb.New(sess)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"game_state_id": {
				S: aws.String(game_state_uuid),
			},
		},
		TableName: aws.String("GameStates"),
	}

	result, err := svc.GetItem(input)

	fmt.Printf("ERROR: %v", err)

	if err != nil {
		panic(fmt.Sprintf("failed to get Record from DynamoDB table: GameStates, uuid: %s, %v", err, game_state_uuid))
	}

	fmt.Printf("result: %v\n", result)

	var game_state zombie_dice.GameState

	dynamodbattribute.UnmarshalMap(result.Item, &game_state)

	return game_state, err
}

func DeleteGameStateFromDynamoDB(game_state_uuid string) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsWest2RegionID),
	}))

	svc := dynamodb.New(sess)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"game_state_id": {
				S: aws.String(game_state_uuid),
			},
		},
		TableName: aws.String("GameStates"),
	}

	_, err := svc.DeleteItem(input)

	return err

}
