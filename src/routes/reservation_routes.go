package routes

import (
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func setupReservationRoutes(g *echo.Group, svc *services.ReservationService) {
	r := g.Group("/reservations")

	r.GET("", listReservations(svc))
	r.POST("", createReservation(svc))
	r.GET("/:id", getReservation(svc))
	r.PUT("/:id", updateReservation(svc))
	r.DELETE("/:id", deleteReservation(svc))
	r.POST("/batch-import/preview", previewBatchImport(svc))
	r.POST("/batch-import", executeBatchImport(svc))
	r.POST("/:id/status", changeReservationStatus(svc))
	r.GET("/:id/transitions", getAllowedTransitions(svc))
}

func listReservations(svc *services.ReservationService) echo.HandlerFunc {
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
		memberBenefitID, _ := strconv.ParseInt(c.QueryParam("member_benefit_id"), 10, 64)

		result, err := svc.List(page, perPage, status, memberBenefitID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func createReservation(svc *services.ReservationService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req schemas.ReservationRecordCreateRequest
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

func getReservation(svc *services.ReservationService) echo.HandlerFunc {
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
				Message: "Reservation not found",
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func updateReservation(svc *services.ReservationService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid ID",
			})
		}

		var req schemas.ReservationRecordUpdateRequest
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

func deleteReservation(svc *services.ReservationService) echo.HandlerFunc {
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
			Message: "Reservation deleted successfully",
		})
	}
}

func previewBatchImport(svc *services.ReservationService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req schemas.BatchImportRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid request body",
			})
		}

		result, err := svc.BatchImportPreview(&req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func executeBatchImport(svc *services.ReservationService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req schemas.BatchImportRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid request body",
			})
		}

		result, err := svc.BatchImport(&req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func changeReservationStatus(svc *services.ReservationService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid ID",
			})
		}

		var req schemas.StatusChangeRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid request body",
			})
		}

		req.EntityType = "reservation"
		req.EntityID = id

		result, err := svc.ChangeStatus(&req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}

func getAllowedTransitions(svc *services.ReservationService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Code:    400,
				Message: "Invalid ID",
			})
		}

		reservation, err := svc.GetByID(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
		}

		if reservation == nil {
			return c.JSON(http.StatusNotFound, schemas.ErrorResponse{
				Code:    404,
				Message: "Reservation not found",
			})
		}

		result := svc.GetAllowedTransitions(reservation.Status)
		return c.JSON(http.StatusOK, result)
	}
}
