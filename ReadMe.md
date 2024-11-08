# GridStream API

api to allow frontend to communicate with database


# Directory Structure


```
my-go-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── handlers/
│   │   ├── middlewares/
│   │   ├── repositories/
│   │   └── usecases/
│   └── config/
│       └── config.go
├── pkg/
│   ├── database/
│   ├── errors/
│   ├── logging/
│   └── web/
│       ├── request.go
│       └── response.go
└── tests/
    ├── integration/
    └── unit/
```