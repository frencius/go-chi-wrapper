package server

import (
	"context"
	"fmt"
	net_http "net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListen(t *testing.T) {
	mockContext, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	mockServer := New(mockContext)
	mockPort := 8080

	t.Run("positive", func(t *testing.T) {
		expectedSrv := &net_http.Server{Addr: fmt.Sprintf(":%v", mockPort), Handler: mockServer.router}

		srv := mockServer.Listen(mockPort)
		defer func() {
			mockServer.Close(srv)
		}()

		assert.Equal(t, srv, expectedSrv)
	})

	t.Run("ListenAndServe Failed", func(t *testing.T) {
		mockServer.router = nil
		expectedSrv := &net_http.Server{Addr: fmt.Sprintf(":%v", mockPort), Handler: mockServer.router}

		srv := mockServer.Listen(mockPort)
		defer func() {
			mockServer.Close(srv)
		}()

		assert.Equal(t, srv, expectedSrv)
	})
}

func TestClose(t *testing.T) {
	mockContext, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	mockServer := New(mockContext)
	mockPort := 8080

	t.Run("positive", func(t *testing.T) {
		mockSrv := &net_http.Server{Addr: fmt.Sprintf(":%v", mockPort), Handler: mockServer.router}

		err := mockServer.Close(mockSrv)
		assert.Nil(t, err, "should be nil")
	})

	// t.Run("negative - close failed", func(t *testing.T) {
	// 	var convertInterface = func(server mock.IMockServer) *net_http.Server {
	// 		return server.(*net_http.Server)
	// 	}

	// 	mockSrv := &net_http.Server{Addr: fmt.Sprintf(":%v", mockPort), Handler: mockServer.router}
	// 	mockSrvNew := convertInterface(mockSrv)

	// 	mockCtrl := gomock.NewController(t)
	// 	defer mockCtrl.Finish()

	// 	mockHTTPServer := mock.NewMockHTTPServer(mockCtrl)
	// 	mockHTTPServer.EXPECT().Shutdown(mockSrvNew).Return(errors.New("an error"))

	// 	err := mockServer.Close(mockSrvNew)
	// 	assert.Error(t, err, "should be error")
	// })
}

func TestMethods(t *testing.T) {
	mockContext, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	mockServer := New(mockContext)
	mockPattern := "/"
	mockHandlerFunc := func(w net_http.ResponseWriter, r *net_http.Request) {}

	t.Run("positive", func(t *testing.T) {
		mockServer.Connect(mockPattern, mockHandlerFunc)
		mockServer.Delete(mockPattern, mockHandlerFunc)
		mockServer.Get(mockPattern, mockHandlerFunc)
		mockServer.Head(mockPattern, mockHandlerFunc)
		mockServer.Options(mockPattern, mockHandlerFunc)
		mockServer.Patch(mockPattern, mockHandlerFunc)
		mockServer.Post(mockPattern, mockHandlerFunc)
		mockServer.Put(mockPattern, mockHandlerFunc)
		mockServer.Trace(mockPattern, mockHandlerFunc)
	})
}

func TestMiddlewares(t *testing.T) {
	mockContext, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	mockServer := New(mockContext)

	t.Run("positive - WithValue() middleware", func(t *testing.T) {
		mockServer.WithValue("any", "any")
	})

	t.Run("positive - AllowCORS() middleware", func(t *testing.T) {
		mockServer.AllowCORS()
	})

	t.Run("positive - Swagger() middleware", func(t *testing.T) {
		mockServer.Swagger("/", "/")
	})
}
