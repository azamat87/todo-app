package handler

import (
	"bytes"
	"errors"
	todoapp "golang_ninja/todo-app"
	"golang_ninja/todo-app/pkg/service"
	mock_service "golang_ninja/todo-app/pkg/service/mocks"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavier func(s *mock_service.MockAuthorization, user todoapp.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           todoapp.User
		mockBehavier        mockBehavier
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name": "Test", "username": "test", "password":"qwerty"}`,
			inputUser: todoapp.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavier: func(s *mock_service.MockAuthorization, user todoapp.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"name": "Test", "password":"qwerty"}`,
			mockBehavier:        func(s *mock_service.MockAuthorization, user todoapp.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"name": "Test", "username": "test", "password":"qwerty"}`,
			inputUser: todoapp.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavier: func(s *mock_service.MockAuthorization, user todoapp.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavier(auth, testCase.inputUser)

			service := &service.Service{Authorization: auth}
			handler := NewHandler(service)

			// Test server
			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(testCase.inputBody))

			// request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, req.Body)
		})
	}
}
