package middlewares

import (
	"net/http"
	"github.com/Caknoooo/go-gin-clean-starter/modules/user/dto"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthorizeRole(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_DENIED_ACCESS, nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		userRole := role.(string)
		isAllowed := false
		for _, r := range allowedRoles {
			if userRole == r {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, "Forbidden: you do not have the required role", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}
		ctx.Next()
	}
}
