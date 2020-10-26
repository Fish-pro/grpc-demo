package helper

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/Fish-pro/grpc-demo/config"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

// 获取服务端证书
func GetServerCred(c *config.Certificate) credentials.TransportCredentials {
	cert, _ := tls.LoadX509KeyPair(c.PemPath, c.KeyPath)
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile(c.CaPath)
	certPool.AppendCertsFromPEM(ca)

	cred := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},        // 服务端证书
		ClientAuth:   tls.RequireAndVerifyClientCert, // 双向验证
		ClientCAs:    certPool,
	})
	return cred
}

// 获取客户端证书
func GetClientCred() credentials.TransportCredentials {
	cert, _ := tls.LoadX509KeyPair(
		"/Users/york/go/src/github.com/Fish-pro/grpc-demo/cert/client.pem",
		"/Users/york/go/src/github.com/Fish-pro/grpc-demo/cert/client.key",
	)
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("/Users/york/go/src/github.com/Fish-pro/grpc-demo/cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	cred := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert}, // 客户端证书
		ServerName:   "localhost",             // 域名
		RootCAs:      certPool,
	})
	return cred
}
