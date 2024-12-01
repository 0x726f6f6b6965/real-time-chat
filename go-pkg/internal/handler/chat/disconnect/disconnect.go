package disconnect

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/0x726f6f6b6965/real-time-chat/go-pkg/iaws"
	"github.com/0x726f6f6b6965/real-time-chat/go-pkg/internal/handler/chat"
	"github.com/0x726f6f6b6965/real-time-chat/go-pkg/internal/handler/common"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func DisconnectHandler(dbClient iaws.IdynamoDB) common.Handler[events.APIGatewayWebsocketProxyRequest] {
	return func(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
		// Remove the connection ID from DynamoDB
		roomId := req.Headers[common.ROOM]

		_, err := dbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(os.Getenv(chat.TABLE_NAME)),
			Key: map[string]types.AttributeValue{
				common.PK: &types.AttributeValueMemberS{Value: fmt.Sprintf(chat.RoomIdPrefix, roomId)},
				common.SK: &types.AttributeValueMemberS{Value: fmt.Sprintf(chat.ConnectionPrefix, req.RequestContext.ConnectionID)},
			},
			ConditionExpression: aws.String(common.PkExists),
		})

		if err != nil {
			log.Printf("Error deleting connection ID: %s", err)
			return common.ApiResponse(common.ErrorInternal, nil), nil
		}

		log.Printf("Client disconnected: %s", req.RequestContext.ConnectionID)
		return common.ApiResponse(common.Success, nil), nil
	}
}
