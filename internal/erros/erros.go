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

var ErroCredenciaisObrigatorias = errors.New("email e senha são obrigatórios")

var ErroCredenciaisInvalidas = errors.New("email ou senha inválidos")

var ErroConfiguracaoToken = errors.New("segredo do token não configurado")

var ErroTokenInvalido = errors.New("token inválido")

var ErroUsuarioNaoAutorizado = errors.New("usuário não autorizado")
