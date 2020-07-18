package test

import (
	"github.com/fabric-app/routers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var router *gin.Engine

func init() {
	gin.SetMode("test")
	router = routers.InitRouter()
}

func TestTestGetRouter(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"test\"}\n", w.Body.String())
}
