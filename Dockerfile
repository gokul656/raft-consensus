FROM golang:1.21

WORKDIR /data/trade-engine

RUN apt-get update
RUN apt-get install -y protobuf-compiler protoc-gen-go protoc-gen-go-grpc

# Copying from host to image, since WORKDIR is already defined, the second dot denotes /data/trade-engine
COPY . .

RUN make build

CMD ./bin/engine
