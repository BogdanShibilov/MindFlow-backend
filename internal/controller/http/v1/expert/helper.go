package expertroutes

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/bogdanshibilov/mindflowbackend/internal/entity"
)

func getStatusQuery(ctx *gin.Context) entity.Status {
	statusString := ctx.Query("status")

	switch strings.ToLower(statusString) {
	case "pending", "0":
		return entity.Pending
	case "approved", "1":
		return entity.Approved
	case "rejected", "2":
		return entity.Rejected
	default:
		return -1
	}
}
