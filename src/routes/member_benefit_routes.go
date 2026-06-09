package routes

import (
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func setupMemberBenefitRoutes(g *echo.Group, svc *services.MemberBenefitService) {
	r := g.Group("/member-benefits")

	r.GET("", listMemberBenefits(svc))
	r.POST("", createMemberBenefit(svc))
	r.GET("/:id", getMemberBenefit(svc))
	r.PUT("/:id", updateMemberBenefit(svc))
	r.DELETE("/:id", deleteMemberBenefit(svc))
	r.POST("/check-expiry", checkBenefitExpiry(svc))
}

func listMemberBenefits(svc *services.MemberBenefitService) echo.HandlerFunc {
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
		memberLevel := c.QueryParam("member_level")

		result, err := svc.List(page, perPage, status, memberLevel)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func createMemberBenefit(svc *services.MemberBenefitService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req schemas.MemberBenefitCreateRequest
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

func getMemberBenefit(svc *services.MemberBenefitService) echo.HandlerFunc {
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
				Message: "Member benefit not found",
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func updateMemberBenefit(svc *services.MemberBenefitService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid ID",
			})
		}

		var req schemas.MemberBenefitUpdateRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid request body",
			})
		}

		result, err := svc.Update(id, &req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func deleteMemberBenefit(svc *services.MemberBenefitService) echo.HandlerFunc {
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
			Message: "Member benefit deleted successfully",
		})
	}
}

func checkBenefitExpiry(svc *services.MemberBenefitService) echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := svc.CheckExpiry()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}
