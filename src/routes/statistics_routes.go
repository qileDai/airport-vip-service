package routes

import (
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func setupStatisticsRoutes(g *echo.Group, svc *services.StatisticsService) {
	r := g.Group("/statistics")

	r.GET("/verification", getVerificationStatistics(svc))
	r.GET("/waitlist", getWaitlistStatistics(svc))
	r.GET("/usage", getUsageStatistics(svc))
}

func getVerificationStatistics(svc *services.StatisticsService) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &schemas.StatisticsRequest{
			StartDate: c.QueryParam("start_date"),
			EndDate:   c.QueryParam("end_date"),
			BatchNo:   c.QueryParam("batch_no"),
			Role:      c.QueryParam("role"),
		}

		result, err := svc.GetVerificationStatistics(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func getWaitlistStatistics(svc *services.StatisticsService) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &schemas.StatisticsRequest{
			StartDate: c.QueryParam("start_date"),
			EndDate:   c.QueryParam("end_date"),
			BatchNo:   c.QueryParam("batch_no"),
			Role:      c.QueryParam("role"),
		}

		result, err := svc.GetWaitlistStatistics(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func getUsageStatistics(svc *services.StatisticsService) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &schemas.StatisticsRequest{
			StartDate: c.QueryParam("start_date"),
			EndDate:   c.QueryParam("end_date"),
			BatchNo:   c.QueryParam("batch_no"),
			Role:      c.QueryParam("role"),
		}

		result, err := svc.GetUsageStatistics(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}
