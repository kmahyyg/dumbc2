.PHONY: certgen output certgen_debug dumbc2 dumbc2_debug

DEBUGFLAGS = -race
DEBUGENV = CGO_ENABLED=1
RELEASEFLAGS = -trimpath -ldflags '-s -w'
RELEASEENV = CGO_ENABLED=0

clean:
	rm -rf output/*

certgen:
	${RELEASEENV} go build ${RELEASEFLAGS} -o output/certgen cmd/certgen/main.go

certgen_debug:
	${DEBUGENV} go build ${DEBUGFLAGS} -o output/certgen cmd/certgen/main.go

dumbc2_debug:
	${DEBUGENV} go build ${DEBUGFLAGS} -o output/dumbyc2 cmd/control/main.go

dumbc2:
	${RELEASEENV} go build ${RELEASEFLAGS} -o output/dumbyc2 cmd/control/main.go
