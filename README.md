# reward-management-system

The Reward Management System is a loyalty points management API for an e-commerce platform. The system enables users to earn points for each purchase, which can then be redeemed for discounts on future purchases. It is designed to handle high transaction volumes, ensure data consistency, and maintain a detailed history of points earned, redeemed, and expired.

# Required
Before setting up the system, make sure you have the following installed:

  1.  Docker Compose - For running PostgreSQL in a container.
  2.  Golang - For running the main service and scheduler.

# Set Up process

Follow the steps below to set up the project locally:
1. Set Up Configuration File
  Create and configure your config.yaml file in the main directory. This file should include necessary information like database credentials, server configurations, etc.

2. Switch to the Main Directory
  Navigate to the main directory where the docker-compose.yml and main.go files are.

3. Start PostgreSQL with Docker Compose
  Run the following command to start your PostgreSQL container and initialize the database tables:
  `docker-compose up -d`
4. Start Main service commander -
  `go mod tidy`
  `go run main.go`
5. Start reward-expiration-schedular 
  `cd reward-expiration-schedular`
  `go run main.go` 
