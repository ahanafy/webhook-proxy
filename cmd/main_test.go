package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"gotest.tools/v3/assert"
	"k8s.io/client-go/kubernetes"
)

type Suite struct {
	suite.Suite
	wh        *WebHandler
	Context   context.Context
	ClientSet *kubernetes.Clientset
}

func (suite *Suite) SetupTest() {
	suite.wh = &WebHandler{
		context:   suite.Context,
		clientSet: &kubernetes.Clientset{},
	}
}
func (suite *Suite) TestWebhookHandlerNoBody() {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/sleuth", nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.wh.WebhookHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status is expected
	assert.Equal(suite.T(), rr.Code, http.StatusBadRequest, "handler returned unexpected status")

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
