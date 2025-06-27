package ginrouter

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

// FizzAndBuzzHandler
// @Description fizz and ... buzz! The below controller returns the results for the provided request.
// @ID fizz-and-buzz
// @Produce json
// @Success 200 {array} string
// @Failure 400 {object} errorStatus
// @Failure 422 {object} errorStatus
// @Failure 500 {object} errorStatus
// @Router /fizz-and-buzz [get]
func (r *ginRouter) FizzAndBuzzHandler(context *gin.Context) {

	var uriParams fizzAndBuzzRequestParams
	if err := context.BindQuery(&uriParams); err != nil {
		log.Error().Err(err).Interface("uri_params", uriParams).Msg("Invalid parameter(s) provided")
		context.AbortWithStatusJSON(http.StatusBadRequest, newErrorStatus(err))
		return
	}

	result, err := r.fnbService.GetFizzAndBuzz(uriParams.N1, uriParams.S1, uriParams.N2, uriParams.S2, uriParams.Limit)
	if err != nil {
		log.Error().Err(err).Interface("uri_params", uriParams).Msg("Couldn't get fizz and buzz computation results")
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, newErrorStatus(err))
		return
	}

	context.JSON(http.StatusOK, fizzAndBuzzResponse{
		Request: uriParams,
		Result:  result,
	})
	return
}
