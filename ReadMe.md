# GridStream API

api to allow frontend to communicate with database


# Directory Structure


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
│   │   └── logic/
│   └── config/
│       └── config.go
├── pkg/
└── tests/
    ├── integration/
    └── unit/
```

**cmd/api/main.go**: Entry point of the application where we set up the HTTP server and wire up components.  

**internal/app**: Contains the core application logic, organized into subfolders:  
    &emsp; **handlers**: Handle the HTTP request and maps the route to business logic  
    &emsp; **middlewares**: Optional functions that can be applied to HTTP requests. For example, making sure a user is authenticated. (Adapter pattern).  
    &emsp; **repositories**: Data access layer that interacts directly with the database.  
    &emsp; **logic**: Holds business logic that needs to be applied to data  
    
**internal/config**: Holds Application configuration, such as database connection details, API keys, etc.  

**pkg**: Holds reusuable libraries to be used accross different projects, may not be necessary in our case  

**tests**: Unit and integration tests related to the application.  
