package serializers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang-ethereum-example-api/pkg/util"
	"net/http"
	"regexp"
)

var addressValidationRegex = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

type ErrorResponse struct {
	Message string
}

type Serializer struct {
	C *gin.Context
}

func (s *Serializer) ErrorResponse(e *util.ErrorInfo) {
	s.C.JSON(e.HttpCode, ErrorResponse{
		Message: e.Message,
	})
}

func validate(ctx context.Context, form interface{}) *util.ErrorInfo {
	validate := validator.New()
	err := validate.StructCtx(ctx, form)
	if err != nil {
		return &util.ErrorInfo{
			HttpCode: http.StatusBadRequest,
			Err:      validate.StructCtx(ctx, form),
			Message:  util.ValidationErrorMessage,
		}
	}
	return nil
}

func (s *Serializer) ShouldBindUri(obj interface{}) *util.ErrorInfo {
	err := s.C.ShouldBindUri(obj)
	if err != nil {
		return &util.ErrorInfo{
			HttpCode: http.StatusInternalServerError,
			Err:      err,
			Message:  util.BindingErrorMessage,
		}
	}
	return nil
}
func (s *Serializer) ShouldBindJSON(obj interface{}) *util.ErrorInfo {
	err := s.C.ShouldBindJSON(obj)
	if err != nil {
		return &util.ErrorInfo{
			HttpCode: http.StatusInternalServerError,
			Err:      err,
			Message:  util.BindingErrorMessage,
		}
	}
	return nil
}

func (s *Serializer) SuccessfulResponse(httpCode int, data interface{}) {
	s.C.JSON(httpCode, data)
}
