package test

import (
	"fmt"
	"github.com/qwerty22121998/go-util/hclient"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"testing"
)

type DataFormat struct {
	A string `json:"a"`
}

func Test_HTTPClient(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%s", b)
	}))
	defer ts.Close()
	ts.Client()
	remoteURL := ts.URL

	f, err := os.Open("a.txt")
	if err != nil {
		t.Fatal(err)
	}

	client := hclient.New(http.MethodPost, remoteURL, hclient.WithHTTPClient(ts.Client()))

	if err := client.FormFile(map[string]io.Reader{
		"file": f,
	}).Error(); err != nil {
		t.Fatal(err)
	}

}
