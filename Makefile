.PHONY: certgen output certgen_debug dumbc2 dumbc2_debug all debugserver all_platform

DEBUGFLAGS = -race
DEBUGENV = CGO_ENABLED=1
RELEASEFLAGS = -trimpath -ldflags '-s -w'
RELEASEENV = CGO_ENABLED=0

debugserver:
	cd output; python3 -m http.server 8080 &

all: clean certgen copycerts generate dumbc2 agent

all_platform: generate
	OUTPUT=./output/certgen ./scripts/build-all.sh cmd/certgen/main.go
	OUTPUT=./output/dumbc2	./scripts/build-all.sh cmd/control/main.go
	echo "Due the TLS Authentication, If you need agent, compile it yourself." > ./output/dumbc2_agent
	bash ./scripts/archive-release.sh
	bash ./scripts/sha256sums.sh

copycerts:
	./output/certgen
	cp -af ~/.dumbyc2/cacert.pem buildtime/certs/
	cp -af ~/.dumbyc2/clientpk.pem buildtime/certs/
	cp -af ~/.dumbyc2/clientcert.pem buildtime/certs/
	cp -af ~/.dumbyc2/serverpin.txt buildtime/certs/

clean:
	rm -rf output/*
	rm -rf statik/*

generate:
	go get -u github.com/rakyll/statik
	go generate cmd/agent/main.go

prune:    clean
	rm -rf buildtime/certs/*.pem
	rm -rf buildtime/certs/*.txt
	rm -rf ~/.dumbyc2

certgen:
	${RELEASEENV} go build ${RELEASEFLAGS} -o output/certgen cmd/certgen/main.go

dumbc2:
	GOOS=linux GOARCH=amd64 ${RELEASEENV} go build ${RELEASEFLAGS} -o output/dumbyc2 cmd/control/main.go

agent: generate
	GOOS=linux GOARCH=amd64 ${RELEASEENV} go build ${RELEASEFLAGS} -o output/dumbyc2_agent cmd/agent/main.go

dumbc2_debug:
	${DEBUGENV} go build ${DEBUGFLAGS} -o output/dumbyc2 cmd/control/main.go

certgen_debug:
	${DEBUGENV} go build ${DEBUGFLAGS} -o output/certgen cmd/certgen/main.go

agent_debug:    generate
	${DEBUGENV} go build ${DEBUGFLAGS} -o output/dumbyc2_agent cmd/agent/main.go
