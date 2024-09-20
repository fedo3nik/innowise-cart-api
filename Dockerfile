FROM golang:1.22.5 AS build-stage

WORKDIR /app    

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go test -v ./...

RUN CGO_ENABLED=0 G00S=Linux go build -o /docker-cart-api ./internal/cmd/main.go    

FROM build-stage AS run-test-stage
RUN go test -v ./...

EXPOSE 3000

ENTRYPOINT [ "/docker-cart-api" ]