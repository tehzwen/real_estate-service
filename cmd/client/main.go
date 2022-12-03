package main

import (
	"context"
	"log"
	"time"

	pb "github.com/tehzwen/real_estate-service/proto/golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic(err)
	}

	client := pb.NewRealEstateClient(conn)
	var token string
	var listings []*pb.Listing
	var since time.Time = time.Date(2022, 05, 13, 0, 0, 0, 0, time.UTC)

	for {
		results, err := client.GetListings(context.Background(), &pb.GetListingsRequest{
			Filter: &pb.GetListingsFilter{
				TimeSpan: &pb.TimeSpan{
					Since: timestamppb.New(since),
				},
			},
			NextToken: token,
			Limit:     1000,
		})
		if err != nil {
			log.Panic(err)
		}
		listings = append(listings, results.Listings...)

		log.Println(len(results.Listings), token)
		if results.NextToken == "" {
			break
		}
		token = results.NextToken
	}

	log.Printf("Retrieved %d listings.", len(listings))
}
