package ginrouter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
	"gitlab.com/emixam23/fizz-and-buzz/tests/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func Test_ginRouter_GetMostUsedCombinationHandler(t *testing.T) {

	tests := []struct {
		name                     string
		serviceInteractionResult *models.FnbRequestInputStats
		serviceInteractionError  error
		wantResult               interface{}
		wantStatusCode           int
	}{
		{
			name: "Scenario 1 - OK",
			serviceInteractionResult: &models.FnbRequestInputStats{
				N1:    2,
				S1:    "two",
				N2:    3,
				S2:    "three",
				Limit: 4,
				Count: 1,
			},
			wantResult: &fnbRequestInputStats{
				N1:    2,
				S1:    "two",
				N2:    3,
				S2:    "three",
				Limit: 4,
				Count: 1,
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:                    "Scenario 2 - KO - service returns an error",
			serviceInteractionError: errors.New("any error"),
			wantResult:              errorStatus{Reason: "any error"},
			wantStatusCode:          http.StatusUnprocessableEntity,
		},
		{
			name:                     "Scenario 3 - KO - no entries retrieved, returns a 404 with empty body",
			serviceInteractionResult: nil,
			serviceInteractionError:  nil,
			wantResult:               struct{}{},
			wantStatusCode:           http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			statsServiceMock := mocks.NewMockIStatsService(ctrl)

			router := ginRouter{
				httpServer: &http.Server{
					Addr:    ":9000",
					Handler: gin.Default(),
				},
				servingURL:   ":9000",
				fnbService:   &mocks.MockIFnbService{},
				statsService: statsServiceMock,
			}

			ctxURL, err := url.Parse("/stats/most_used")
			if err != nil {
				t.Fatalf("an error '%s' was not expected when testing handler", err)
			}
			responseRecorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(responseRecorder)
			ctx.Request = &http.Request{
				Header: make(http.Header),
				Method: http.MethodGet,
				URL:    ctxURL,
			}

			statsServiceMock.EXPECT().GetFnbRequestsMostUsedCombination().Times(1).Return(tt.serviceInteractionResult, tt.serviceInteractionError)

			router.GetMostUsedCombinationHandler(ctx)
			if !reflect.DeepEqual(tt.wantStatusCode, responseRecorder.Code) {
				t.Errorf("GetMostUsedCombinationHandler() got = %v, want %v", responseRecorder.Code, tt.wantStatusCode)
				return
			}
			wantResultBytes, err := json.Marshal(tt.wantResult)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when testing handler", err)
			}
			responseBytes := responseRecorder.Body.Bytes()
			if !reflect.DeepEqual(wantResultBytes, responseBytes) {
				fmt.Println(string(responseBytes))
				fmt.Println(string(wantResultBytes))
				t.Errorf("GetMostUsedCombinationHandler() got = %s, want %s", string(responseBytes), string(wantResultBytes))
				return
			}
		})
	}
}

