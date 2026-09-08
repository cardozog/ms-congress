package erros

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	StatusCode int   `json:"status_code"`
	Error      error `json:"message"`
}

func HandleErrorWithStatus(c *gin.Context, apiErr *ApiError) {
	c.AbortWithError(apiErr.StatusCode, apiErr.Error)

}

var ErroEmailJaCadastrado = errors.New("email já cadastrado")

var ErroEventoNaoEncontrado = errors.New("evento não encontrado")
