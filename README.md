# Microservice Implementation with gRPC, Golang, and GORM


## Table of Contents
* **Introduction**
* **Folder Structure**
* **Application Configuration**
* **DB Package**
* **Internal Package**
* **Testing**
  
## Introduction
> Developed a product microservice that exposes gRPC endpoints. The microservice manage different types of products, each potentially having specific fields and associated subscription plans. 

# Folder Structure 
```
C:.
├───config
├───db
├───internal
│   ├───domain
│   ├───repository
│   ├───service
│   └───transport
│       └───grpc
├───pkg
│   └───logger 
├───proto      
│   ├───product
│   └───subscription
└───test

```
> I named my base folder product-micoservice.

## Application Configuration 
> Config holds the application database configuration and info from env file.

## DB Package 
> Holds the application database connections, I'm using Postgresql.

## Internal Package
> Holds four important packages to this application setup
- Packages Under:
    - domain package: which represent (model classes `product and subscription`) that holds the database tables entities.
    - repository package: This package holds classes hides the details of how data is fetched or persisted in the database.
    - service package: This package holds classes responsible for implementing the business logic of the application.
    - transport package: The package holds a sub package called `grpc` and the role is to mediate between the gRPC server and the business logic layer. It receives incoming gRPC requests, calls the necessary business logic, and sends back the responses.

## Proto Package
> The proto define the structure of the data being sent over the wire and the service methods that can be invoked remotely. They are used to generate client and server code in multiple programming languages, making it possible for services to communicate with each other.

## Test Package
> This package holds classes to test our grpc endpoints.

### Clone the Repository
```
git clone <repository_url>
cd <repository_name>
```
> Install Dependencies
- Make sure you have Go installed and GORM set up with a SQL-based database (e.g., PostgreSQL, MySQL).
    - Run the following command to install the necessary dependencies:

```
go mod tidy
```
### Generate the protocol buffers using protoc:
- open your base project directory navigate to proto folder then run the commands belows:
```
protoc --go_out=../ --go-grpc_out=../ product.proto

protoc --go_out=../ --go-grpc_out=../ subscription.proto
```

### gRPC Endpoints
#### Product Service
- CreateProduct:
    - Description: Create a new product.
        - Request:
```
message CreateProductRequest {
  string name = 1;
  string description = 2;
  float price = 3;
}

```

- Response:
```
message CreateProductResponse {
  string id = 1;
}
```
- GetProduct:
```
message GetProductRequest {
  string id = 1;
}
```
- Response:

```
message GetProductResponse {
  string name = 1;
  string description = 2;
  float price = 3;
}

```
- UpdateProduct:
    - Description: Update product details
        - Request:

```
message UpdateProductRequest {
  string id = 1;
  string name = 2;
  string description = 3;
  float price = 4;
}

```

- Response:

```
message UpdateProductResponse {
  string id = 1;
}
```
- DeleteProduct:
    - Description: Delete a product by ID.
        - Request:
```
message DeleteProductRequest {
  string id = 1;
}
```

- Response:
```
message DeleteProductResponse {
  string id = 1;
}
```
#### Subscription Service
- CreateSubscriptionPlan:
    - Description: Create a subscription plan for a product.
        - Request:
```
message CreateSubscriptionPlanRequest {
  string product_id = 1;
  string plan_name = 2;
  int32 duration = 3;
  float price = 4;
}
```

-  Response

```
message CreateSubscriptionPlanResponse {
  string id = 1;
}
```

- GetSubscriptionPlan:
    - Description: Get a subscription plan by ID.
        - Request:
```
message GetSubscriptionPlanRequest {
  string id = 1;
}

```

- Response:
```
message GetSubscriptionPlanResponse {
  string product_id = 1;
  string plan_name = 2;
  int32 duration = 3;
  float price = 4;
}

```
## Dockerization Process

> The Default Golang version installed was go 1.23.1 in go.mod file rename to 1.23 or Just make sure the golang version inside go.mod matches with the Dockerfile FROM golang:1.23-alpine vision.

- Created Dockerfile to containerize our application to make it easier to deploy and manage.
    - Also added ENV instructions for each environment variable that your application needs.

- Building and Running the Docker Image:
```
docker build -t your-app-name .
```
> In my case i use docker build -t product-microservice .

- Run the Docker container: After building the image, run the container:
```
docker run -p 8080:8080 product-microservice
```
