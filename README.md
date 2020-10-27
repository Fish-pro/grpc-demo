# grpc-demo
It's a `gRPC` demo, include server client and client-rest, run server, can start `gRPCServer` and `httpServer`
## make .go file
```bash
sh gen.sh
```
## Certificate
### Mutual authentication
Generate root certificate
```bash
genrsa -out ca.key 2048
req -new -x509 -days 3650 -key ca.key -out ca.pem
```
server certificate
```bash
genrsa -out server.key 2048
req -new -key server.key -out server.csr // 本地使用localhost作为域名
x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in server.csr -out server.pem
```
client certificate
```bash
ecparam -genkey -name secp384r1 -out client.key
req -new -key client.key -out client.csr
x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.pem
```
## Swagger
Run server and access [apiList](http://localhost:8080/api/)
