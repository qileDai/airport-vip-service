package routes

import (
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func setupAuditRoutes(g *echo.Group, svc *services.AuditService) {
	r := g.Group("/audit-logs")

	r.GET("", listAuditLogs(svc))
	r.GET("/:id", getAuditLog(svc))
	r.GET("/entity/:type/:id", getAuditLogsByEntity(svc))
}

func listAuditLogs(svc *services.AuditService) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &schemas.AuditLogQueryRequest{}

		req.Page, _ = strconv.Atoi(c.QueryParam("page"))
		if req.Page <= 0 {
			req.Page = 1
		}
		req.PerPage, _ = strconv.Atoi(c.QueryParam("per_page"))
		if req.PerPage <= 0 {
			req.PerPage = 20
		}

		req.EntityType = c.QueryParam("entity_type")
		req.OperationType = c.QueryParam("operation_type")
		req.Operator = c.QueryParam("operator")

		result, err := svc.List(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func getAuditLog(svc *services.AuditService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid ID",
			})
		}

		result, err := svc.GetByID(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		if result == nil {
			return c.JSON(http.StatusNotFound, schemas.ErrorResponse{
				Code:    404,
				Message: "Audit log not found",
			})
		}

		return c.JSON(http.StatusOK, schemas.AuditLogResponse{
			ID:                result.ID,
			LogNo:             result.LogNo,
			OperationType:     result.OperationType,
			EntityType:        result.EntityType,
			EntityID:          result.EntityID,
			EntityNo:          result.EntityNo,
			OldValue:          result.OldValue,
			NewValue:          result.NewValue,
			Operator:          result.Operator,
			OperatorRole:      result.OperatorRole,
			IPAddress:         result.IPAddress,
			UserAgent:         result.UserAgent,
			RequestID:         result.RequestID,
			Status:            result.Status,
			ResponsiblePerson: result.ResponsiblePerson,
			BatchNo:           result.BatchNo,
			Remarks:           result.Remarks,
			CreatedAt:         result.CreatedAt,
		})
	}
}

func getAuditLogsByEntity(svc *services.AuditService) echo.HandlerFunc {
	return func(c echo.Context) error {
		entityType := c.Param("type")
		entityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid entity ID",
			})
		}

		result, err := svc.GetByEntity(entityType, entityID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"entity_type": entityType,
			"entity_id":   entityID,
			"total":       len(result),
			"data":        result,
		})
	}
}
