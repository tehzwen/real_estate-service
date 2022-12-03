#!/opt/homebrew/bin/python3
import grpc
import real_estate_service_pb2_grpc
import real_estate_service_pb2
import datetime

from google.protobuf.timestamp_pb2 import Timestamp

def get_listings(stub, limit):
    since_time = datetime.datetime(2022, 5, 13)
    timestamp = Timestamp()
    timestamp.FromDatetime(since_time)

    token = ""
    listings = []
    while(True):
        resp = stub.GetListings(real_estate_service_pb2.GetListingsRequest(
            next_token=token,
            limit=limit,
            filter=real_estate_service_pb2.GetListingsFilter(
                time_span = real_estate_service_pb2.TimeSpan(
                    since = timestamp
                ),
                min_price = 200000,
                max_price = 500000
            )
        ))
        print(".")

        if resp.nextToken == "":
            break

        listings.extend(resp.listings)
        token = resp.nextToken

    return listings

def main():
    channel = grpc.insecure_channel('localhost:50051')
    stub = real_estate_service_pb2_grpc.RealEstateStub(channel)
    listings = get_listings(stub, 1000)
    print(f"Retrieved {len(listings)} listings.")
    

if __name__ == "__main__":
    main()
