package ginrouter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models/errormodels"
	"gitlab.com/emixam23/fizz-and-buzz/tests/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func Test_ginRouter_GetHistoryHandler(t *testing.T) {

	now := time.Now()
	limit2 := uint64(2)

	type params struct {
		uriParams         historyRequestParams
		uriLimitParamStr  string
		expectServiceCall bool
	}
	tests := []struct {
		name                     string
		params                   params
		serviceInteractionResult []*models.FnbRequest
		serviceInteractionError  error
		wantResult               interface{}
		wantStatusCode           int
	}{
		{
			name: "Scenario 1 - OK - limit (2) provided",
			params: params{
				uriParams: historyRequestParams{
					Limit: &limit2,
				},
				uriLimitParamStr:  "",
				expectServiceCall: true,
			},
			serviceInteractionResult: []*models.FnbRequest{
				{
					ID:          1,
					RequestDate: &now,
					N1:          2,
					S1:          "two",
					N2:          3,
					S2:          "three",
					Limit:       4,
				},
				{
					ID:          5,
					RequestDate: &now,
					N1:          6,
					S1:          "six",
					N2:          7,
					S2:          "seven",
					Limit:       8,
				},
			},
			wantResult: []*fnbRequest{
				{
					ID:          1,
					RequestDate: &now,
					N1:          2,
					S1:          "two",
					N2:          3,
					S2:          "three",
					Limit:       4,
				},
				{
					ID:          5,
					RequestDate: &now,
					N1:          6,
					S1:          "six",
					N2:          7,
					S2:          "seven",
					Limit:       8,
				},
			},
			wantStatusCode: http.StatusOK,
		}, {
			name: "Scenario 2 - OK - no limit provided",
			params: params{
				uriParams: historyRequestParams{
					Limit: nil,
				},
				uriLimitParamStr:  "",
				expectServiceCall: true,
			},
			serviceInteractionResult: []*models.FnbRequest{
				{
					ID:          1,
					RequestDate: &now,
					N1:          2,
					S1:          "two",
					N2:          3,
					S2:          "three",
					Limit:       4,
				},
				{
					ID:          5,
					RequestDate: &now,
					N1:          6,
					S1:          "six",
					N2:          7,
					S2:          "seven",
					Limit:       8,
				},
				{
					ID:          9,
					RequestDate: &now,
					N1:          10,
					S1:          "ten",
					N2:          11,
					S2:          "eleven",
					Limit:       12,
				},
			},
			wantResult: []*fnbRequest{
				{
					ID:          1,
					RequestDate: &now,
					N1:          2,
					S1:          "two",
					N2:          3,
					S2:          "three",
					Limit:       4,
				},
				{
					ID:          5,
					RequestDate: &now,
					N1:          6,
					S1:          "six",
					N2:          7,
					S2:          "seven",
					Limit:       8,
				},
				{
					ID:          9,
					RequestDate: &now,
					N1:          10,
					S1:          "ten",
					N2:          11,
					S2:          "eleven",
					Limit:       12,
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Scenario 3 - KO - validating params",
			params: params{
				uriParams: historyRequestParams{
					Limit: nil,
				},
				uriLimitParamStr:  "test", // invalid param as uint64 is expected if provided
				expectServiceCall: false,
			},
			wantResult:     errorStatus{Reason: "strconv.ParseUint: parsing \"test\": invalid syntax"},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Scenario 4 - KO - service returns an error",
			params: params{
				uriParams: historyRequestParams{
					Limit: &limit2,
				},
				uriLimitParamStr:  "",
				expectServiceCall: true,
			},
			serviceInteractionError: errormodels.NewUnprocessableError(errors.New("any error")),
			wantResult:              errorStatus{Reason: "any error"},
			wantStatusCode:          http.StatusUnprocessableEntity,
		},
		{
			name: "Scenario 5 - KO - no entries retrieved, returns a 404 with empty array",
			params: params{
				uriParams: historyRequestParams{
					Limit: &limit2,
				},
				uriLimitParamStr:  "",
				expectServiceCall: true,
			},
			serviceInteractionResult: []*models.FnbRequest{},
			wantResult:               []struct{}{},
			wantStatusCode:           http.StatusNotFound,
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

			var urlStr string
			if tt.params.uriLimitParamStr != "" {
				urlStr = fmt.Sprintf("/history?limit=%s", tt.params.uriLimitParamStr)
			} else if tt.params.uriParams.Limit != nil {
				urlStr = fmt.Sprintf("/history?limit=%d", tt.params.uriParams.Limit)
			} else {
				urlStr = "/history"
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
				fnbServiceMock.EXPECT().GetFnbRequestsHistory(gomock.Any()).Times(1).Return(tt.serviceInteractionResult, tt.serviceInteractionError)
			}

			router.GetHistoryHandler(ctx)
			if !reflect.DeepEqual(tt.wantStatusCode, responseRecorder.Code) {
				t.Errorf("GetHistoryHandler() got = %v, want %v", responseRecorder.Code, tt.wantStatusCode)
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
				t.Errorf("GetHistoryHandler() got = %s, want %s", string(responseBytes), string(wantResultBytes))
				return
			}
		})
	}
}
