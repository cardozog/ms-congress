package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const UsuarioIDContextKey = "usuarioID"
const PapeisContextKey = "papeis"

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.TrimSpace(secret) == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"erro": "segredo do token não configurado"})
			return
		}

		parts := strings.Fields(c.GetHeader("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "token Bearer obrigatório"})
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "token inválido"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "token inválido"})
			return
		}
		subject, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "token inválido"})
			return
		}
		userID, err := strconv.ParseUint(subject, 10, 64)
		if err != nil || userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "token inválido"})
			return
		}

		c.Set(UsuarioIDContextKey, uint(userID))
		if roles, ok := claims["papeis"].([]interface{}); ok {
			papeis := make([]string, 0, len(roles))
			for _, role := range roles {
				if value, ok := role.(string); ok {
					papeis = append(papeis, value)
				}
			}
			c.Set(PapeisContextKey, papeis)
		}
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get(PapeisContextKey)
		roles, ok := value.([]string)
		if !exists || !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"erro": "permissão insuficiente"})
			return
		}
		for _, current := range roles {
			if current == role {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"erro": "permissão insuficiente"})
	}
}

func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get(PapeisContextKey)
		currentRoles, ok := value.([]string)
		if !exists || !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"erro": "permissão insuficiente"})
			return
		}
		for _, current := range currentRoles {
			for _, allowed := range roles {
				if current == allowed {
					c.Next()
					return
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"erro": "permissão insuficiente"})
	}
}
