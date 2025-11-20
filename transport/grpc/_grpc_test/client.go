package test

import (
	"context"

	"google.golang.org/grpc"

	"git-codecommit.ap-southeast-1.amazonaws.com/v1/repos/be-base/endpoint"
	grpctransport "git-codecommit.ap-southeast-1.amazonaws.com/v1/repos/be-base/transport/grpc"
	"git-codecommit.ap-southeast-1.amazonaws.com/v1/repos/be-base/transport/grpc/_grpc_test/pb"
)

type clientBinding struct {
	test endpoint.Endpoint
}

func (c *clientBinding) Test(ctx context.Context, a string, b int64) (context.Context, string, error) {
	response, err := c.test(ctx, TestRequest{A: a, B: b})
	if err != nil {
		return nil, "", err
	}
	r := response.(*TestResponse)
	return r.Ctx, r.V, nil
}

func NewClient(cc *grpc.ClientConn) Service {
	return &clientBinding{
		test: grpctransport.NewClient(
			cc,
			"pb.Test",
			"Test",
			encodeRequest,
			decodeResponse,
			&pb.TestResponse{},
			grpctransport.ClientBefore(
				injectCorrelationID,
			),
			grpctransport.ClientBefore(
				displayClientRequestHeaders,
			),
			grpctransport.ClientAfter(
				displayClientResponseHeaders,
				displayClientResponseTrailers,
			),
			grpctransport.ClientAfter(
				extractConsumedCorrelationID,
			),
		).Endpoint(),
	}
}
