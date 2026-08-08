package grpcserver

import (
	"context"

	pbapi "github.com/oktavarium/go-gauger/api"
	"github.com/oktavarium/go-gauger/internal/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GrpcServer) Get(ctx context.Context, in *pbapi.GetRequest) (*pbapi.GetResponse, error) {
	resp := &pbapi.GetResponse{}
	switch in.GetType() {
	case shared.GaugeType:
		value, ok := s.storage.GetGauger(ctx, in.GetName())
		if !ok {
			return resp, status.Errorf(codes.NotFound, "value not found")
		}
		resp.Value = value
	case shared.CounterType:
		value, ok := s.storage.GetCounter(ctx, in.GetName())
		if !ok {
			return resp, status.Errorf(codes.NotFound, "value not found")
		}
		resp.Value = float64(value)
	}

	return resp, nil
}
