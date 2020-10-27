# grpc-demo
It's a `gRPC` demo, include server client and client-rest, run server, can start `gRPCServer` and `httpServer`
## Run server
```bash
make run
```
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
Server certificate
```bash
genrsa -out server.key 2048
req -new -key server.key -out server.csr // 本地使用localhost作为域名
x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in server.csr -out server.pem
```
Client certificate
```bash
ecparam -genkey -name secp384r1 -out client.key
req -new -key client.key -out client.csr
x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.pem
```
## Swagger
Run server and access <http://localhost:8080/api/>
## Configuration
The configuration can be modified by configuring environment variables

| Environment variable name | Description                      | Default                     |
| ------------------------- | -------------------------------- | --------------------------- |
| RUN_HOST                  | host name                        | localhost                   |
| GRPC_PORT                 | gRPC server port                 | 8081                        |
| HTTP_PORT                 | http srever port                 | 8080                        |
| OPEN_PEM                  | if Mutual authentication         | true                        |
| LOG_LEVEL                 | log level[-1(debug)--->3(error)] | -1                          |
| TIME_FORMAT               | log time format                  | 2006-01-02 15:04:05         |
| GRPC_DB_HOST              | database host                    | 127.0.0.1:3306              |
| GRPC_DB_USER              | database user                    |                             |
| GRPC_DB_PASSWORD          | database password                |                             |
| GRPC_DB_NAME              | database name                    | grpc-demo                   |
| CA_PATH                   | root certificate path            | ~/grpc-demo/cert/ca.pem     |
| PEM_PATH                  | server pem path                  | ~/grpc-demo/cert/server.pem |
| KEY_PATH                  | server key path                  | ~/grpc-demo/cert/server.key |

