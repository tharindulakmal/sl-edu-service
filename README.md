sl-edu-service

A production-ready backend API service built with Go
 and Gin
.
This project follows standard Go project structure, Dockerized for portability, and ready for CI/CD + Kubernetes deployments.

Features

REST API with Gin

Config via .env (12-Factor style)

Dockerized (multi-stage build for small image size)

Health check endpoint (/health)

Ready for cloud-native deployment

Clone the Repo
git clone https://github.com/<your-username>/sl-edu-service.git
cd sl-edu-service

Install Dependencies

Ensure you have Go 1.24+ installed:

go mod tidy

Run Locally
go run cmd/server/main.go


API available at: http://localhost:8080/health

### Test the admin catalog endpoint

Once the server is running and connected to your MySQL instance, you can fetch the full catalog payload with:

```bash
curl -X GET "http://localhost:8080/api/v1/admin/catalog" \
  -H "Accept: application/json"
```

This returns grades, subjects, grade-subject mappings, lessons, topics, and subtopics in a single response for frontend filtering.

Run with Docker
Build Image
docker build -t sl-edu-service .

Run Container

Without .env:

docker run -p 8080:8080 sl-edu-service


With .env:

docker run -p 8080:8080 --env-file .env sl-edu-service

Environment Variables

Create a .env file in the project root:

PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=secret
DB_NAME=goapi

Run migration 
migrate -path db/migrations -database "mysql://admin:Tharindu#1995#11#16@tcp(127.0.0.1:3307)/sledu" up

Pushing to Docker Hub
# Tag
docker tag sl-edu-service <your-dockerhub-username>/sl-edu-service:latest

# Push
docker push <your-dockerhub-username>/sl-edu-service:latest

Testing
go test ./... -v