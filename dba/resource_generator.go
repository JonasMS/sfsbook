package dba

import (
	"log"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
)

// ResourceResultsGenerator finds a specific resource by its uuid.
type ResourceResultsGenerator struct {
	index bleve.Index
}

func MakeResourceResultsGenerator(bi bleve.Index) *ResourceResultsGenerator {
	return &ResourceResultsGenerator{index: bi}
}

type resourceResults struct {
	// Indicates if the query was successful. (i.e. that it produced data.)
	Success bool

	// Show this if something went wrong.
	FailureText string

	// The requested resource.
	Uuid string

	// The actual fields in the document.
	Document map[string]interface{}

	// Should we display additional debugging info
	Debug bool
}

// ForRequest generates the data comprising a result page showing a single
// resource guide entry.
func (qr *ResourceResultsGenerator) ForRequest(req interface{}) interface{} {
	uuid := req.(string)

	log.Println("uuid", uuid)

	// Code quality comment: Writing the templates requires knowing what I've
	// produced here. I feel that I have not layered this code very well.
	results := &resourceResults{
		Success: false,
		FailureText: "query had a sad",
		Uuid: uuid,
		Document: make(map[string]interface{}),
		Debug: true,
	}	

	doc, err :=  qr.index.Document(uuid)
	if err != nil || doc == nil {
		log.Println("query failed", err)
		return results
	}

	log.Println("succeeded", doc)
	results.Success = true
	
	// Because template code operates on maps, I can build a generic solution that
	// can work for any future change in the format of documents.
	for _, f := range doc.Fields {
		// TODO(rjk): I want the debug output to be dynamically configurable.
		// log.Println("resource name", f.Name())

		switch t := f.(type) {
		default:
			log.Println("object",  f, "has unexpected type");		

		case *document.NumericField:
			// It is preferable to ship a number here so that I generate JSON.
			v, err := t.Number()
			if err != nil {
				// I don't really know what goes here. I need to figure it out.
				// Error handling in general needs to be treated correctly.
				log.Println("couldn't convert field", t.Name(), "to number", err)				
				continue
			}
			results.Document[t.Name()] = v
		case *document.BooleanField:
			// log.Println("found boolean field", f.Name())
			v, err := t.Boolean()
			if err != nil {
				log.Println("couldn't convert field", t.Name(), "to bool", err)	
				continue
			}
			results.Document[t.Name()] = v
		case *document.TextField:
			results.Document[t.Name()] = string(t.Value())
		}
	}
	
	// Need to support showing the comments.
	return results
}
