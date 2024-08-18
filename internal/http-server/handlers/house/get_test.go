package house

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"realty-avito/internal/http-server/handlers"
	"realty-avito/internal/lib/logger"
	"realty-avito/internal/repositories/flat"
	"realty-avito/internal/repositories/flat/mocks"
)

func TestGetFlatsInHouseHandler(t *testing.T) {
	mockFlatsRepo := new(mocks.FlatsRepository)
	log := logger.SetupLogger("local")

	r := chi.NewRouter()
	handler := GetFlatsInHouseHandler(log, mockFlatsRepo) // ваш хэндлер

	r.Get("/house/{id}", handler)

	tests := []struct {
		name               string
		userType           string
		houseID            string
		prepareMock        func()
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:     "valid client with approved flats",
			userType: "client",
			houseID:  "1",
			prepareMock: func() {
				mockFlatsRepo.On("GetApprovedFlatsByHouseID", mock.Anything, int64(1)).Return([]flat.FlatEntity{
					{ID: 1, HouseID: 1, Status: "approved"},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: Response{
				Flats: []handlers.Flat{
					{ID: 1, HouseID: 1, Status: "approved"},
				},
			},
		},
		{
			name:     "valid moderator with all flats",
			userType: "moderator",
			houseID:  "1",
			prepareMock: func() {
				mockFlatsRepo.On("GetFlatsByHouseID", mock.Anything, int64(1)).Return([]flat.FlatEntity{
					{ID: 1, HouseID: 1, Status: "created"},
					{ID: 2, HouseID: 1, Status: "approved"},
					{ID: 3, HouseID: 1, Status: "declined"},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: Response{
				Flats: []handlers.Flat{
					{ID: 1, HouseID: 1, Status: "created"},
					{ID: 2, HouseID: 1, Status: "approved"},
					{ID: 3, HouseID: 1, Status: "declined"},
				},
			},
		},
		{
			name:               "missing user type",
			userType:           "",
			houseID:            "1",
			prepareMock:        func() {},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   nil,
		},
		{
			name:               "unauthorized user type",
			userType:           "invalid",
			houseID:            "1",
			prepareMock:        func() {},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMock()

			req, err := http.NewRequest("GET", fmt.Sprintf("/house/%s", tt.houseID), nil)
			require.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.houseID)

			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			ctx = context.WithValue(req.Context(), "user_type", tt.userType)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)

			// Если ожидаемый ответ не пустой, проверяем его
			if tt.expectedResponse != nil {
				var actualResponse interface{}
				if rr.Code == http.StatusOK {
					var resp Response
					require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
					actualResponse = resp
				} else {
					actualResponse = rr.Body.String()
				}
				require.Equal(t, tt.expectedResponse, actualResponse)
			}
		})
	}
}
