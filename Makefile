.PHONY: all build check deps help
all: check build

## build: build colordiff
build:
	go build .

## check: execute gofumpt and golangci-lint
check:
	gofumpt -e -extra -w .
	golangci-lint run -E asciicheck,bidichk,bodyclose,canonicalheader,containedctx,contextcheck,copyloopvar,decorder,dogsled,dupl,dupword,durationcheck,embeddedstructfieldcheck,err113,errcheck,errchkjson,errname,errorlint,exhaustive,exptostd,fatcontext,forcetypeassert,funcorder,gocheckcompilerdirectives,gochecknoglobals,gochecksumtype,gocognit,goconst,gocritic,gocyclo,godoclint,godot,godox,goheader,gomoddirectives,gomodguard,goprintffuncname,gosec,govet,grouper,iface,importas,inamedparam,ineffassign,interfacebloat,intrange,iotamixing,ireturn,lll,loggercheck,maintidx,makezero,mirror,misspell,mnd,musttag,nakedret,nestif,nilerr,nilnesserr,nilnil,noctx,nolintlint,nonamedreturns,nosprintfhostport,perfsprint,prealloc,predeclared,promlinter,protogetter,reassign,recvcheck,revive,rowserrcheck,sloglint,sqlclosecheck,tagalign,tagliatelle,testableexamples,thelper,tparallel,unconvert,unparam,unqueryvet,unused,usestdlibvars,usetesting,wastedassign,whitespace,wrapcheck,zerologlint --show-stats --color always | \less -iMRFX
	@printf "Press Enter to continue..."; read dummy

## deps: get linting dependencies, if necessary
deps:
	go mod tidy
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
