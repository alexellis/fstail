LDFLAGS := "-s -w"

export GO111MODULE=on

.PHONY: dist
dist:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS)  -o bin/fstail
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS)  -o bin/fstail-darwin
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags $(LDFLAGS) -installsuffix cgo -o bin/fstail-darwin-arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS)  -o bin/fstail-arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOVER=7 go build -ldflags $(LDFLAGS)  -o bin/fstail-armhf
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS)  -o bin/fstail.exe
