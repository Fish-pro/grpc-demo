cd api/proto/v1/
protoc --go_out=plugins=grpc:. *.proto