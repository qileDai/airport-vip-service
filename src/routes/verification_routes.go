package routes

import (
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func setupVerificationRoutes(g *echo.Group, svc *services.VerificationService) {
	r := g.Group("/verifications")

	r.GET("", listVerifications(svc))
	r.POST("", createVerification(svc))
	r.GET("/:id", getVerification(svc))
	r.POST("/verify-eligibility", verifyEligibility(svc))
}

func listVerifications(svc *services.VerificationService) echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
		if perPage <= 0 {
			perPage = 20
		}

		status := c.QueryParam("status")
		resultFilter := c.QueryParam("result")

		result, err := svc.List(page, perPage, status, resultFilter)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func createVerification(svc *services.VerificationService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req schemas.VerificationResultCreateRequest
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

func getVerification(svc *services.VerificationService) echo.HandlerFunc {
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
				Message: "Verification result not found",
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func verifyEligibility(svc *services.VerificationService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req schemas.EligibilityVerificationRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid request body",
			})
		}

		result, err := svc.VerifyEligibility(&req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}
