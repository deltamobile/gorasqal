package gorasqal

// #cgo CFLAGS: -I/usr/include/rasqal -I/usr/include/raptor2
// #cgo LDFLAGS: -lrasqal -lraptor2
// #include <stdlib.h>
// #include <rasqal.h>
// #include "crasqal.h"
import "C"

import (
	"bytes"
	"errors"
	"github.com/deltamobile/goraptor"
	"log"
	"os"
	"sync"
	"text/template"
	"unsafe"
)

type World struct {
	rasqal_world *C.rasqal_world
	err          error
}

func NewWorld() (w *World) {
	rasqal_world := C.rasqal_new_world()
	if rasqal_world == nil {
		return
	}
	if C.rasqal_world_open(rasqal_world) != 0 {
		return
	}
	w = &World{rasqal_world: rasqal_world}
	C.gorasqal_set_log_handler(rasqal_world, unsafe.Pointer(w))
	return
}

//export GoRasqal_log_handler
func GoRasqal_log_handler(user_data unsafe.Pointer, messagep unsafe.Pointer) {
	message := (*C.raptor_log_message)(messagep)
	world := (*World)(user_data)
	text := C.GoString(message.text)
	if C.int(message.level) > C.RAPTOR_LOG_LEVEL_INFO {
		world.err = errors.New(text)
	}
	log.Print(text)
}

func (w *World) Free() {
	C.rasqal_free_world(w.rasqal_world)
}

type Query struct {
	world *World
	query *C.rasqal_query
}

func NewQuery(w *World) *Query {
	qlang := C.CString("sparql11")
	rq := C.rasqal_new_query(w.rasqal_world, qlang, (*C.uchar)(nil))
	C.free(unsafe.Pointer(qlang))
	return &Query{w, rq}
}

func (q *Query) Free() {
	C.rasqal_free_query(q.query)
}

func (q *Query) Prepare(query string) (err error) {
	cq := C.CString(query)
	result := C.rasqal_query_prepare(q.query, (*C.uchar)(unsafe.Pointer(cq)), (*C.raptor_uri)(nil))
	C.free(unsafe.Pointer(cq))
	if result != 0 {
		err = q.world.err
	}
	return
}

func (q *Query) Print() {
	mode := C.CString("w")
	ofp := C.fdopen(C.int(os.Stdout.Fd()), mode)
	C.free(unsafe.Pointer(mode))
	C.rasqal_query_print(q.query, ofp)
	C.fclose(ofp)
}

func QueryPrint(query string) {
	w := NewWorld()
	q := NewQuery(w)
	err := q.Prepare(query)
	if err == nil {
		q.Print()
	}
	q.Free()
	w.Free()
}

// The rasqal service enables queries against remote SPARQL endpoints.
// It is a one-off construct, used to execute a single query.
type Service struct {
	mutex        sync.Mutex
	world        *World
	endpoint     *C.raptor_uri
	endpoint_str string
	orig_query   string
	actual_query string
	svc          *C.rasqal_service
	dg           *C.raptor_sequence
	www          *C.raptor_www
	format		 string
	user_agent	 string
	proxy		 string
	query_ready	 bool
}

func NewService(world *World, endpoint string, query string) *Service {
	s := &Service{world: world}
	s.endpoint_str = endpoint
	s.orig_query = query
	return s
}

func (s *Service) generateQueryFromTemplate(data interface{}) error {
	buf := new(bytes.Buffer)
	tmpl := template.Must(template.New("query").Parse(s.orig_query))
	err := tmpl.Execute(buf, data)
	if err != nil {
		return err
	}
	s.actual_query = buf.String()
	return nil
}


/* Set up C rasqal object based on current values. */
func (s *Service) prepQuery(data ... interface{}) (err error) {

	if len(data) == 1 {
		s.query_ready = false
		err = s.generateQueryFromTemplate(data[0])
		if err != nil {
			log.Println("Generation of query from template failed. ", err)
			return
		}
	} else if len(data) == 0 && s.actual_query == "" {
		s.query_ready = false
		s.actual_query = s.orig_query
	} else if s.query_ready {
		return nil
	}

	s.free()
	raptor_world := C.rasqal_world_get_raptor(s.world.rasqal_world)

	cep := (*C.uchar)(unsafe.Pointer(C.CString(s.endpoint_str)))
	s.endpoint = C.raptor_new_uri(raptor_world, cep)
	C.free(unsafe.Pointer(cep))

	cquery := (*C.uchar)(unsafe.Pointer(C.CString(s.actual_query)))
	defer C.free(unsafe.Pointer(cquery))

	s.dg = C.goraptor_new_sequence()
	s.svc = C.rasqal_new_service(s.world.rasqal_world, s.endpoint, cquery, s.dg)

	if s.svc == nil {
		C.raptor_free_uri(s.endpoint)
		return errors.New("Failed to create service.")
	}

	s.www = C.raptor_new_www(raptor_world)
	if s.www == nil {
		s.free()
		return errors.New("Failed to create www.")
	}

	if C.rasqal_service_set_www(s.svc, s.www) != 0 {
		s.free()
		return errors.New("Failed to set www.")
	}

	if s.format != "" {
		cformat := C.CString(s.format)
		C.rasqal_service_set_format(s.svc, cformat)
		C.free(unsafe.Pointer(cformat))
	}

	if s.user_agent == "" {
		s.SetUserAgent("gorasqal hello world")
	}
	cua := C.CString(s.user_agent)
	C.raptor_www_set_user_agent(s.www, cua)
	C.free(unsafe.Pointer(cua))

	if s.proxy != "" {
		cproxy := C.CString(s.proxy)
		C.raptor_www_set_proxy(s.www, cproxy)
		C.free(unsafe.Pointer(cproxy))
	}
	s.query_ready = true

	return nil
}

