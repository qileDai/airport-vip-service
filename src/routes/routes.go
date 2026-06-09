package routes

import (
	"airport-vip-service/src/repositories"
	"airport-vip-service/src/services"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Services struct {
	MemberBenefitSvc  *services.MemberBenefitService
	ReservationSvc    *services.ReservationService
	FlightScheduleSvc *services.FlightScheduleService
	WaitlistSvc       *services.WaitlistService
	VerificationSvc   *services.VerificationService
	ExceptionSvc      *services.ExceptionService
	AuditSvc          *services.AuditService
	StatisticsSvc     *services.StatisticsService
}

func SetupRoutes(e *echo.Echo, db *sql.DB) {
	repos := initRepositories(db)
	svcs := initServices(repos)

	api := e.Group("/api/v1")

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	setupMemberBenefitRoutes(api, svcs.MemberBenefitSvc)
	setupReservationRoutes(api, svcs.ReservationSvc)
	setupFlightScheduleRoutes(api, svcs.FlightScheduleSvc)
	setupWaitlistRoutes(api, svcs.WaitlistSvc)
	setupVerificationRoutes(api, svcs.VerificationSvc)
	setupExceptionRoutes(api, svcs.ExceptionSvc)
	setupAuditRoutes(api, svcs.AuditSvc)
	setupStatisticsRoutes(api, svcs.StatisticsSvc)
}

func initRepositories(db *sql.DB) map[string]interface{} {
	return map[string]interface{}{
		"memberBenefit":    repositories.NewMemberBenefitRepository(db),
		"reservation":      repositories.NewReservationRecordRepository(db),
		"flightSchedule":   repositories.NewFlightScheduleRepository(db),
		"companion":        repositories.NewCompanionRepository(db),
		"usageVoucher":     repositories.NewUsageVoucherRepository(db),
		"waitlist":         repositories.NewWaitlistEntryRepository(db),
		"verification":     repositories.NewVerificationResultRepository(db),
		"statusTransition": repositories.NewStatusTransitionRepository(db),
		"ruleConfig":       repositories.NewRuleConfigRepository(db),
		"exceptionEvent":   repositories.NewExceptionEventRepository(db),
		"auditLog":         repositories.NewAuditLogRepository(db),
	}
}

func initServices(repos map[string]interface{}) *Services {
	memberBenefitRepo := repos["memberBenefit"].(*repositories.MemberBenefitRepository)
	reservationRepo := repos["reservation"].(*repositories.ReservationRecordRepository)
	flightScheduleRepo := repos["flightSchedule"].(*repositories.FlightScheduleRepository)
	companionRepo := repos["companion"].(*repositories.CompanionRepository)
	waitlistRepo := repos["waitlist"].(*repositories.WaitlistEntryRepository)
	verificationRepo := repos["verification"].(*repositories.VerificationResultRepository)
	statusTransitionRepo := repos["statusTransition"].(*repositories.StatusTransitionRepository)
	ruleConfigRepo := repos["ruleConfig"].(*repositories.RuleConfigRepository)
	exceptionEventRepo := repos["exceptionEvent"].(*repositories.ExceptionEventRepository)
	auditLogRepo := repos["auditLog"].(*repositories.AuditLogRepository)

	return &Services{
		MemberBenefitSvc: services.NewMemberBenefitServiceWithRepos(
			memberBenefitRepo, exceptionEventRepo, auditLogRepo,
		),
		ReservationSvc: services.NewReservationServiceWithRepos(
			reservationRepo, memberBenefitRepo, flightScheduleRepo,
			companionRepo, statusTransitionRepo, auditLogRepo,
		),
		FlightScheduleSvc: services.NewFlightScheduleServiceWithRepos(
			flightScheduleRepo, auditLogRepo,
		),
		WaitlistSvc: services.NewWaitlistServiceWithRepos(
			waitlistRepo, memberBenefitRepo, flightScheduleRepo,
			ruleConfigRepo, auditLogRepo,
		),
		VerificationSvc: services.NewVerificationServiceWithRepos(
			verificationRepo, memberBenefitRepo, reservationRepo,
			flightScheduleRepo, companionRepo, ruleConfigRepo, auditLogRepo,
		),
		ExceptionSvc: services.NewExceptionServiceWithRepos(
			exceptionEventRepo, auditLogRepo,
		),
		AuditSvc: services.NewAuditServiceWithRepos(auditLogRepo),
		StatisticsSvc: services.NewStatisticsServiceWithRepos(
			verificationRepo, waitlistRepo, reservationRepo,
		),
	}
}
