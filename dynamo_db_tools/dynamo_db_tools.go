package dynamo_db_tools

import (
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
