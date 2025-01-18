# VeritasEngine
This is a flexible and customizable rule engine designed to accomodate the
rule evaluation needs of other systems.

## Technologies
- Go version 1.19
- Gin (web framework)
- Grule (rule engline library)
- PostgreSQL

# Usage
When running the application, you can supply a `mode` flag (optional) with values `dev` or `prod` to indiciate whether to run in production mode or development mode.

The system will assume the `dev` mode in the absence of the flag.

## Running in producion mode
In development mode, the system will use different configurations as defined in the `config.dev.yaml` file and make optimizations that are suitable for development.

To run the application in development mode, execute the following make utility command:

```bash
make run-dev
```

Or

```bash
go run cmd/main.go -mode=dev
```

## Running in production mode
In production mode, the system will use different configurations as defined in the `config.prod.yaml` file and make optimizations that are suitable for production ready systems.

To run the application in production mode, execute the following make utility command:

```bash
make run-prod
```

Or

```bash
go run cmd/main.go -mode=prod
```
