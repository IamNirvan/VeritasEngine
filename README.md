# VeritasEngine

This is a flexible and customizable rule engine designed to accomodate the
rule evaluation needs of other systems.

## Technologies
- Go version 1.19
- Gin (web framework)
- Grule (rule engline library)
- PostgreSQL

# Usage

When running the application, you can supply a `mode` flag (optional) to indiciate whether to run in production mode or development mode.
The system will assume the dev mode if nothin is provided.

To run the application, execute the following make utility command:

```bash
make run
```

Alternatively, you can this command the start the service. This starts it in `dev` mode. Use `prod` to run in production mode

```bash
go run cmd/main.go -mode=dev
```
