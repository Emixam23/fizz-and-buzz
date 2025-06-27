package ginrouter

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

// GetStatsHandler
// @Description returns the statistics of the keys/strings usage based on fizz and buzz users requests.
// @ID get-stats
// @Accept json
// @Produce json
// @Success 200 {array} fnbRequestInputStats
// @Failure 400 {object} errorStatus
// @Failure 404 {array} struct
// @Failure 422 {object} errorStatus
// @Router /stats/{sorted} [get]
func (r *ginRouter) GetStatsHandler(context *gin.Context) {

	var uriParams statsRequestParams
	if err := context.BindQuery(&uriParams); err != nil {
		log.Error().Err(err).Interface("uri_params", uriParams).Msg("Invalid parameter(s) provided")
		context.AbortWithStatusJSON(http.StatusBadRequest, newErrorStatus(err))
		return
	}

	result, err := r.statsService.GetFnbRequestsInputsStats(uriParams.Sorted)
	if err != nil {
		log.Error().Err(err).Msg("Couldn't retrieve stats")
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, newErrorStatus(err))
		return
	} else if len(result) == 0 {
		context.AbortWithStatusJSON(http.StatusNotFound, []struct{}{})
		return
	}

	context.JSON(http.StatusOK, fromDomainFnbRequestsInputsStatsToUIFnbRequestsInputsStats(result))
	return
}

// GetMostUsedCombinationHandler
// @Description returns the most used combination based on fizz and buzz users requests.
// @ID get-most-used-combination
// @Accept json
// @Produce json
// @Success 200 {object} fnbRequestInputStats
// @Failure 404 {array} struct
// @Failure 422 {object} errorStatus
// @Router /stats/most_used [get]
func (r *ginRouter) GetMostUsedCombinationHandler(context *gin.Context) {

	result, err := r.statsService.GetFnbRequestsMostUsedCombination()
	if err != nil {
		log.Error().Err(err).Msg("Couldn't retrieve stats")
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, newErrorStatus(err))
		return
	} else if result == nil {
		context.AbortWithStatusJSON(http.StatusNotFound, struct{}{})
		return
	}

	context.JSON(http.StatusOK, fromDomainFnbRequestInputStatsToUIFnbRequestInputStats(result))
	return
}
