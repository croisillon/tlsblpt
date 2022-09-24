# pgoxy
Golang TLS Client/Server boilerplate (with SSC)

## How to run

```sh 
$ cd tls
$ bash generate-certificate.sh
```

### Server
```sh
$ cd cmd/server
$ go run main.go -sslcert=../../certs/server/cert.pem -sslkey=../../certs/server/key.pem -sslcacert=../../certs/ca/cert.pem
```

### Client
```sh
$ cd cmd/client
$ go run main.go -dial=127.0.0.1:5433 -sslcert=../../certs/server/cert.pem -sslkey=../../certs/server/key.pem -sslcacert=../../certs/ca/cert.pem
```