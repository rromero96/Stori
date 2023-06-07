package system

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func GetHTMLInfoV1(ctx context.Context, request events.APIGatewayProxyRequest, htmlProcessTransactions HTMLProcessTransactions) (events.APIGatewayProxyResponse, error) {
	html, err := htmlProcessTransactions(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("error getting html info: %w", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(html),
	}, nil
}
