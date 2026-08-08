package grpcserver

import (
	"context"
	"log"
	"net"
	"testing"

	pbapi "github.com/oktavarium/go-gauger/api"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
)

const buffSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(buffSize)
	s := grpc.NewServer()
	storage, _ := storage.NewInMemoryStorage(context.Background(), "/tmp/test.db", false, 1000)
	pbapi.RegisterGaugerServer(s, &GrpcServer{storage: storage})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGrpcServer(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pbapi.NewGaugerClient(conn)

	updateTests := []struct {
		name string
		req  *pbapi.UpdateRequest
	}{
		{
			name: "update counter",
			req: &pbapi.UpdateRequest{
				Metric: &pbapi.Metric{
					Id:    "test1",
					Type:  "counter",
					Value: 28,
				},
			},
		},
		{
			name: "update gauger",
			req: &pbapi.UpdateRequest{
				Metric: &pbapi.Metric{
					Id:    "test2",
					Type:  "gauge",
					Value: 28.28,
				},
			},
		},
	}

	for _, test := range updateTests {
		_, err := client.Update(ctx, test.req)
		require.NoError(t, err, test.name)
	}

	getTests := []struct {
		name     string
		req      *pbapi.GetRequest
		wantResp *pbapi.GetResponse
	}{
		{
			name: "get counter",
			req: &pbapi.GetRequest{
				Name: "test1",
				Type: "counter",
			},
			wantResp: &pbapi.GetResponse{
				Value: 28,
			},
		},
		{
			name: "get gauger",
			req: &pbapi.GetRequest{
				Name: "test2",
				Type: "gauge",
			},
			wantResp: &pbapi.GetResponse{
				Value: 28.28,
			},
		},
	}

	for _, test := range getTests {
		resp, err := client.Get(ctx, test.req)
		require.NoError(t, err, test.name)
		require.Equal(t, test.wantResp.GetValue(), resp.GetValue(), test.name)
	}

	updatesTests := []struct {
		name string
		req  *pbapi.UpdatesRequest
	}{
		{
			name: "updates test",
			req: &pbapi.UpdatesRequest{
				Metrics: []*pbapi.Metric{
					{
						Id:    "test3",
						Type:  "gauge",
						Value: 29.28,
					},
					{
						Id:    "test4",
						Type:  "gauge",
						Value: 30.28,
					},
				},
			},
		},
	}

	for _, test := range updatesTests {
		_, err := client.Updates(ctx, test.req)
		require.NoError(t, err, test.name)
	}

	getAllTests := []struct {
		name     string
		wantResp *pbapi.GetAllResponse
	}{
		{
			name: "updates test",
			wantResp: &pbapi.GetAllResponse{
				Metrics: []*pbapi.Metric{
					{
						Id:    "test1",
						Type:  "counter",
						Value: 28,
					},
					{
						Id:    "test2",
						Type:  "gauge",
						Value: 28.28,
					},
					{
						Id:    "test3",
						Type:  "gauge",
						Value: 29.28,
					},
					{
						Id:    "test4",
						Type:  "gauge",
						Value: 30.28,
					},
				},
			},
		},
	}

	for _, test := range getAllTests {
		resp, err := client.GetAll(ctx, &emptypb.Empty{})
		require.NoError(t, err, test.name)
		for _, mt := range test.wantResp.GetMetrics() {
			for _, mr := range resp.GetMetrics() {
				if mt.GetId() == mr.GetId() {
					require.Equal(t, &mt, &mr, test.name)
					continue
				}
			}
		}
	}
}
