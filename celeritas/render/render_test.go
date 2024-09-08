package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRender_Page(t *testing.T) {
	tests := []struct {
		name          string
		renderer      string
		template      string
		errorExpected bool
		errorMessage  string
	}{
		{"go_page", "go", "home", false, "error rendering go template"},
		{"go_page_no_template", "go", "no-file", true, "no error rendering non-existent go template, when one is expected"},
		{"jet_page", "go", "home", false, "error rendering jet template"},
		{"jet_page_no_template", "go", "no-file", true, "no error rendering non-existent jet template, when one is expected"},
		{"invalid_render_engine", "foo", "home", true, "no error rendering with non-existent template engine"},
	}

	for _, tt := range tests {
		r, err := http.NewRequest("GET", "some-url", nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()

		testRenderer.Renderer = tt.renderer
		testRenderer.RootPath = "./testdata"

		err = testRenderer.Page(w, r, tt.template, nil, nil)
		if tt.errorExpected {
			if err == nil {
				t.Errorf("%s: %s", tt.name, tt.errorMessage)
			}
		} else {
			if err != nil {
				t.Errorf("%s: %s: %s", tt.name, tt.errorMessage, err.Error())
			}
		}
	}
}

func TestRender_GoPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/url", nil)
	if err != nil {
		t.Error(err)
	}

	testRenderer.Renderer = "go"
	testRenderer.RootPath = "./testdata"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error rendering page", err)
	}
}

func TestRender_JetPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/url", nil)
	if err != nil {
		t.Error(err)
	}

	testRenderer.Renderer = "jet"
	testRenderer.RootPath = "./testdata"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error rendering page", err)
	}
}