func Test_ginRouter_GetStatsHandler(t *testing.T) {

	type params struct {
		uriParams         *statsRequestParams
		uriSortedParamStr string
		expectServiceCall bool
	}
	tests := []struct {
		name                     string
		params                   params
		serviceInteractionResult []*models.FnbRequestInputStats
		serviceInteractionError  error
		wantResult               interface{}
		wantStatusCode           int
	}{
		{
			name: "Scenario 1 - OK - sorted",
			params: params{
				uriParams: &statsRequestParams{
					Sorted: true,
				},
				expectServiceCall: true,
			},
			serviceInteractionResult: []*models.FnbRequestInputStats{
				{
					N1:    2,
					S1:    "two",
					N2:    3,
					S2:    "three",
					Limit: 4,
					Count: 30,
				},
				{
					N1:    5,
					S1:    "five",
					N2:    6,
					S2:    "six",
					Limit: 7,
					Count: 25,
				},
			},
			wantResult: []*fnbRequestInputStats{
				{
					N1:    2,
					S1:    "two",
					N2:    3,
					S2:    "three",
					Limit: 4,
					Count: 30,
				},
				{
					N1:    5,
					S1:    "five",
					N2:    6,
					S2:    "six",
					Limit: 7,
					Count: 25,
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Scenario 2 - OK - not sorted",
			params: params{
				uriParams: &statsRequestParams{
					Sorted: false,
				},
				expectServiceCall: true,
			},
			serviceInteractionResult: []*models.FnbRequestInputStats{
				{
					N1:    5,
					S1:    "five",
					N2:    6,
					S2:    "six",
					Limit: 7,
					Count: 25,
				},
				{
					N1:    2,
					S1:    "two",
					N2:    3,
					S2:    "three",
					Limit: 4,
					Count: 30,
				},
			},
			wantResult: []*fnbRequestInputStats{
				{
					N1:    5,
					S1:    "five",
					N2:    6,
					S2:    "six",
					Limit: 7,
					Count: 25,
				},
				{
					N1:    2,
					S1:    "two",
					N2:    3,
					S2:    "three",
					Limit: 4,
					Count: 30,
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Scneario 3 - KO - validating params",
			params: params{
				uriParams:         nil,
				uriSortedParamStr: "test", // not a bool
				expectServiceCall: false,
			},
			wantResult:     errorStatus{Reason: "strconv.ParseBool: parsing \"test\": invalid syntax"},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Scenario 4 - KO - service returns an error",
			params: params{
				expectServiceCall: true,
			},
			serviceInteractionError: errors.New("any error"),
			wantResult:              errorStatus{Reason: "any error"},
			wantStatusCode:          http.StatusUnprocessableEntity,
		},
		{
			name: "Scenario 5 - KO - no entries retrieved, returns a 404 with empty body",
			params: params{
				expectServiceCall: true,
			},
			serviceInteractionResult: nil,
			serviceInteractionError:  nil,
			wantResult:               []struct{}{},
			wantStatusCode:           http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			statsServiceMock := mocks.NewMockIStatsService(ctrl)

			router := ginRouter{
				httpServer: &http.Server{
					Addr:    ":9000",
					Handler: gin.Default(),
				},
				servingURL:   ":9000",
				fnbService:   &mocks.MockIFnbService{},
				statsService: statsServiceMock,
			}

			var urlStr string
			if tt.params.uriSortedParamStr != "" {
				urlStr = fmt.Sprintf("/stats?sorted=%s", tt.params.uriSortedParamStr)
			} else if tt.params.uriParams != nil {
				urlStr = fmt.Sprintf("/stats?sorted=%v", tt.params.uriParams.Sorted)
			} else {
				urlStr = "/stats"
			}

			ctxURL, err := url.Parse(urlStr)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when testing handler", err)
			}
			responseRecorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(responseRecorder)
			ctx.Request = &http.Request{
				Header: make(http.Header),
				Method: http.MethodGet,
				URL:    ctxURL,
			}

			if tt.params.expectServiceCall {
				var sorted bool
				if tt.params.uriParams != nil {
					sorted = tt.params.uriParams.Sorted
				}
				statsServiceMock.EXPECT().GetFnbRequestsInputsStats(sorted).Times(1).Return(tt.serviceInteractionResult, tt.serviceInteractionError)
			}

			router.GetStatsHandler(ctx)
			if !reflect.DeepEqual(tt.wantStatusCode, responseRecorder.Code) {
				t.Errorf("GetStatsHandler() got = %v, want %v", responseRecorder.Code, tt.wantStatusCode)
				return
			}
			wantResultBytes, err := json.Marshal(tt.wantResult)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when testing handler", err)
			}
			responseBytes := responseRecorder.Body.Bytes()
			if !reflect.DeepEqual(wantResultBytes, responseBytes) {
				fmt.Println(string(responseBytes))
				fmt.Println(string(wantResultBytes))
				t.Errorf("GetStatsHandler() got = %s, want %s", string(responseBytes), string(wantResultBytes))
				return
			}
		})
	}
}
