# go-broker

Simple Message broker written in Go. Being lightweight, this can be used for testing pub-sub models, especially while writing a subscriber or publisher. 

Go-broker uses REST for publishing and requires subscriber API to notify (in **non-dev mode**, more ahead).

Please note the following information:

**main** branch is the database version of go-broker. If you do not want to use database and only work with memory, you can switch to branch **go-broker-no-database**.

## Requirements

We use postgres to store registered Publishers and Subscribers. Don't want to bother with setting up the database? Fear not! We got you covered. Use the dockerized installation as described in next section. 

However, we still require the following environment flags.
- IMAGE_NAME: Name of the go-broker image (default 'gobroker')
- POSTGRES_URL: Postgresql url (default 'localhost')
- POSTGRES_PORT: Postgres port (default '5432')
- POSTGRES_USER: Postgres user (default 'gobroker')
- POSTGRES_PASSWORD: Postgres user password (default 'admin')
- POSTGRES_DB: Postgres database (default 'postgres')
- JWT_KEY: Your secret JWT key

The default values can be used in docker compose setup. 

## Installation

Setting the Environment Variable

```bash
# In case of Linux/macOS
export JWT_KEY="your_secret_key"

# In case of Windows Command Prompt
set JWT_KEY=your_secret_key

# In case of Windows PowerShell
$env:JWT_KEY="your_secret_key"
```

Similarly, set the following environment flags: IMAGE_NAME, POSTGRES_URL, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, and POSTGRES_DB

### Running the Message Broker directly

Runnning the message broker is as easy as any other go app. Just use the command

```bash
go run cmd/go-broker/main.go
```

**NOTE:** In this case, make sure you have the postgres server running. 

Want to develop this code further? Say no more. It can be rather hectic to have a subscriber with a listener API always running. The solution? Running the app in development mode. How? Just add the dev flag in the run command.

```bash
go run cmd/go-broker/main.go --dev
```

In this mode, the app prints all the published messages instead of posting it to listener. For example:

```bash
Subscriber sub1 received message: {topic1 Hello, World!}
```

### Dockerized Installation (Preferred)

Dockerized installation is made easy for you. It is just like saying 1..2..3.. Use the following command for building the image and running all containers.

```bash
make all
```

## APIs

### POST APIs

- **Register Publisher**: 

```bash
curl -X POST -d '{"id":"pub1"}' -H "Content-Type: application/json" http://localhost:8080/register-publisher
```
  Response: 
```json
{
    "token": "JWT_TOKEN_HERE"
}
```

- **Register Subscriber**:

```bash
curl -X POST -d '{"id":"sub1", "topic":"topic1", "listener":"http://localhost:8081/listener"}' -H "Content-Type: application/json" http://localhost:8080/register-subscriber
```

Response:

```json
{
    "token": "JWT_TOKEN_HERE"
}
```

- Publish:

```bash
curl -X POST -d '{"topic":"topic1", "payload":"Hello, World!"}' -H "Content-Type: application/json" -H "Authorization: JWT_TOKEN_HERE" http://localhost:8080/publish
```
