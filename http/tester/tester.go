package tester

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"go.mindeco.de/logging"
	"go.mindeco.de/logging/logtest"
)

type Tester struct {
	mux http.Handler
	t   *testing.T

	jar *cookiejar.Jar

	extraHeaders http.Header
}

func New(mux *http.ServeMux, t *testing.T) *Tester {
	l, _ := logtest.KitLogger("http/tester", t)
	tester := Tester{
		mux: logging.InjectHandler(l)(mux),
		t:   t,
	}

	var err error
	tester.jar, err = cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	tester.ClearHeaders()

	return &tester
}

func (t *Tester) ClearHeaders() {
	t.extraHeaders = make(http.Header)
}

func (t *Tester) SetHeaders(h http.Header) {
	for k, vals := range h {
		for _, v := range vals {
			t.extraHeaders.Add(k, v)
		}
	}
}

func (t *Tester) ClearCookies() {
	var err error
	t.jar, err = cookiejar.New(nil)
	if err != nil {
		t.t.Fatal("failed to clear cookies:", err)
	}
}

func (t *Tester) constructHeader(h *http.Header, u *url.URL) {
	*h = t.extraHeaders.Clone()

	cookies := t.jar.Cookies(u)
	for _, c := range cookies {
		cstr := c.String()
		h.Add("Cookie", cstr)
	}
}

func (t *Tester) GetHTML(u *url.URL) (*goquery.Document, *httptest.ResponseRecorder) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		t.t.Fatal(err)
	}
	t.constructHeader(&req.Header, u)

	rw := httptest.NewRecorder()
	t.mux.ServeHTTP(rw, req)

	t.jar.SetCookies(u, rw.Result().Cookies())

	doc, err := goquery.NewDocumentFromReader(rw.Body)
	if err != nil {
		t.t.Fatal(err)
	}

	return doc, rw
}

func (t *Tester) GetBody(u *url.URL) (rw *httptest.ResponseRecorder) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		t.t.Fatal(err)
	}
	t.constructHeader(&req.Header, u)

	rw = httptest.NewRecorder()
	t.mux.ServeHTTP(rw, req)

	t.jar.SetCookies(u, rw.Result().Cookies())
	return
}

func (t *Tester) GetJSON(u *url.URL, v interface{}) (rw *httptest.ResponseRecorder) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		t.t.Fatal(err)
	}
	t.constructHeader(&req.Header, u)

	rw = httptest.NewRecorder()
	t.mux.ServeHTTP(rw, req)
	t.jar.SetCookies(u, rw.Result().Cookies())

	body := rw.Body.Bytes()
	if rw.Code == 200 {
		if err = json.Unmarshal(body, v); err != nil {
			t.t.Log("Body:", string(body))
			t.t.Fatal(err)
		}
	}

	return
}

func (t *Tester) SendJSON(u *url.URL, v interface{}) (rw *httptest.ResponseRecorder) {
	blob, err := json.Marshal(v)
	if err != nil {
		t.t.Fatal(err)
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(blob))
	if err != nil {
		t.t.Fatal(err)
	}
	t.constructHeader(&req.Header, u)

	req.Header.Set("Content-Type", "application/json")

	rw = httptest.NewRecorder()
	t.mux.ServeHTTP(rw, req)
	t.jar.SetCookies(u, rw.Result().Cookies())
	return
}

func (t *Tester) PostForm(u *url.URL, v url.Values) (rw *httptest.ResponseRecorder) {
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(v.Encode()))
	if err != nil {
		t.t.Fatal(err)
	}
	t.constructHeader(&req.Header, u)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rw = httptest.NewRecorder()
	t.mux.ServeHTTP(rw, req)
	t.jar.SetCookies(u, rw.Result().Cookies())
	return
}
