# Fret Files

```
my-go-project/
│
├── cmd/                # Main applications for this project
│   ├── myapp/         # Application-specific code
│   │   └── main.go    # Entry point of the application
│   └── anotherapp/
│       └── main.go
│
├── internal/          # Private application and library code
│   ├── feature1/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   └── feature2/
│       ├── handler.go
│       ├── service.go
│       └── model.go
│
├── pkg/               # Public libraries for use by other projects
│   └── utils/
│       └── utils.go
│
├── api/               # API definitions (if applicable)
│   └── routes.go      # API routes definitions
│
├── tests/             # Test files and directories
│   └── feature1_test.go  # Tests for feature1
│
├── go.mod             # Module definition file
└── README.md          # Project documentation
```
