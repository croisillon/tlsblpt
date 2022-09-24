package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
)

var options struct {
	listen     string
	serverCert string
	serverKey  string
	caCert     string
}

func main() {
	flag.StringVar(&options.listen, "listen", "0.0.0.0:5433", "Listen address")
	flag.StringVar(&options.serverCert, "sslcert", "cert.pem", "Server signed certificate")
	flag.StringVar(&options.serverKey, "sslkey", "key.pem", "Server key certificate")
	flag.StringVar(&options.caCert, "sslcacert", "ca.pem", "Certificate of the CA who signed server's certificate")
	flag.Parse()

	pemCACert, err := ioutil.ReadFile(options.caCert)
	if err != nil {
		log.Fatalf("failed to read ca certificate: %s", err.Error())
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemCACert) {
		log.Fatalf("failed to add CA's certificate")
	}

	cert, err := tls.LoadX509KeyPair(options.serverCert, options.serverKey)
	if err != nil {
		log.Fatalf("server: loadkeys %s", err.Error())
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	ln, err := tls.Listen("tcp", options.listen, tlsConfig)
	if err != nil {
		log.Fatalf("server: listen %s", err.Error())
	}

	log.Println("Listening on ", ln.Addr())

	defer ln.Close()

	for {
		clientConn, err := ln.Accept()
		if err != nil {
			log.Printf("%v", err)
			continue
		}

		log.Println("Connection from ", clientConn.RemoteAddr())

		go func(cl net.Conn) {
			buf := make([]byte, 1024)
			n, err := cl.Read(buf)
			if err != nil {
				if err == io.ErrUnexpectedEOF {
					log.Fatalln(io.ErrUnexpectedEOF)
				}
				log.Fatalf("error frontend read: %s", err)
			}

			log.Printf("%s\n", buf[:n])

			cl.Write([]byte("HELLO Mr. TLS Client"))
			cl.Close()

		}(clientConn)
	}
}
