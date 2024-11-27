# GridStream API

api to allow frontend to communicate with database

## 🧭 Table of Contents

- [API](#GridStream-API)
  - [Table of Contents](#-table-of-contents)
  - [Team](#-team)
  - [Directory Structure](#-directory-structure)
  - [Contributing](#-contributing)
  - [Local Run](#-local-run)
    - [Prerequisites](#prerequisites)

   
## 👥 Team

| Team Member     | Role Title                | Description                                                                                                                                             |
| --------------- | ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Matthew Collett | Technical Lead/Developer  | Focus on architecture design and solving complex problems, with a focus on the micro-batching process.                                                  |
| Cooper Dickson  | Project Manager/Developer | Ensure that the scope and timeline are feasible and overview project status, focus on UI and real-time transmission.                                    |
| Eric Cuenat     | Scrum Master/Developer    | In charge of agile methods for the team such as organizing meetings, removing blockers, and team communication, focus on UI and web socket interaction. |
| Sam Keays       | Product Owner/Developer   | Manager of product backlog and updating board to reflect scope changes and requirements, focus on database operations and schema design.                |


## 🏗️ Directory Structure


```
api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── handlers/
│   │   ├── middlewares/
│   │   ├── repositories/
│   │   ├── server/
│   │   │    └── routes.go
│   │   │    └── server.go
│   │   ├── models/
│   │   └── logic/
│   └── config/
│       └── config.go
├── pkg/
└── tests/
    ├── integration/
    └── unit/
```

- `cmd/api/main.go`
    - Entry point of the application where we set up the HTTP server and wire up components. Runs HTTP server created by server.go in a goroutine.

- `internal/app` Contains the core application logic, organized into subfolders
    - `handlers/`
        - Handle the HTTP request; validate schema, apply business logic from logic package (optional), uses repository package to insert and returns response.  
    - `middlewares/`
        - Optional functions that can be applied to HTTP requests. For example, making sure a user is authenticated. (Adapter pattern).  
    - `repositories/`
        - Data access layer that interacts directly with the database. Called from handler package and returns response.  
    - `logic/`
        - Holds complex business logic that needs to be applied to data if necessary.
    - `models/`
        - Contains structs representing database schemas. Can be used to insert into database and validate schemas.
    - `server/` Contains initialization for server and registering routes.
        - `server.go`
            - Creates and returns server struct for main.go to run. Calls routes.go to register routes for requests.
        - `routes.go`
            - Initializes router, registers routes and applies necessary middleware to requests. Called by server.go
    
- `internal/config/config.go`
    - Holds Application configuration, such as database connection details, API keys, etc.  

- `pkg`
    - Holds reusuable libraries to be used accross different projects. Example, Firebase client, Big Query client, generic logger

- `tests`
    - Unit and integration tests related to the application.
 
## ⛑️ Contributing
 TODO

## 🚀 Local Run
TODO - update makefile
- Navigate to root directory of repository
- Run `go mod download` to download packages
- Run `make run` to run api

### Prerequisites
- Install Go version >= 1.23.0 https://go.dev/doc/install
