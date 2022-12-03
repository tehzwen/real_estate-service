# Real Estate Service

## Description
Started as a passion project of scraping real estate listings while looking to buy a home when the market was insane. Locally another system continues to fetch and inserts data on a regular cadence however I had no easy way of sharing this data with friends/others. 

This project serves as both an example of gRPC server & client implementation but also a means by which I can allow others to consume this data if I so desire.

## Contents
- `/cmd` contains binaries (clients, server, etc).
- `/internal` contains server internal (database, models, secret management).
- `/proto` contains the base proto definition aswell as some generated code used by the go server.
