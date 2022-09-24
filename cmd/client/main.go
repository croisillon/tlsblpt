package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"
)

var options struct {
	dial       string
	clientCert string
	clientKey  string
	caCert     string
}

const (
	TxDialTimeout   = 10 * time.Second
	TxDialKeepAlive = 30 * time.Second
	TxDialDual      = false
	TxDialNetwork   = "tcp4"
)

func main() {
	flag.StringVar(&options.dial, "dial", "0.0.0.0:5433", "Dial address")
	flag.StringVar(&options.clientCert, "sslcert", "cert.pem", "Client signed certificate")
	flag.StringVar(&options.clientKey, "sslkey", "key.pem", "Client key certificate")
	flag.StringVar(&options.caCert, "sslcacert", "ca.pem", "Certificate of the CA who signed client's certificate")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pemCACert, err := ioutil.ReadFile(options.caCert)
	if err != nil {
		log.Fatalf("failed to read ca certificate: %s", err.Error())
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemCACert) {
		log.Fatalf("failed to add CA's certificate")
	}

	cert, err := tls.LoadX509KeyPair(options.clientCert, options.clientKey)
	if err != nil {
		log.Fatalf("client: loadkeys %s", err.Error())
	}

	tlsDialer := (tls.Dialer{
		NetDialer: &net.Dialer{
			Timeout:   TxDialTimeout,
			KeepAlive: TxDialKeepAlive,
			DualStack: TxDialDual,
		},
		Config: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      certPool,
		},
	})

	tlsConn, err := tlsDialer.DialContext(ctx, TxDialNetwork, options.dial)
	if err != nil {
		log.Fatalf("TLS dial failed: %s", err.Error())
	}

	_, _ = tlsConn.Write([]byte("HELLO Mr TLS Server"))

	buf := make([]byte, 1024)
	n, _ := tlsConn.Read(buf)

	fmt.Printf("%s\n", buf[:n])
	tlsConn.Close()
}
