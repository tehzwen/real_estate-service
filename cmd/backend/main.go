package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tehzwen/real_estate-service/internal/db"
	"github.com/tehzwen/real_estate-service/internal/secrets"
	pb "github.com/tehzwen/real_estate-service/proto/golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	reg         = prometheus.NewRegistry()
	grpcMetrics = grpc_prometheus.NewServerMetrics()
	grpcPort    = flag.Int("grpc_port", 50051, "the port you wish to run grpc server on")
	metricsPort = flag.Int("metrics_port", 9092, "the port you wish to run prom metrics on")
)

func init() {
	reg.MustRegister(grpcMetrics)
}

type realEstateServer struct {
	pb.UnimplementedRealEstateServer
	DBWorker db.Worker
}

func (s *realEstateServer) GetListings(ctx context.Context, r *pb.GetListingsRequest) (*pb.GetListingsResponse, error) {
	f := &db.GetListingsFilter{}
	f.FromProto(r.Filter)
	f.PageToken = r.NextToken

	if r.Limit > 0 && r.Limit < 10000 {
		f.Limit = int(r.Limit)
	}

	listings, err := s.DBWorker.GetListings(ctx, f)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoListings []*pb.Listing
	for _, v := range listings {
		protoListings = append(protoListings, v.ToProto())
	}

	var token string
	if len(listings) > 0 {
		lastDate := listings[len(listings)-1].AddedDate
		lastId := listings[len(listings)-1].ID
		token = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d,%d", lastDate.Unix(), lastId)))
	}

	return &pb.GetListingsResponse{
		Listings:  protoListings,
		NextToken: token,
	}, nil
}

func main() {
	ctx := context.Background()

	f, err := os.Open("local-secrets.json")
	if err != nil {
		log.Fatalf("failed to open secrets file: %v", err)
	}

	s, err := secrets.FromJson(f)
	if err != nil {
		log.Fatalf("failed to load secrets from json: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption = []grpc.ServerOption{
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	}
	grpcServer := grpc.NewServer(opts...)
	d, err := db.NewPostgresWorker(ctx, db.Credentials{
		Password: s.GetSecret("DB_PASS"),
		User:     s.GetSecret("DB_USER"),
		Host:     s.GetSecret("DB_HOST"),
		Port:     s.GetSecret("DB_PORT"),
		Database: s.GetSecret("DB_DATABASE"),
	})
	if err != nil {
		log.Fatal(err)
	}

	service := &realEstateServer{
		DBWorker: d,
	}
	pb.RegisterRealEstateServer(grpcServer, service)
	grpcMetrics.EnableHandlingTimeHistogram()
	grpcMetrics.InitializeMetrics(grpcServer)

	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", *metricsPort)}
	log.Printf("Starting http server for metrics on port %d", *metricsPort)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("unable to start metrics server")
		}
	}()

	log.Printf("Starting real estate server on port %d", *grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
