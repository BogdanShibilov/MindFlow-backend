package expertroutes

import (
	"errors"
	"strings"

	"github.com/bogdanshibilov/mindflowbackend/internal/entity"
	expertrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/expert"
	"github.com/gin-gonic/gin"
)

func getStatusQuery(ctx *gin.Context) (entity.Status, error) {
	statusString := ctx.Query("status")

	switch strings.ToLower(statusString) {
	case "pending":
		return entity.Pending, nil
	case "approved":
		return entity.Approved, nil
	case "rejected":
		return entity.Rejected, nil
	case "":
		return expertrepo.AllStatus, nil
	default:
		return -1, errors.New("invalid 'status' query parameter")
	}
}
