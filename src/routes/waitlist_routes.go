package routes

import (
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func setupWaitlistRoutes(g *echo.Group, svc *services.WaitlistService) {
	r := g.Group("/waitlist")

	r.GET("", listWaitlistEntries(svc))
	r.POST("", createWaitlistEntry(svc))
	r.GET("/:id", getWaitlistEntry(svc))
	r.PUT("/:id", updateWaitlistEntry(svc))
	r.DELETE("/:id", deleteWaitlistEntry(svc))
	r.POST("/calculate-priority", calculatePriority(svc))
	r.POST("/arrange", arrangeWaitlist(svc))
}

func listWaitlistEntries(svc *services.WaitlistService) echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
		if perPage <= 0 {
			perPage = 20
		}

		flightScheduleID, _ := strconv.ParseInt(c.QueryParam("flight_schedule_id"), 10, 64)
		status := c.QueryParam("status")

		result, err := svc.List(page, perPage, flightScheduleID, status)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func createWaitlistEntry(svc *services.WaitlistService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req schemas.WaitlistEntryCreateRequest
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

func getWaitlistEntry(svc *services.WaitlistService) echo.HandlerFunc {
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
				Message: "Waitlist entry not found",
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func updateWaitlistEntry(svc *services.WaitlistService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid ID",
			})
		}

		var req schemas.WaitlistEntryUpdateRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid request body",
			})
		}

		entry, err := svc.GetByID(id)
		if err != nil || entry == nil {
			return c.JSON(http.StatusNotFound, schemas.ErrorResponse{
				Code:    404,
				Message: "Waitlist entry not found",
			})
		}

		if req.PriorityScore != nil {
			entry.PriorityScore = *req.PriorityScore
		}
		if req.EstimatedWaitMins != nil {
			entry.EstimatedWaitMins = *req.EstimatedWaitMins
		}
		if req.Status != "" {
			entry.Status = req.Status
		}

		return c.JSON(http.StatusOK, entry)
	}
}

func deleteWaitlistEntry(svc *services.WaitlistService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid ID",
			})
		}

		if err := svc.Delete(id); err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, schemas.SuccessResponse{
			Success: true,
			Message: "Waitlist entry deleted successfully",
		})
	}
}

func calculatePriority(svc *services.WaitlistService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req schemas.PriorityCalculationRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid request body",
			})
		}

		result, err := svc.CalculatePriority(&req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func arrangeWaitlist(svc *services.WaitlistService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req schemas.WaitlistArrangementRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid request body",
			})
		}

		result, err := svc.ArrangeWaitlist(&req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}
