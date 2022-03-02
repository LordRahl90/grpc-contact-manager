FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o grpc-contact ./cmd


FROM gcr.io/distroless/base-debian10
# FROM scratch

WORKDIR /
COPY --from=build /app/grpc-contact grpc-contact


EXPOSE 3500
EXPOSE 3501
EXPOSE 5200

ENTRYPOINT ["/grpc-contact"]