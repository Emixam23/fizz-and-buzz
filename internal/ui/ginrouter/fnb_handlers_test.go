package ginrouter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models/errormodels"
	"gitlab.com/emixam23/fizz-and-buzz/tests/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func Test_ginRouter_FizzAndBuzzHandler(t *testing.T) {

	type params struct {
		uriParams         fizzAndBuzzRequestParams
		expectServiceCall bool
	}
	tests := []struct {
		name                     string
		params                   params
		serviceInteractionResult []string
		serviceInteractionError  error
		wantResult               interface{}
		wantStatusCode           int
	}{
		{
			name: "Scneario 1 - OK",
			params: params{
				uriParams: fizzAndBuzzRequestParams{
					N1:    3,
					S1:    "fizz",
					N2:    5,
					S2:    "buzz",
					Limit: 19,
				},
				expectServiceCall: true,
			},
			serviceInteractionResult: []string{
				"1",
				"2",
				"fizz",
				"4",
				"buzz",
				"fizz",
				"7",
				"8",
				"fizz",
				"buzz",
				"11",
				"fizz",
				"13",
				"14",
				"fizzbuzz",
				"16",
				"17",
				"fizz",
				"19",
			},
			wantResult: fizzAndBuzzResponse{
				Request: fizzAndBuzzRequestParams{
					N1:    3,
					S1:    "fizz",
					N2:    5,
					S2:    "buzz",
					Limit: 19,
				},
				Result: []string{
					"1",
					"2",
					"fizz",
					"4",
					"buzz",
					"fizz",
					"7",
					"8",
					"fizz",
					"buzz",
					"11",
					"fizz",
					"13",
					"14",
					"fizzbuzz",
					"16",
					"17",
					"fizz",
					"19",
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Scneario 2 - KO - validating params",
			params: params{
				uriParams: fizzAndBuzzRequestParams{
					N1:    0, // will appear as not present -> but it's required  -> error will be thrown
					S1:    "fizz",
					N2:    5,
					S2:    "buzz",
					Limit: 19,
				},
				expectServiceCall: false,
			},
			wantResult:     errorStatus{Reason: "Key: 'fizzAndBuzzRequestParams.N1' Error:Field validation for 'N1' failed on the 'required' tag"},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Scneario 3 - KO - validating 0 and limit",
			params: params{
				uriParams: fizzAndBuzzRequestParams{
					N1:    3,
					S1:    "fizz",
					N2:    5,
					S2:    "buzz",
					Limit: 19,
				},
				expectServiceCall: true,
			},
			serviceInteractionError: errormodels.NewUnprocessableError(errors.New("any error")),
			wantResult:              errorStatus{Reason: "any error"},
			wantStatusCode:          http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			fnbServiceMock := mocks.NewMockIFnbService(ctrl)

			router := ginRouter{
				httpServer: &http.Server{
					Addr:    ":9000",
					Handler: gin.Default(),
				},
				servingURL:   ":9000",
				fnbService:   fnbServiceMock,
				statsService: &mocks.MockIStatsService{},
			}

			ctxURL, err := url.Parse(fmt.Sprintf("/fizz-and-buzz?n1=%d&s1=%s&n2=%d&s2=%s&limit=%d", tt.params.uriParams.N1, tt.params.uriParams.S1, tt.params.uriParams.N2, tt.params.uriParams.S2, tt.params.uriParams.Limit))
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
				fnbServiceMock.EXPECT().GetFizzAndBuzz(tt.params.uriParams.N1, tt.params.uriParams.S1, tt.params.uriParams.N2, tt.params.uriParams.S2, tt.params.uriParams.Limit).Times(1).Return(tt.serviceInteractionResult, tt.serviceInteractionError)
			}

			router.FizzAndBuzzHandler(ctx)
			if !reflect.DeepEqual(tt.wantStatusCode, responseRecorder.Code) {
				t.Errorf("FizzAndBuzzHandler() got = %v, want %v", responseRecorder.Code, tt.wantStatusCode)
				return
			}
			wantResultBytes, err := json.Marshal(tt.wantResult)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when testing handler", err)
			}
			responseBytes := responseRecorder.Body.Bytes()
			if !reflect.DeepEqual(wantResultBytes, responseBytes) {
				t.Errorf("FizzAndBuzzHandler() got = %s, want %s", string(responseBytes), string(wantResultBytes))
				return
			}
		})
	}
}
