package web_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/web"
)

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/version/", gin.WrapF(web.VersionHandler()))
	return r
}

func Test_VersionHandler(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	r := newTestRouter()
	a.NotNil(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/version/", nil)
	r.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code)

	var version map[string]string
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&version)
	a.NoError(err)
	a.Equal("development", version["build_version"])
}
