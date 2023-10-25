# multistage docker build. This redices the size of the final docker image.
# stage 1 to build the app
FROM golang:1.21.3

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# stage 2 deploys the app built in stage 1
EXPOSE 3000

CMD ["./main"]