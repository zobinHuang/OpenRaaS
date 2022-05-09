FROM golang:alpine as builder

# setup work directory
WORKDIR /go/src/app

# enable golang module
ENV GO111MODULE=on

# set up proxy in mainland China
RUN go env -w GOPROXY=https://goproxy.cn

# get reflex for live reload in dev
RUN go install github.com/cespare/reflex@latest

# copy go.mod & go.sum to WORKDIR
COPY go.mod .
COPY go.sum .

# download all necessary go module
RUN go mod download

# copy all source code to WORKDIR
COPY .. .

# build golang application
RUN go build -o ./run .

FROM alpine:latest

# update CA certificates for https
RUN apk --no-cache add ca-certificates

# setup work directory
WORKDIR /root/

# copy executable from builder
COPY --from=builder /go/src/app/run .

# expose ports
EXPOSE 8080

# run application
CMD ["./run"]