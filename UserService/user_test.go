package main

import (
	controller "deneme.com/bng-go/Controller"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterUserRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	uc := &controller.UserController{}

	rg := r.Group("/")
	uc.RegisterUserRoutes(rg)
	for i := 0; i < 20; i++ {

		req, err := http.NewRequest(http.MethodGet, "/user", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	}
}
