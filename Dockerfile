FROM golang:latest
LABEL maintainer="Kyle Keller <kylegk@gmail.com>"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o goapp .
EXPOSE 8080
ENV SSE_SERVER_URL="https://live-test-scores.herokuapp.com/scores"
ENV APPLICATION_PORT=":8080"
CMD ["./goapp"]