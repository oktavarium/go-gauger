package grpcserver

import (
	"context"
	"fmt"

	pbapi "github.com/oktavarium/go-gauger/api"
	"github.com/oktavarium/go-gauger/internal/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GrpcServer) Update(ctx context.Context, in *pbapi.UpdateRequest) (*emptypb.Empty, error) {
	switch in.GetMetric().Type {
	case shared.GaugeType:
		err := s.storage.SaveGauge(ctx, in.GetMetric().GetId(), in.GetMetric().GetValue())
		if err != nil {
			return &emptypb.Empty{}, status.Errorf(codes.Internal, fmt.Sprintf("error of saving gauge: %s", err))
		}

	case shared.CounterType:
		_, err := s.storage.UpdateCounter(ctx, in.GetMetric().GetId(), int64(in.GetMetric().GetValue()))
		if err != nil {
			return &emptypb.Empty{}, status.Errorf(codes.Internal, fmt.Sprintf("error of saving counter: %s", err))
		}
	}

	return &emptypb.Empty{}, nil
}
