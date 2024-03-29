package logrus

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goroute/route"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := logrus.New()
	logger.SetOutput(buf)
	entry := logrus.NewEntry(logger)

	mux := route.NewServeMux()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := mux.NewContext(req, rec)
	h := func(c route.Context) error {
		return c.String(http.StatusOK, "test")
	}
	mw := New(Entry(entry))

	// Status 2xx
	mw(c, h)

	str := buf.String()
	assert.Contains(t, str, "bytes_out=4")
	assert.Contains(t, str, "host=example.com")
	assert.Contains(t, str, "method=GET")
	assert.Contains(t, str, "path=/")
	assert.Contains(t, str, "status=200")
}

func TestError(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := logrus.New()
	logger.SetOutput(buf)
	entry := logrus.NewEntry(logger)

	mux := route.NewServeMux()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := mux.NewContext(req, rec)
	h := func(c route.Context) error {
		return errors.New("ups")
	}
	mw := New(Entry(entry))

	// Status 5xx
	mw(c, h)

	str := buf.String()
	assert.Contains(t, str, "msg=ups")
	assert.Contains(t, str, "bytes_out=35")
	assert.Contains(t, str, "host=example.com")
	assert.Contains(t, str, "method=GET")
	assert.Contains(t, str, "path=/")
	assert.Contains(t, str, "status=500")
}
