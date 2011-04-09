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
The rasqal service enables queries against remote SPARQL endpoints.
It is a one-off construct, used to execute a single query.

func NewService(world *World, endpoint string, query string) *Service

func (s *Service) Execute() (err os.Error)
Perform the operation as a query and return a set of results. This is usually
used for SPARUL INSERT/DELETE queries.

func (s *Service) Free()

func (s *Service) Query() (results chan map[string]goraptor.Term, err os.Error)
Perform the operation as a query and return a set of results.

func (s *Service) SetFormat(format string)

func (s *Service) SetProxy(proxy string)

func (s *Service) SetUserAgent(user_agent string)

type World struct {
    // contains unexported fields
}

func NewWorld() (w *World)

func (w *World) Free()


SUBDIRECTORIES

	.hg
