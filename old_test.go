package swaggering_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/savaki/swaggering"
)

type Once struct {
	counter int32
	once    sync.Once
}

func (o *Once) Increment() {
	o.once.Do(func() {
		atomic.AddInt32(&o.counter, 1)
	})
}

func TestOnce(t *testing.T) {
	o := &Once{}

	o.Increment()
	if o.counter != 1 {
		t.Errorf("expected counter == 0 ; got %v", o.counter)
		return
	}

	o.Increment()
	if o.counter != 1 {
		t.Errorf("expected counter == 0 ; got %v", o.counter)
		return
	}
}

func TestServeHTTP_SetsProducerAndConsumer(t *testing.T) {
	api := &swaggering.OldApi{
		Endpoints: swaggering.OldEndpoints{
			{
				Method: "get",
				Path:   "/",
			},
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	api.ServeHTTP(w, req)

	content := &swaggering.Api{}
	err := json.Unmarshal(w.Body.Bytes(), content)
	if err != nil {
		t.Errorf("unable to unmarshal content - %v", err)
		return
	}

	endpoint := content.Paths["/"].Get
	expected := "application/json"

	// check produces
	//
	if got := len(endpoint.Produces); got != 1 {
		t.Errorf("expected produces to default to application/json - %v", err)
		return
	}
	if got := endpoint.Produces[0]; got != expected {
		t.Errorf("expected %v, got %v", expected, got)
		return
	}

	// check consumes
	//
	if got := len(endpoint.Consumes); got != 1 {
		t.Errorf("expected consumes to default to application/json - %v", err)
		return
	}
	if got := endpoint.Consumes[0]; got != expected {
		t.Errorf("expected %v, got %v", expected, got)
		return
	}
}

func TestServeHTTP_SetsHostAndScheme(t *testing.T) {
	api := &swaggering.OldApi{}

	// test from one host
	//

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	api.ServeHTTP(w, req)

	content := &swaggering.Api{}
	err := json.Unmarshal(w.Body.Bytes(), content)
	if err != nil {
		t.Errorf("unable to unmarshal content - %v", err)
		return
	}

	if got := len(content.Schemes); got != 1 {
		t.Errorf("expected %#v; got %#v", 1, got)
		return
	}

	if expected := "https"; content.Schemes[0] != expected {
		t.Errorf("expected %#v; got %#v", expected, content.Schemes)
		return
	}

	if expected := "example.com"; expected != content.Host {
		t.Errorf("expected %v; got %v", expected, content.Host)
		return
	}

	// and then another
	//

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "http://othersite.com", nil)
	api.ServeHTTP(w, req)

	content = &swaggering.Api{}
	err = json.Unmarshal(w.Body.Bytes(), content)
	if err != nil {
		t.Errorf("unable to unmarshal content - %v", err)
		return
	}

	if got := len(content.Schemes); got != 1 {
		t.Errorf("expected %#v; got %#v", 1, got)
		return
	}

	if expected := "http"; content.Schemes[0] != expected {
		t.Errorf("expected %#v; got %#v", expected, content.Schemes)
		return
	}

	if expected := "othersite.com"; expected != content.Host {
		t.Errorf("expected %v; got %v", expected, content.Host)
		return
	}
}
