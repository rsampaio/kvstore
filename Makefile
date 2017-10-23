certs:
	openssl genrsa -out cert.key
	openssl req -x509 -key cert.key -out cert.pem -subj "/C=US/ST=CA/L=San Francisco/O=kvstore.test"

install:
	go install -v ./...

test:
	go test -race -cover ./...

benchmarks:
	go test -race -bench=Server ./server

godoc:
	@echo "Open documentation at http://localhost:6060/pkg/github.com/rsampaio/kvstore"
	@godoc -http=:6060

run-tls: install certs
	@${GOPATH}/bin/kvserver --enable-tls --tls-cert=cert.pem --tls-key=cert.key
