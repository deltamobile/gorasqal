include $(GOROOT)/src/Make.inc

TARG=bitbucket.org/ww/gorasqal
GOFILES=doc.go
CGOFILES=gorasqal.go
CGO_OFILES=crasqal.o

include $(GOROOT)/src/Make.pkg

format:
	gofmt -w *.go

docs:
	gomake clean
	godoc ${TARG} > README.txt
