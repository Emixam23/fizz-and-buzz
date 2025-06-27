package ginrouter

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/services/fnbservice"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/services/statsservice"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/ui"
	"gitlab.com/emixam23/fizz-and-buzz/tests/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestModes(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "Scenario 1 - OK",
			want: []string{gin.DebugMode, gin.TestMode, gin.ReleaseMode},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ginEngineModes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ginEngineModes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {

	ctrl := gomock.NewController(t)
	fnbService := mocks.NewMockIFnbService(ctrl)
	statsService := mocks.NewMockIStatsService(ctrl)

	type params struct {
		config       *ui.Config
		fnbService   fnbservice.IFnbService
		statsService statsservice.IStatsService
	}
	tests := []struct {
		name                     string
		params                   params
		wantGinRouterInitialized bool
		wantErr                  error
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				config: &ui.Config{
					Host: "0.0.0.0",
					Port: 9000,
					Mode: ginEngineModes()[1], // test
				},
				fnbService:   fnbService,
				statsService: statsService,
			},
			wantGinRouterInitialized: true,
		},
		{
			name: "Scenario 2 - KO - config is nil",
			params: params{
				config:       nil,
				fnbService:   fnbService,
				statsService: statsService,
			},
			wantErr: errors.New("no config provided"),
		},
		{
			name: "Scenario 3 - KO - fnbService is nil",
			params: params{
				config: &ui.Config{
					Host: "0.0.0.0",
					Port: 9000,
					Mode: ginEngineModes()[1], // test
				},
				fnbService:   nil,
				statsService: statsService,
			},
			wantErr: errors.New("provided fnb service is not initialized (nil)"),
		},
		{
			name: "Scenario 4 - KO - statsService is nil",
			params: params{
				config: &ui.Config{
					Host: "0.0.0.0",
					Port: 9000,
					Mode: ginEngineModes()[1], // test
				},
				fnbService:   fnbService,
				statsService: nil,
			},
			wantErr: errors.New("provided stats service is not initialized (nil)"),
		},
		{
			name: "Scenario 5 - KO - config port not defined",
			params: params{
				config: &ui.Config{
					Host: "0.0.0.0",
					Port: 0,
					Mode: ginEngineModes()[1], // test
				},
				fnbService:   fnbService,
				statsService: statsService,
			},
			wantErr: errors.New("no port config provided"),
		},
		{
			name: "Scenario 6 - KO - config gin router mod not defined",
			params: params{
				config: &ui.Config{
					Host: "0.0.0.0",
					Port: 9000,
					Mode: "",
				},
				fnbService:   fnbService,
				statsService: statsService,
			},
			wantErr: errors.New("no router mode provided"),
		},
		{
			name: "Scenario 7 - KO - config gin router mod is invalid",
			params: params{
				config: &ui.Config{
					Host: "0.0.0.0",
					Port: 9000,
					Mode: "azerty",
				},
				fnbService:   fnbService,
				statsService: statsService,
			},
			wantErr: errors.New("invalid router mode provided"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.params.config, tt.params.fnbService, tt.params.statsService)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantGinRouterInitialized {
				assert.NotNil(t, got, "router should be initialized")
				assert.NotNil(t, got.(*ginRouter).fnbService, "fnb service should be initialized")
				assert.NotNil(t, got.(*ginRouter).statsService, "stats service should be initialized")
				assert.Equal(t, "0.0.0.0:9000", got.(*ginRouter).servingURL, "serving URL should be initialized")
				assert.NotNil(t, got.(*ginRouter).httpServer, "http server should be initialized")
				assert.NotNil(t, got.(*ginRouter).httpServer.Handler, "gin engine should be initialized")
				assert.NotNil(t, got.(*ginRouter).httpServer.Handler.(*gin.Engine).Routes(), "gin engine should be initialized")
				assert.Equal(t, len(got.(*ginRouter).httpServer.Handler.(*gin.Engine).Routes()), 6, "gin engine should have 6 defined routes")
			}
		})
	}
}

func Test_ginRouter_HealthCheckHandler(t *testing.T) {

	tests := []struct {
		name           string
		healthURL      string
		wantResult     healthStatus
		wantStatusCode int
	}{
		{
			name:           "Scenario 1 - OK",
			healthURL:      "",
			wantResult:     healthStatus{Status: "Ok"},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "Scenario 2 - OK",
			healthURL:      "/health",
			wantResult:     healthStatus{Status: "Ok"},
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := ginRouter{
				httpServer: &http.Server{
					Addr:    ":9000",
					Handler: gin.Default(),
				},
				servingURL:   ":9000",
				fnbService:   &mocks.MockIFnbService{},
				statsService: &mocks.MockIStatsService{},
			}

			ctxURL, err := url.Parse(tt.healthURL)
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

			router.HealthCheckHandler(ctx)
			if !reflect.DeepEqual(responseRecorder.Code, tt.wantStatusCode) {
				t.Errorf("HealthCheckHandler() got = %v, want %v", responseRecorder.Code, tt.wantStatusCode)
				return
			}
			wantResultBytes, err := json.Marshal(tt.wantResult)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when testing handler", err)
			}
			responseBytes := responseRecorder.Body.Bytes()
			if !reflect.DeepEqual(responseBytes, wantResultBytes) {
				t.Errorf("HealthCheckHandler() got = %s, want %s", string(responseBytes), string(wantResultBytes))
				return
			}
		})
	}
}

func Test_ginRouter_ListenAndServe_Then_Shutdown(t *testing.T) {

	router := ginRouter{
		servingURL:   ":9000",
		httpServer:   &http.Server{Addr: ":9000", Handler: gin.Default()},
		fnbService:   &mocks.MockIFnbService{},
		statsService: &mocks.MockIStatsService{},
	}

	var listenAndServeError, shutdownError error
	routerRunning := make(chan struct{})
	routerDone := make(chan struct{})
	go func() {
		close(routerRunning)
		listenAndServeError = router.ListenAndServe()
		defer close(routerDone)
	}()

	<-routerRunning
	shutdownError = router.Shutdown()
	<-routerDone

	assert.Nil(t, listenAndServeError, "no error should be thrown if shutdown")
	assert.Nil(t, shutdownError, "no error should be thrown, it should be working")
}
