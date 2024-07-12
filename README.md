# Argus

![Argus Banner](./docs/images/argus-banner.jpeg)

[![codecov](https://codecov.io/github/Amin-MAG/Argus/graph/badge.svg?token=AL8RSEOJ2C)](https://codecov.io/github/Amin-MAG/Argus)

## API Documentation

The `swagger` file (`api.yml`) exists in `docs` directory. You can see All APIs of Argus service in this file.

## Code Architecture

+ `api`: API documents like Swagger file.
+ `cmd` : Service entry points (Argus service).
+ `config` : Files containing configuration structs for Argus service.
+ `internal` : all application specific logic are implemented here.
  + `db`: Database schema and queries.
  + `handlers`: All Gin handlers.
  + `iputil`: Customized wrapper for IPInfo service.
  + `routes`: Creating Gin server and Routing different requests. 
+ `pkg`: General purpose packages like `logger`, `otel`.
+ `test`: Contains scripts for load testing

## System Architecture

![Argus Design](./docs/images/argus-design.jpg)

## Developers Guide

### Swagger

To update the swagger docs you can use

```bash
swag init -g ./cmd/argus/main.go -o api
```
### Formatting

```bash
# We need to exclude the vendor directory
find . -name '*.go' ! -path './vendor/*' -type f -exec go fmt {} \;
```

```bash
# We need to exclude the vendor directory
find . -name '*.go' ! -path './vendor/*' -type f | xargs goimports -w
```

### Git Hooks

You can create a pre-commit hook that automatically runs each one of linting commands before git commits. To enable the
project's git hook use this command.

```bash
cp githooks/pre-commit .git/hooks/pre-commit
```