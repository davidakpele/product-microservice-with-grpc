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
> The proto define the structure of the data being sent over the wire and the service methods that can be invoked remotely. They are used to generate client and server.

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
> During your installation GORM and Database Driver will be install but you prefere  manually installation, use the below command.
```
# Install GORM and PostgreSQL driver
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

### Install & Generate the protocol buffers using protoc:
####  If You Don't have Protocol Buffers compiler install.

- Install the necessary packages for gRPC support:
    - To install the Protocol Buffers compiler (protoc) on Windows, follow these steps:
    - Step 1: Download Protocol Buffers
    - Visit the official Protocol Buffers GitHub releases page:
        - [Download Docker Desktop (macOS/Windows)](https://github.com/protocolbuffers/protobuf/releases)
    - Download the latest version of the precompiled binaries for Windows:
        - Look for a file named something like protobuf-29.3.zip
    - Step 2: Extract the ZIP File
      - Extract the contents of the downloaded ZIP file to a folder on your system, such as:
```
C:\protobuf
```
- The folder should contain:
- Step 3: Add to System PATH 
    - Add the bin directory of the extracted folder to your system's PATH environment variable:
        - Open the Start Menu and search for Environment Variables.
        - Click Edit the system environment variables.
        - In the System Properties window, click the Environment Variables button.
        - Under System Variables, locate the Path variable, select it, and click Edit.
        - Click New and add the path to the bin directory (e.g., C:\protobuf\bin).
        - Click OK to save and close all windows.
- Step 4: Verify Installation 
    - Open a new Command Prompt (cmd) or PowerShell window.
    - Run the following command to check if protoc is installed:
```
protoc --version
```
### If You already have Protocol Buffers compiler installed in your system, move to the next stage below
# Install gRPC and Protocol Buffers

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
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
    - Also added ENV instructions for each environment variable that is needs.

- Building and Running the Docker Image:
```
docker build -t your-app-name .
```
> In my case i use docker build -t product-microservice .

- Run the Docker container: After building the image, run the container:
```
docker run -p 8080:8080 product-microservice
```

## Assumptions and Constraints
- Database: Ensure the database is properly configured and the product and subscription models are correctly related.
- Deployment: This microservice can be deployed using Docker for easy management.
