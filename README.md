# Project api-gateway

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing
purposes. See deployment for notes on how to deploy the project on a live system.

## Environment Configuration

Before running the application, you need to create an `.env` file in the root directory of the project. This file should contain all the necessary environment variables for the application to run properly.

Create a `.env` file in the project root:

```bash
touch .env
```

Add the required environment variables to your `.env` file. You can use the `.env.example` file as a template if available.

## MakeFile

Run the application locally

```bash
make run
```

Run the dockerized application

```bash
make docker-run
```

Shutdown the dockerized application

```bash
make docker-down
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Build the application

```bash
make build
```

Clean up binary from the last build:

```bash
make clean
```