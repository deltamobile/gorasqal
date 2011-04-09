PACKAGE

package gorasqal
import "bitbucket.org/ww/gorasqal"

Rasqal bindings for Go. Written by William Waites in 2011.

SPARQL and related RDF query language bindings and utilities.

These are minimal and reflect only the parts of rasqal that I
need at the moment. Contributions to make it more complete are
more than welcome.


FUNCTIONS

func GoRasqal_log_handler(user_data unsafe.Pointer, messagep unsafe.Pointer)
export GoRasqal_log_handler

func QueryPrint(query string)


TYPES

type Query struct {
    // contains unexported fields
}

func NewQuery(w *World) *Query

func (q *Query) Free()

func (q *Query) Prepare(query string) (err os.Error)

func (q *Query) Print()

type Service struct {
    // contains unexported fields
}

func NewService(world *World, endpoint string, query string) *Service

func (s *Service) Execute() (results chan map[string]goraptor.Term)

func (s *Service) Free()

func (s *Service) SetFormat(format string)

type World struct {
    // contains unexported fields
}

func NewWorld() (w *World)

func (w *World) Free()


SUBDIRECTORIES

	.hg
