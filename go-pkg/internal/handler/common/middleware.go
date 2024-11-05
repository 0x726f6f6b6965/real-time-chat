package common

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func TokenMiddleware[input Request](next Handler[input]) Handler[input] {
	return func(ctx context.Context, req input) (events.APIGatewayProxyResponse, error) {
		headers := getHeaders(req)
		data := strings.Split(headers[AUTH], " ")
		if len(data) < 2 {
			log.Printf("Error getting token, %s", data)
			return ApiResponse(ErrorValidation, nil), nil
		}
		token := data[1]
		email, roomId, err := ExtractToken(token)
		if err != nil {
			log.Printf("Error getting token, %s, secret: %s", token, os.Getenv(SECRET))
			return ApiResponse(ErrorValidation, nil), nil
		}
		setUser(req, email)
		setRoomId(req, roomId)
		return next(ctx, req)
	}
}

func getHeaders(data any) map[string]string {
	switch data := data.(type) {
	case events.APIGatewayWebsocketProxyRequest:
		return data.Headers
	case events.APIGatewayProxyRequest:
		return data.Headers
	default:
		return map[string]string{}
	}
}

func setRoomId(data any, roomId int) {
	switch data := data.(type) {
	case events.APIGatewayWebsocketProxyRequest:
		data.Headers[ROOM] = fmt.Sprintf("%d", roomId)
	case events.APIGatewayProxyRequest:
		data.Headers[ROOM] = fmt.Sprintf("%d", roomId)
	}
}

func setUser(data any, user string) {
	switch data := data.(type) {
	case events.APIGatewayWebsocketProxyRequest:
		data.Headers[USER] = user
	case events.APIGatewayProxyRequest:
		data.Headers[USER] = user
	}
}
