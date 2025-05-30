# Go HTTP Server Integration API Testing

## Quick Start

### Run Tests
```bash
go test ./httpserver -v
```

### Start Server
```bash
go run main.go
# Server running on http://localhost:3333
```

## API Endpoints

### User CRUD Operations

#### Create User (POST)
```bash
curl -X POST http://localhost:3333/user \
  -H "Content-Type: application/json" \
  -d '{"name":"New User"}'
```

#### Get User (GET)
```bash
curl "http://localhost:3333/user?id=1"
```

#### Update User (PATCH)
```bash
curl -X PATCH http://localhost:3333/user \
  -H "Content-Type: application/json" \
  -d '{"id":1,"name":"Updated Name"}'
```

#### Delete User (DELETE)
```bash
curl -X DELETE "http://localhost:3333/user?id=1"
```

### Utility Endpoints

#### Root Endpoint
```bash
curl http://localhost:3333/
```

#### Hello Endpoint
```bash
curl http://localhost:3333/hello
```

## Response Codes

| Code | Description         |
|------|---------------------|
| 200  | OK                  |
| 201  | Created             |
| 204  | No Content          |
| 400  | Bad Request         |
| 404  | Not Found           |
| 405  | Method Not Allowed  |
| 409  | Conflict            |
