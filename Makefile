
TARG=bitbucket.org/ww/gorasqal
GOFILES=doc.go
CGOFILES=gorasqal.go
CGO_OFILES=crasqal.o

format:
	gofmt -w *.go

docs:
	godoc ${TARG} > README.txt
