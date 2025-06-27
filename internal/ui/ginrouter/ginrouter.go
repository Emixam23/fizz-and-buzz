package ginrouter

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/services/fnbservice"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/services/statsservice"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/ui"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/utils/slice"
	"net/http"
)

type ginRouter struct {
	servingURL string
	httpServer *http.Server

	fnbService   fnbservice.IFnbService
	statsService statsservice.IStatsService
}

const (
	healthCheckEndpoint = "/health"

	fizzAndBuzzEndpoint              = "/fizz-and-buzz"
	historyEndpoint                  = "/history"
	statsEndpoint                    = "/stats"
	statsMostUsedCombinationEndpoint = "/stats/most_used"
)

func ginEngineModes() []string {
	return []string{gin.DebugMode, gin.TestMode, gin.ReleaseMode}
}

// New creates and initializes a new router (gin api)
func New(config *ui.Config, fnbService fnbservice.IFnbService, statsService statsservice.IStatsService) (ui.IRestAPI, error) {

	if config == nil {
		return nil, errors.New("no config provided")
	} else if fnbService == nil {
		return nil, errors.New("provided fnb service is not initialized (nil)")
	} else if statsService == nil {
		return nil, errors.New("provided stats service is not initialized (nil)")
	} else if config.Port == 0 {
		return nil, errors.New("no port config provided")
	} else if config.Mode == "" {
		return nil, errors.New("no router mode provided")
	} else if !slice.IsStringInSlice(config.Mode, ginEngineModes()) {
		return nil, errors.New("invalid router mode provided")
	}

	ginEngine := gin.New()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	ginEngine.Use(cors.New(corsConfig))
	ginEngine.Use(gin.Recovery())
	gin.SetMode(config.Mode)

	servingURL := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// initialize struct
	router := &ginRouter{
		servingURL: servingURL,
		httpServer: &http.Server{
			Handler: ginEngine,
			Addr:    servingURL,
		},
		fnbService:   fnbService,
		statsService: statsService,
	}

	ginEngine.GET("/", router.HealthCheckHandler)
	ginEngine.GET(healthCheckEndpoint, router.HealthCheckHandler)

	ginEngine.GET(fizzAndBuzzEndpoint, router.FizzAndBuzzHandler)
	ginEngine.GET(historyEndpoint, router.GetHistoryHandler)
	ginEngine.GET(statsEndpoint, router.GetStatsHandler)
	ginEngine.GET(statsMostUsedCombinationEndpoint, router.GetMostUsedCombinationHandler)

	return router, nil
}

// ListenAndServe trigger the process for the router to start listen for new incoming requests
func (r *ginRouter) ListenAndServe() error {
	if err := r.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("an error happened while listening and serving")
	}

	return nil
}

// Shutdown gracefully shutdown the router
func (r *ginRouter) Shutdown() error {
	return r.httpServer.Shutdown(context.Background())
}

// ServeHTTP is a pass through to the initial handler ServeHTTP
// This function is implemented for integration testing purpose
func (r *ginRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.httpServer.Handler.ServeHTTP(w, req)
}

// HealthCheckHandler
// @Description Return a string message meaning the service is running (and so healthy).
// @ID health
// @Accept json
// @Produce json
// @Success 200 {object} response.Status
// @Router /health [get]
func (r *ginRouter) HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, healthStatus{Status: "Ok"})
}
