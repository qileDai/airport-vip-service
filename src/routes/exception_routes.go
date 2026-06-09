package routes

import (
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func setupExceptionRoutes(g *echo.Group, svc *services.ExceptionService) {
	r := g.Group("/exceptions")

	r.GET("", listExceptions(svc))
	r.POST("", createException(svc))
	r.GET("/:id", getException(svc))
	r.POST("/handle", handleException(svc))
	r.GET("/open", getOpenExceptions(svc))
}

func listExceptions(svc *services.ExceptionService) echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
		if perPage <= 0 {
			perPage = 20
		}

		eventType := c.QueryParam("event_type")
		status := c.QueryParam("status")
		severity := c.QueryParam("severity")

		result, err := svc.List(page, perPage, eventType, status, severity)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func createException(svc *services.ExceptionService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req schemas.ExceptionEventCreateRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid request body",
			})
		}

		result, err := svc.Create(&req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, result)
	}
}

func getException(svc *services.ExceptionService) echo.HandlerFunc {
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
				Message: "Exception event not found",
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func handleException(svc *services.ExceptionService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req schemas.ExceptionHandleRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid request body",
			})
		}

		result, err := svc.Handle(&req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func getOpenExceptions(svc *services.ExceptionService) echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := svc.GetOpenEvents()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		var data []schemas.ExceptionEventResponse
		for _, e := range result {
			data = append(data, schemas.ExceptionEventResponse{
				ID:                e.ID,
				EventNo:           e.EventNo,
				EventType:         e.EventType,
				EntityType:        e.EntityType,
				EntityID:          e.EntityID,
				TriggerField:      e.TriggerField,
				ThresholdValue:    e.ThresholdValue,
				ActualValue:       e.ActualValue,
				Severity:          e.Severity,
				Handler:           e.Handler,
				HandlingDeadline:  e.HandlingDeadline,
				HandledAt:         e.HandledAt,
				HandlingResult:    e.HandlingResult,
				Status:            e.Status,
				ResponsiblePerson: e.ResponsiblePerson,
				BatchNo:           e.BatchNo,
				Remarks:           e.Remarks,
				CreatedAt:         e.CreatedAt,
				UpdatedAt:         e.UpdatedAt,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"total": len(data),
			"data":  data,
		})
	}
}
