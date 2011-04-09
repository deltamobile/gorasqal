package gorasqal

import (
	"log"
	"testing"
)

/*
func TestPrint(t *testing.T) {
	QueryPrint("SELECT DISTINCT ?s WHERE { { ?s <a> ?o } UNION { ?s <b> <c> } }")
}
*/

func TestService(t *testing.T) {
	w := NewWorld()
	query := `
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
SELECT * 
WHERE { 
    ?s rdfs:comment ?o
    FILTER ( isLiteral(?o) )
}
LIMIT 2
`
	svc := NewService(w, "http://dbpedia.org/sparql", query)
	defer svc.Free()

	if svc == nil {
		t.Fatal("could not make service")
	}

	results := svc.Execute()
	for {
		row, ok := <-results
		if !ok {
			break
		}
		log.Print(row)
	}
}