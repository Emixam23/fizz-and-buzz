package ginrouter

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

// GetHistoryHandler
// @Description returns history of fizz and buzz users requests.
// @ID get-history
// @Accept json
// @Produce json
// @Success 200 {array} fnbRequest
// @Failure 400 {object} errorStatus
// @Failure 404 {array} struct
// @Failure 422 {object} errorStatus
// @Router /history [get]
func (r *ginRouter) GetHistoryHandler(context *gin.Context) {

	var uriParams historyRequestParams
	if err := context.BindQuery(&uriParams); err != nil {
		log.Error().Err(err).Interface("uri_params", uriParams).Msg("Invalid parameter(s) provided")
		context.AbortWithStatusJSON(http.StatusBadRequest, newErrorStatus(err))
		return
	}

	fnbRequestsHistory, err := r.fnbService.GetFnbRequestsHistory(uriParams.Limit)
	if err != nil {
		log.Error().Err(err).Msg("Couldn't retrieve fizz and buzz requests history")
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, newErrorStatus(err))
		return
	} else if len(fnbRequestsHistory) == 0 {
		context.AbortWithStatusJSON(http.StatusNotFound, []struct{}{})
		return
	}

	context.JSON(http.StatusOK, fromDomainFnbRequestsToUIFnbRequests(fnbRequestsHistory))
	return
}
