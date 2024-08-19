FROM golang:latest

WORKDIR /app

COPY . .

RUN git config --global --add safe.directory /app
RUN go mod tidy
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN apt update && apt install -y postgresql-client


EXPOSE 8000