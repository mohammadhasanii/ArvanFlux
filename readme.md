# ArvanFlux 1.1.0

If you want to have a traffic distributor with a round-robin algorithm and distribute your frontend or backend requests to your own servers, you can use this project. Because the project is Dockerized, you can easily use this on the [ArvanCloud](https://arvancloud.ir) service.

- Round Robin Load Balancing: Distributes traffic evenly across a list of backend servers.
- Simple Configuration: Easy to set up with minimal configuration.
- Efficient Performance: Optimized for high performance and low latency
- Automatic server failover with configurable downtime: Temporarily removes non-responsive servers from the rotation for a set duration (1 minutes) to ensure high availability

# Getting Started

To get started with ArvanFlux, you need to have Docker installed on your machine. Follow the instructions below to build and run the load balancer.

# Prerequisites

- Docker

# Building the Docker Image

1. Clone the repository or download the source code.

2. Navigate to the project directory where the Dockerfile is located.

3. Build the Docker image using the following command:

4. Clone the repository:

```bash
   docker build -t arvanflux .
```

# Running the Load Balancer

1. Run the Docker container using the following command:

```bash
  docker run -p 8080:8080 arvanflux

```

This command maps port 8080 of the container to port 8080 on your host machine.

2. Access the load balancer by navigating to http://localhost:8080 in your web browser or using an HTTP client.

# Configuration

Update the `servers` slice in the `main.go` file with the addresses of your backend servers.

```go
servers := []string{
    "http://localhost:8081", // Replace with your backend servers or frontend
    "http://localhost:8082",
    "http://localhost:8083",
}

```

# Example Usage

Once running, ArvanFlux will distribute incoming requests to the servers specified in the `servers` slice using a Round Robin approach.

# Contributing

If you would like to contribute to this project, please fork the repository and submit a pull request.
