cd api/proto/v1/ &&
protoc --go_out=plugins=grpc:. *.proto
protoc --grpc-gateway_out=logtostderr=true:. *.proto
protoc --swagger_out=logtostderr=true:../swagger *.proto