func (s *Service) Free() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.free()
}

func (s *Service) free() {
	if s.svc != nil {
		C.rasqal_free_service(s.svc)
		s.svc = nil
	}
	if s.endpoint != nil {
		s.endpoint = nil /* C object freed by rasqal_free_service */
	}
	if s.dg != nil {
		s.dg = nil /* C object freed by rasqal_free_service */
	}
	/*
		The rasqal_free_service does not free the www object, but
		trying to do so here causes a crash.  Not sure why.

		if s.www != nil {
			C.raptor_free_www(s.www)
		}

	*/
	s.www = nil /* Questions: Is it OK to re-use the www object? */
	s.query_ready = false
}

func (s *Service) SetFormat(format string) {
	s.format = format
	s.query_ready = false
}

func (s *Service) SetUserAgent(user_agent string) {
	s.user_agent = user_agent
	s.query_ready = false
}

func (s *Service) SetProxy(proxy string) {
	s.proxy = proxy
	s.query_ready = false
}

// Perform the operation as a query and return a set of results. This is usually
// used for SPARUL INSERT/DELETE queries.
func (s *Service) Execute(data ... interface {}) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(data) > 0 {
		s.query_ready = false
	}
	if ! s.query_ready {
		err = s.prepQuery(data...)
		if err != nil {
			log.Println("Could not prepare query: ", err)
			return
		}
	}

	query_results := C.rasqal_service_execute(s.svc)
	if query_results == nil {
		// xxx when this fails, svc gets freed???
		s.svc = nil
		err = errors.New("could not execute the query. inspect the log for details")
	}
	return
}

// Perform the operation as a query and return a set of results.
func (s *Service) Query(data ... interface {}) (results []map[string]goraptor.Term, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(data) > 0 {
		s.query_ready = false
	}
	if ! s.query_ready {
		err = s.prepQuery(data...)
		if err != nil {
			log.Println("Could not prepare query.")
			return
		}
	}

	query_results := C.rasqal_service_execute(s.svc)
	if query_results == nil {
		// xxx when this fails, svc gets freed???
		s.free()
		err = errors.New("could not execute the query. inspect the log for details")
		return
	}

	rows := int(C.rasqal_query_results_get_count(query_results))
	results = make([]map[string]goraptor.Term, 0, rows)

	columns := int(C.rasqal_query_results_get_bindings_count(query_results))
	bindings := make([]string, 0, columns)
	for i := 0; i < columns; i++ {
		ucbinding := C.rasqal_query_results_get_binding_name(query_results, C.int(i))
		binding := C.GoString((*C.char)(unsafe.Pointer(ucbinding)))
		bindings = append(bindings, binding)
	}

	for {
		if C.rasqal_query_results_finished(query_results) != 0 {
			break
		}

		row := make(map[string]goraptor.Term)
		for i := 0; i < columns; i++ {
			rasqal_literal := C.rasqal_query_results_get_binding_value(query_results, C.int(i))
			ucvalue := C.rasqal_literal_as_string(rasqal_literal)
			value := C.GoString((*C.char)(unsafe.Pointer(ucvalue)))
			term_type := C.rasqal_literal_get_rdf_term_type(rasqal_literal)

			var term goraptor.Term
			switch {
			case term_type == C.RASQAL_LITERAL_BLANK:
				blank := goraptor.Blank(value)
				term = &blank
			case term_type == C.RASQAL_LITERAL_URI:
				uri := goraptor.Uri(value)
				term = &uri
			default: // literal
				dturi := C.rasqal_literal_datatype(rasqal_literal)
				var datatype string
				if dturi != nil {
					dtstr := C.raptor_uri_as_string(dturi)
					datatype = C.GoString((*C.char)(unsafe.Pointer(dtstr)))
				}
				var language string
				if rasqal_literal.language != nil {
					language = C.GoString(rasqal_literal.language)
				}
				term = &goraptor.Literal{Value: value, Lang: language, Datatype: datatype}
			}

			row[bindings[i]] = term
		}

		results = append(results, row)

		if C.rasqal_query_results_next(query_results) != 0 {
			break
		}
	}

	C.rasqal_free_query_results(query_results)
	return
}
