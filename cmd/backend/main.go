package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/tehzwen/real_estate-service/internal/db"
	"github.com/tehzwen/real_estate-service/internal/secrets"
	pb "github.com/tehzwen/real_estate-service/proto/golang"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedRealEstateServer
	DBWorker db.Worker
}

func (s *server) GetListings(ctx context.Context, r *pb.GetListingsRequest) (*pb.GetListingsResponse, error) {
	f := &db.GetListingsFilter{}
	f.FromProto(r.Filter)
	f.PageToken = r.NextToken

	if r.Limit > 0 && r.Limit < 10000 {
		f.Limit = int(r.Limit)
	}

	listings, err := s.DBWorker.GetListings(ctx, f)
	if err != nil {
		return nil, err
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

	f, err := os.Open("../../local-secrets.json")
	if err != nil {
		log.Fatalf("failed to open secrets file: %v", err)
	}

	s, err := secrets.FromJson(f)
	if err != nil {
		log.Fatalf("failed to load secrets from json: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
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

	service := &server{
		DBWorker: d,
	}
	pb.RegisterRealEstateServer(grpcServer, service)

	log.Println("Starting real estate server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
