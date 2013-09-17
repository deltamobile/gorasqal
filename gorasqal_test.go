package gorasqal

import (
	"log"
	"testing"
)

func TestPrint(t *testing.T) {
	QueryPrint("SELECT DISTINCT ?s WHERE { { ?s <a> ?o } UNION { ?s <b> <c> } }")
}

func TestService(t *testing.T) {
	log.Println("TestService basic started.")
	w := NewWorld()
	query := `
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
SELECT *
WHERE {
    ?s rdfs:comment ?o
}
LIMIT 2
`
	svc := NewService(w, "http://dbpedia.org/sparql", query)
	defer svc.Free()

	if svc == nil {
		t.Fatal("could not make service")
	}

	results, err := svc.Query()
	if err != nil {
		t.Fatal(err)
	}

	for _, row := range results {
		log.Print(row)
	}
	log.Println("TestService basic finished.\n")
}

type rdfVar  struct {
	Predicate string
}

func TestServiceTemplate(t *testing.T) {
	log.Println("TestService template started.")
	w := NewWorld()
	query := `
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
SELECT *
WHERE {
    ?s {{.Predicate}} ?o
}
LIMIT 2
`
	rv := []rdfVar{ rdfVar{"rdfs:label"}, rdfVar{"rdfs:comment"}, rdfVar{"rdfs:range"}}
	svc := NewService(w, "http://dbpedia.org/sparql", query)
	defer svc.Free()

	if svc == nil {
		t.Fatal("could not make service")
	}

	for i := 0; i < 3; i++ {

		results, err := svc.Query(rv[i])
		if err != nil {
			t.Fatal(err)
		}

		for _, row := range results {
			log.Print(row)
		}
	}
	log.Println("TestService template finished.\n")
}
