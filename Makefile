.PHONY: certgen output certgen_debug dumbc2 dumbc2_debug

DEBUGFLAGS = -race
DEBUGENV = CGO_ENABLED=1
RELEASEFLAGS = -trimpath -ldflags '-s -w'
RELEASEENV = CGO_ENABLED=0

clean:
	rm -rf output/*

generate:
	go generate cmd/agent/main.go

prune:    clean
	rm -rf buildtime/certs/*.pem
	rm -rf buildtime/certs/*.txt
	rm -rf statik/

certgen:
	${RELEASEENV} go build ${RELEASEFLAGS} -o output/certgen cmd/certgen/main.go
	. output/certgen

dumbc2:
	${RELEASEENV} go build ${RELEASEFLAGS} -o output/dumbyc2 cmd/server/main.go

agent:
	${RELEASEENV} go build ${RELEASEFLAGS} -o output/dumbyc2_agent cmd/agent/main.go

dumbc2_debug:    generate
	${DEBUGENV} go build ${DEBUGFLAGS} -o output/dumbyc2 cmd/server/main.go

certgen_debug:
	${DEBUGENV} go build ${DEBUGFLAGS} -o output/certgen cmd/certgen/main.go

agent_debug:    generate
	${DEBUGENV} go build ${DEBUGFLAGS} -o output/dumbyc2_agent cmd/agent/main.go

