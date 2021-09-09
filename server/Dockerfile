FROM golang:alpine as build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o server

FROM alpine:latest
EXPOSE 4769
COPY --from=build /app/server .
CMD [ "./server" ]