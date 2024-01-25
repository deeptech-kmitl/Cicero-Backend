FROM golang:1.21-bullseye AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o myapp main.go

## Deploy
FROM gcr.io/distroless/static-debian11

COPY --from=build /app/myapp /bin
COPY .env.prod /bin

EXPOSE 3000

ENTRYPOINT [ "/bin/myapp", "/bin/.env.prod" ]