package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"realty-avito/internal/http-server/middleware"
)

func TestJWTMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		authorization    string
		expectedUserType string
		expectError      bool
	}{
		{
			name:             "Valid client token",
			authorization:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX3R5cGUiOiJjbGllbnQiLCJqdGkiOiIxNzIzOTk5ODA1NDE2NzM5MDAwIn0.f9bXxK0gmrvNulNe8SuiNRw5xBsSj6gRLklKdEZ0Z2o",
			expectedUserType: "client",
			expectError:      false,
		},
		{
			name:             "Valid moderator token",
			authorization:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX3R5cGUiOiJtb2RlcmF0b3IiLCJqdGkiOiIxNzIzOTk5ODc2NzA0MzQyMDAwIn0._GXniBxxwlL6tDno3XnS4JkaHH7IsFsmwyj86iIVr2o",
			expectedUserType: "moderator",
			expectError:      false,
		},
		{
			name:          "Invalid token",
			authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ7777777R5cGUiOiJtb2RlcmF0b3IiLCJqdGkiOiIxNzIzOTk5ODc2NzA0MzQyMDAwIn0._GXniBxxwlL6tDno3XnS4JkaHH7IsFsmwyj86iIVr2o",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", tt.authorization)

			rr := httptest.NewRecorder()

			handler := middleware.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userType := r.Context().Value("user_type").(string)
				require.Equal(t, tt.expectedUserType, userType)
			}))

			handler.ServeHTTP(rr, req)

			if tt.expectError {
				require.Equal(t, http.StatusUnauthorized, rr.Code)
			} else {
				require.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}

func TestJWTModeratorOnlyMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		authorization    string
		expectedUserType string
		expectError      bool
	}{
		{
			name:             "Valid moderator token",
			authorization:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX3R5cGUiOiJtb2RlcmF0b3IiLCJqdGkiOiIxNzIzOTk5ODc2NzA0MzQyMDAwIn0._GXniBxxwlL6tDno3XnS4JkaHH7IsFsmwyj86iIVr2o",
			expectedUserType: "moderator",
			expectError:      false,
		},
		{
			name:             "Valid client token",
			authorization:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX3R5cGUiOiJjbGllbnQiLCJqdGkiOiIxNzIzOTk5ODA1NDE2NzM5MDAwIn0.f9bXxK0gmrvNulNe8SuiNRw5xBsSj6gRLklKdEZ0Z2o",
			expectedUserType: "client",
			expectError:      true,
		},
		{
			name:          "Invalid token",
			authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ7777777R5cGUiOiJtb2RlcmF0b3IiLCJqdGkiOiIxNzIzOTk5ODc2NzA0MzQyMDAwIn0._GXniBxxwlL6tDno3XnS4JkaHH7IsFsmwyj86iIVr2o",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", tt.authorization)

			rr := httptest.NewRecorder()

			handler := middleware.JWTModeratorOnlyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.expectedUserType, r.Context().Value("user_type").(string))
			}))

			handler.ServeHTTP(rr, req)

			if tt.expectError {
				require.Equal(t, http.StatusUnauthorized, rr.Code)
			} else {
				require.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}
