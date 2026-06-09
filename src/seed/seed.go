package seed

import (
	"database/sql"
	"fmt"
	"time"
)

func RunSeed(db *sql.DB) {
	clearTables(db)

	seedMemberBenefits(db)
	seedFlightSchedules(db)
	seedReservationRecords(db)
	seedCompanions(db)
	seedUsageVouchers(db)
	seedWaitlistEntries(db)
	seedVerificationResults(db)
	seedStatusTransitions(db)
	seedRuleConfigs(db)
	seedExceptionEvents(db)
	seedAuditLogs(db)

	fmt.Println("Seed data inserted successfully")
}

func clearTables(db *sql.DB) {
	tables := []string{
		"audit_logs", "exception_events", "status_transitions", "verification_results",
		"waitlist_entries", "usage_vouchers", "companions", "reservation_records",
		"flight_schedules", "member_benefits", "rule_configs",
	}

	for _, table := range tables {
		db.Exec(fmt.Sprintf("DELETE FROM %s", table))
	}
}

func seedMemberBenefits(db *sql.DB) {
	benefits := []struct {
		benefitNo, memberName, memberLevel, status, responsiblePerson, batchNo string
		remainingQuota, totalQuota int
		validFrom, validTo time.Time
	}{
		{"BNF-20240101-0001", "张明远", "platinum", "active", "李运营", "BATCH-2024-001", 15, 20, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)},
		{"BNF-20240101-0002", "王丽华", "gold", "active", "张服务", "BATCH-2024-001", 8, 10, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)},
		{"BNF-20240201-0003", "刘建国", "platinum", "active", "李运营", "BATCH-2024-002", 18, 20, time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)},
		{"BNF-20240215-0004", "陈晓燕", "silver", "active", "王客服", "BATCH-2024-002", 3, 5, time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC)},
		{"BNF-20240301-0005", "赵伟强", "gold", "active", "张服务", "BATCH-2024-003", 6, 10, time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)},
		{"BNF-20240315-0006", "孙美玲", "platinum", "expired", "李运营", "BATCH-2024-003", 0, 20, time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 3, 14, 0, 0, 0, 0, time.UTC)},
		{"BNF-20240401-0007", "周志刚", "regular", "active", "王客服", "BATCH-2024-004", 2, 3, time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)},
		{"BNF-20240415-0008", "吴晓红", "gold", "suspended", "张服务", "BATCH-2024-004", 4, 10, time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC), time.Date(2025, 4, 14, 0, 0, 0, 0, time.UTC)},
		{"BNF-20240501-0009", "郑海涛", "platinum", "active", "李运营", "BATCH-2024-005", 12, 20, time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 4, 30, 0, 0, 0, 0, time.UTC)},
		{"BNF-20240515-0010", "黄雅琴", "silver", "active", "王客服", "BATCH-2024-005", 4, 5, time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC)},
	}

	query := `INSERT INTO member_benefits (benefit_no, member_name, member_level, remaining_quota, total_quota, status, responsible_person, valid_from, valid_to, batch_no, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, b := range benefits {
		now := time.Now()
		db.Exec(query, b.benefitNo, b.memberName, b.memberLevel, b.remainingQuota, b.totalQuota, b.status, b.responsiblePerson, b.validFrom, b.validTo, b.batchNo, now, now)
	}
}

func seedFlightSchedules(db *sql.DB) {
	schedules := []struct {
		flightNo, departure, arrival, flightStatus, status, responsiblePerson, batchNo string
		scheduledDepart, scheduledArrive time.Time
		capacity, used int
	}{
		{"CA1234", "北京首都T3", "上海虹桥T2", "scheduled", "active", "调度员张三", "FLT-2024-001", time.Date(2024, 6, 15, 8, 0, 0, 0, time.UTC), time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC), 50, 35},
		{"MU5678", "上海浦东T1", "广州白云T2", "boarding", "active", "调度员李四", "FLT-2024-001", time.Date(2024, 6, 15, 9, 30, 0, 0, time.UTC), time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC), 40, 40},
		{"CZ9012", "广州白云T1", "成都双流T2", "departed", "active", "调度员王五", "FLT-2024-001", time.Date(2024, 6, 15, 7, 0, 0, 0, time.UTC), time.Date(2024, 6, 15, 9, 30, 0, 0, time.UTC), 45, 30},
		{"HU3456", "深圳宝安T3", "北京大兴T2", "delayed", "active", "调度员赵六", "FLT-2024-002", time.Date(2024, 6, 15, 14, 0, 0, 0, time.UTC), time.Date(2024, 6, 15, 17, 0, 0, 0, time.UTC), 50, 25},
		{"FM7890", "杭州萧山T3", "西安咸阳T3", "scheduled", "active", "调度员孙七", "FLT-2024-002", time.Date(2024, 6, 15, 11, 0, 0, 0, time.UTC), time.Date(2024, 6, 15, 13, 30, 0, 0, time.UTC), 35, 20},
		{"CA2468", "北京首都T3", "深圳宝安T3", "scheduled", "active", "调度员张三", "FLT-2024-003", time.Date(2024, 6, 16, 8, 30, 0, 0, time.UTC), time.Date(2024, 6, 16, 11, 30, 0, 0, time.UTC), 50, 28},
		{"MU1357", "上海虹桥T2", "昆明长水T1", "arrived", "archived", "调度员李四", "FLT-2024-003", time.Date(2024, 6, 14, 10, 0, 0, 0, time.UTC), time.Date(2024, 6, 14, 13, 0, 0, 0, time.UTC), 40, 38},
		{"CZ8642", "成都双流T2", "杭州萧山T3", "scheduled", "active", "调度员王五", "FLT-2024-004", time.Date(2024, 6, 16, 15, 0, 0, 0, time.UTC), time.Date(2024, 6, 16, 17, 30, 0, 0, time.UTC), 45, 15},
	}

	query := `INSERT INTO flight_schedules (flight_no, departure_airport, arrival_airport, scheduled_depart, scheduled_arrive, flight_status, vip_lounge_capacity, vip_lounge_used, status, responsible_person, batch_no, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, s := range schedules {
		now := time.Now()
		db.Exec(query, s.flightNo, s.departure, s.arrival, s.scheduledDepart, s.scheduledArrive, s.flightStatus, s.capacity, s.used, s.status, s.responsiblePerson, s.batchNo, now, now)
	}
}

func seedReservationRecords(db *sql.DB) {
	reservations := []struct {
		reservationNo, memberName, flightNo, vipLounge, status, responsiblePerson, batchNo string
		memberBenefitID, flightScheduleID int64
		reservationTime time.Time
		guestCount int
	}{
		{"RSV-20240615-0001", "张明远", "CA1234", "首都机场VIP厅A区", "confirmed", "服务台小李", "RSV-BATCH-001", 1, 1, time.Date(2024, 6, 14, 15, 0, 0, 0, time.UTC), 2},
		{"RSV-20240615-0002", "王丽华", "MU5678", "浦东机场VIP厅B区", "pending_review", "服务台小王", "RSV-BATCH-001", 2, 2, time.Date(2024, 6, 14, 16, 0, 0, 0, time.UTC), 3},
		{"RSV-20240615-0003", "刘建国", "CZ9012", "白云机场VIP厅C区", "completed", "服务台小张", "RSV-BATCH-001", 3, 3, time.Date(2024, 6, 14, 10, 0, 0, 0, time.UTC), 1},
		{"RSV-20240615-0004", "陈晓燕", "HU3456", "宝安机场VIP厅D区", "pending_supplement", "服务台小李", "RSV-BATCH-002", 4, 4, time.Date(2024, 6, 14, 17, 0, 0, 0, time.UTC), 2},
		{"RSV-20240615-0005", "赵伟强", "FM7890", "萧山机场VIP厅E区", "draft", "服务台小王", "RSV-BATCH-002", 5, 5, time.Date(2024, 6, 14, 18, 0, 0, 0, time.UTC), 1},
		{"RSV-20240615-0006", "周志刚", "CA2468", "首都机场VIP厅A区", "rejected", "服务台小张", "RSV-BATCH-003", 7, 6, time.Date(2024, 6, 14, 19, 0, 0, 0, time.UTC), 4},
		{"RSV-20240615-0007", "郑海涛", "CA1234", "首都机场VIP厅A区", "confirmed", "服务台小李", "RSV-BATCH-003", 9, 1, time.Date(2024, 6, 14, 20, 0, 0, 0, time.UTC), 2},
		{"RSV-20240615-0008", "黄雅琴", "MU5678", "浦东机场VIP厅B区", "cancelled", "服务台小王", "RSV-BATCH-003", 10, 2, time.Date(2024, 6, 14, 21, 0, 0, 0, time.UTC), 1},
		{"RSV-20240616-0009", "张明远", "CA2468", "首都机场VIP厅A区", "draft", "服务台小李", "RSV-BATCH-004", 1, 6, time.Date(2024, 6, 15, 8, 0, 0, 0, time.UTC), 3},
		{"RSV-20240616-0010", "刘建国", "CZ8642", "双流机场VIP厅F区", "pending_review", "服务台小张", "RSV-BATCH-004", 3, 8, time.Date(2024, 6, 15, 9, 0, 0, 0, time.UTC), 2},
	}

	query := `INSERT INTO reservation_records (reservation_no, member_benefit_id, member_name, flight_no, flight_schedule_id, vip_lounge_name, reservation_time, guest_count, status, responsible_person, batch_no, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, r := range reservations {
		now := time.Now()
		db.Exec(query, r.reservationNo, r.memberBenefitID, r.memberName, r.flightNo, r.flightScheduleID, r.vipLounge, r.reservationTime, r.guestCount, r.status, r.responsiblePerson, r.batchNo, now, now)
	}
}

func seedCompanions(db *sql.DB) {
	companions := []struct {
		companionNo, name, idCard, relation, verifyStatus, status, responsiblePerson, batchNo string
		reservationID int64
		isVipEligible bool
	}{
		{"CPN-20240615-0001", "李芳", "110101198501011234", "spouse", "passed", "active", "服务台小李", "CPN-BATCH-001", 1, true},
		{"CPN-20240615-0002", "张小明", "110101201001011235", "child", "passed", "active", "服务台小李", "CPN-BATCH-001", 1, false},
		{"CPN-20240615-0003", "王强", "310101198602021236", "colleague", "passed", "active", "服务台小王", "CPN-BATCH-001", 2, true},
		{"CPN-20240615-0004", "王丽", "310101198803031237", "spouse", "passed", "active", "服务台小王", "CPN-BATCH-001", 2, true},
		{"CPN-20240615-0005", "王小明", "310101201202041238", "child", "failed", "inactive", "服务台小王", "CPN-BATCH-001", 2, false},
		{"CPN-20240615-0006", "陈父", "440101195501051239", "parent", "passed", "active", "服务台小张", "CPN-BATCH-002", 4, true},
		{"CPN-20240615-0007", "郑妻", "330101199001061240", "spouse", "passed", "active", "服务台小李", "CPN-BATCH-003", 7, true},
		{"CPN-20240615-0008", "张同事", "110101198707071241", "colleague", "pending", "active", "服务台小李", "CPN-BATCH-004", 9, false},
	}

	query := `INSERT INTO companions (companion_no, reservation_id, companion_name, companion_id_card, relation_type, is_vip_eligible, verification_status, status, responsible_person, batch_no, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, c := range companions {
		now := time.Now()
		isVip := 0
		if c.isVipEligible {
			isVip = 1
		}
		db.Exec(query, c.companionNo, c.reservationID, c.name, c.idCard, c.relation, isVip, c.verifyStatus, c.status, c.responsiblePerson, c.batchNo, now, now)
	}
}

func seedUsageVouchers(db *sql.DB) {
	vouchers := []struct {
		voucherNo, voucherType, qrCode, verifyStatus, status, location, responsiblePerson, batchNo string
		reservationID, memberBenefitID int64
		validFrom, validTo time.Time
	}{
		{"VCH-20240615-0001", "single_use", "QR001-AABBCCDD", "valid", "active", "", "系统生成", "VCH-BATCH-001", 1, 1, time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 15, 23, 59, 59, 0, time.UTC)},
		{"VCH-20240615-0002", "single_use", "QR002-EEFFGGHH", "valid", "active", "", "系统生成", "VCH-BATCH-001", 2, 2, time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 15, 23, 59, 59, 0, time.UTC)},
		{"VCH-20240615-0003", "single_use", "QR003-IIJJKKLL", "used", "used", "白云机场VIP厅C区", "系统生成", "VCH-BATCH-001", 3, 3, time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 15, 23, 59, 59, 0, time.UTC)},
		{"VCH-20240615-0004", "single_use", "QR004-MMNNOOPP", "valid", "active", "", "系统生成", "VCH-BATCH-002", 4, 4, time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 15, 23, 59, 59, 0, time.UTC)},
		{"VCH-20240615-0005", "guest_pass", "QR005-QQRRSSTT", "valid", "active", "", "系统生成", "VCH-BATCH-003", 7, 9, time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 15, 23, 59, 59, 0, time.UTC)},
		{"VCH-20240615-0006", "single_use", "QR006-UUVVWWXX", "expired", "expired", "", "系统生成", "VCH-BATCH-003", 8, 10, time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 14, 23, 59, 59, 0, time.UTC)},
	}

	query := `INSERT INTO usage_vouchers (voucher_no, reservation_id, member_benefit_id, voucher_type, qr_code, valid_from, valid_to, verification_status, status, used_location, responsible_person, batch_no, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, v := range vouchers {
		now := time.Now()
		db.Exec(query, v.voucherNo, v.reservationID, v.memberBenefitID, v.voucherType, v.qrCode, v.validFrom, v.validTo, v.verifyStatus, v.status, v.location, v.responsiblePerson, v.batchNo, now, now)
	}
}

func seedWaitlistEntries(db *sql.DB) {
	entries := []struct {
		waitlistNo, memberName, memberLevel, status, responsiblePerson, batchNo string
		memberBenefitID, flightScheduleID int64
		waitingSince time.Time
		priorityScore, waitMins int
	}{
		{"WTL-20240615-0001", "孙美玲", "platinum", "waiting", "调度员张三", "WTL-BATCH-001", 6, 1, time.Date(2024, 6, 15, 7, 30, 0, 0, time.UTC), 65, 45},
		{"WTL-20240615-0002", "吴晓红", "gold", "waiting", "调度员李四", "WTL-BATCH-001", 8, 1, time.Date(2024, 6, 15, 8, 0, 0, 0, time.UTC), 45, 60},
		{"WTL-20240615-0003", "周志刚", "regular", "waiting", "调度员王五", "WTL-BATCH-001", 7, 1, time.Date(2024, 6, 15, 8, 30, 0, 0, time.UTC), 45, 30},
		{"WTL-20240615-0004", "周志刚", "regular", "waiting", "调度员王五", "WTL-BATCH-001", 7, 2, time.Date(2024, 6, 15, 9, 0, 0, 0, time.UTC), 25, 60},
		{"WTL-20240615-0005", "黄雅琴", "silver", "seated", "调度员赵六", "WTL-BATCH-001", 10, 2, time.Date(2024, 6, 15, 8, 15, 0, 0, time.UTC), 40, 0},
		{"WTL-20240615-0006", "张明远", "platinum", "waiting", "调度员孙七", "WTL-BATCH-002", 1, 5, time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC), 70, 20},
	}

	query := `INSERT INTO waitlist_entries (waitlist_no, member_benefit_id, flight_schedule_id, member_name, member_level, waiting_since, priority_score, estimated_wait_mins, status, responsible_person, batch_no, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, w := range entries {
		now := time.Now()
		db.Exec(query, w.waitlistNo, w.memberBenefitID, w.flightScheduleID, w.memberName, w.memberLevel, w.waitingSince, w.priorityScore, w.waitMins, w.status, w.responsiblePerson, w.batchNo, now, now)
	}
}

func seedVerificationResults(db *sql.DB) {
	results := []struct {
		verificationNo, verifyType, result, failureReason, status, responsiblePerson, batchNo string
		reservationID, memberBenefitID, flightScheduleID int64
		verifiedQuota, verifiedCompanions int
	}{
		{"VER-20240615-0001", "reservation", "passed", "", "confirmed", "核验员甲", "VER-BATCH-001", 1, 1, 1, 2, 2},
		{"VER-20240615-0002", "reservation", "pending", "", "draft", "核验员乙", "VER-BATCH-001", 2, 2, 2, 3, 3},
		{"VER-20240615-0003", "reservation", "passed", "", "confirmed", "核验员甲", "VER-BATCH-001", 3, 3, 3, 1, 0},
		{"VER-20240615-0004", "reservation", "failed", "缺少航班时段信息", "draft", "核验员乙", "VER-BATCH-002", 4, 4, 4, 2, 1},
		{"VER-20240615-0005", "checkin", "passed", "", "confirmed", "核验员甲", "VER-BATCH-003", 7, 9, 1, 2, 1},
		{"VER-20240615-0006", "companion", "failed", "同行人验证未通过", "rejected", "核验员乙", "VER-BATCH-003", 2, 2, 2, 3, 3},
	}

	query := `INSERT INTO verification_results (verification_no, reservation_id, member_benefit_id, flight_schedule_id, verification_type, result, failure_reason, verified_quota, verified_companions, status, responsible_person, batch_no, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, v := range results {
		now := time.Now()
		var flightID interface{} = nil
		if v.flightScheduleID > 0 {
			flightID = v.flightScheduleID
		}
		db.Exec(query, v.verificationNo, v.reservationID, v.memberBenefitID, flightID, v.verifyType, v.result, v.failureReason, v.verifiedQuota, v.verifiedCompanions, v.status, v.responsiblePerson, v.batchNo, now, now)
	}
}

func seedStatusTransitions(db *sql.DB) {
	transitions := []struct {
		transitionNo, entityType, fromStatus, toStatus, action, reason, operator, status, responsiblePerson, batchNo string
		entityID int64
	}{
		{"TRN-20240615-0001", "reservation", "draft", "pending_review", "submit", "", "服务台小李", "active", "服务台小李", "TRN-BATCH-001", 1},
		{"TRN-20240615-0002", "reservation", "pending_review", "confirmed", "approve", "", "审核员甲", "active", "审核员甲", "TRN-BATCH-001", 1},
		{"TRN-20240615-0003", "reservation", "draft", "pending_review", "submit", "", "服务台小王", "active", "服务台小王", "TRN-BATCH-001", 2},
		{"TRN-20240615-0004", "reservation", "pending_review", "rejected", "reject", "同行人数量超限", "审核员乙", "active", "审核员乙", "TRN-BATCH-002", 6},
		{"TRN-20240615-0005", "member_benefit", "active", "expired", "expire", "权益已过期", "系统自动", "active", "李运营", "TRN-BATCH-002", 6},
		{"TRN-20240615-0006", "waitlist", "waiting", "seated", "arrange", "", "调度员赵六", "active", "调度员赵六", "TRN-BATCH-003", 4},
	}

	query := `INSERT INTO status_transitions (transition_no, entity_type, entity_id, from_status, to_status, action, reason, operator, status, responsible_person, batch_no, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, t := range transitions {
		now := time.Now()
		db.Exec(query, t.transitionNo, t.entityType, t.entityID, t.fromStatus, t.toStatus, t.action, t.reason, t.operator, t.status, t.responsiblePerson, t.batchNo, now, now)
	}
}

func seedRuleConfigs(db *sql.DB) {
	rules := []struct {
		ruleNo, ruleName, ruleType, ruleValue, appliesToLevel, status, responsiblePerson, batchNo string
		thresholdValue float64
		isActive bool
		effectiveDate time.Time
	}{
		{"RUL-001", "白金会员配额限制", "quota_limit", "20", "platinum", "active", "规则管理员", "RUL-BATCH-001", 20, true, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"RUL-002", "金卡会员配额限制", "quota_limit", "10", "gold", "active", "规则管理员", "RUL-BATCH-001", 10, true, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"RUL-003", "银卡会员配额限制", "quota_limit", "5", "silver", "active", "规则管理员", "RUL-BATCH-001", 5, true, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"RUL-004", "白金会员同行人限制", "companion_limit", "4", "platinum", "active", "规则管理员", "RUL-BATCH-001", 4, true, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"RUL-005", "金卡会员同行人限制", "companion_limit", "3", "gold", "active", "规则管理员", "RUL-BATCH-001", 3, true, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"RUL-006", "候补超时时间", "waitlist_timeout", "120", "", "active", "规则管理员", "RUL-BATCH-001", 120, true, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"RUL-007", "白金会员优先级权重", "priority_weight", "40", "platinum", "active", "规则管理员", "RUL-BATCH-002", 40, true, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"RUL-008", "权益过期提前提醒天数", "expiry_days", "7", "", "active", "规则管理员", "RUL-BATCH-002", 7, true, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	query := `INSERT INTO rule_configs (rule_no, rule_name, rule_type, rule_value, threshold_value, applies_to_level, is_active, effective_date, status, responsible_person, batch_no, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, r := range rules {
		now := time.Now()
		isActive := 0
		if r.isActive {
			isActive = 1
		}
		db.Exec(query, r.ruleNo, r.ruleName, r.ruleType, r.ruleValue, r.thresholdValue, r.appliesToLevel, isActive, r.effectiveDate, r.status, r.responsiblePerson, r.batchNo, now, now)
	}
}

func seedExceptionEvents(db *sql.DB) {
	events := []struct {
		eventNo, eventType, entityType, triggerField, thresholdValue, actualValue, severity, handler, handlingResult, status, responsiblePerson, batchNo string
		entityID int64
	}{
		{"EVT-20240615-0001", "benefit_expired", "member_benefit", "valid_to", "2024-06-15", "2024-03-14", "high", "李运营", "已处理", "resolved", "李运营", "EVT-BATCH-001", 6},
		{"EVT-20240615-0002", "quota_exceeded", "member_benefit", "remaining_quota", "1", "0", "medium", "", "", "open", "张服务", "EVT-BATCH-001", 4},
		{"EVT-20240615-0003", "companion_violation", "reservation", "guest_count", "3", "4", "medium", "审核员乙", "已驳回预约", "resolved", "审核员乙", "EVT-BATCH-001", 6},
		{"EVT-20240615-0004", "waitlist_timeout", "waitlist", "waiting_since", "60", "90", "low", "", "", "handling", "调度员张三", "EVT-BATCH-002", 3},
		{"EVT-20240615-0005", "voucher_invalid", "usage_voucher", "valid_to", "2024-06-15", "2024-06-14", "low", "", "", "open", "系统", "EVT-BATCH-002", 6},
	}

	query := `INSERT INTO exception_events (event_no, event_type, entity_type, entity_id, trigger_field, threshold_value, actual_value, severity, handler, handling_result, status, responsible_person, batch_no, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, e := range events {
		now := time.Now()
		db.Exec(query, e.eventNo, e.eventType, e.entityType, e.entityID, e.triggerField, e.thresholdValue, e.actualValue, e.severity, e.handler, e.handlingResult, e.status, e.responsiblePerson, e.batchNo, now, now)
	}
}

func seedAuditLogs(db *sql.DB) {
	logs := []struct {
		logNo, operationType, entityType, entityNo, operator, operatorRole, status, responsiblePerson, batchNo string
		entityID int64
	}{
		{"AUD-20240615-0001", "create", "member_benefit", "BNF-20240101-0001", "李运营", "权益运营", "active", "李运营", "AUD-BATCH-001", 1},
		{"AUD-20240615-0002", "create", "reservation", "RSV-20240615-0001", "服务台小李", "机场服务台", "active", "服务台小李", "AUD-BATCH-001", 1},
		{"AUD-20240615-0003", "status_change", "reservation", "RSV-20240615-0001", "审核员甲", "客服支持", "active", "审核员甲", "AUD-BATCH-001", 1},
		{"AUD-20240615-0004", "batch_import", "reservation", "batch", "服务台小王", "机场服务台", "active", "服务台小王", "AUD-BATCH-001", 0},
		{"AUD-20240615-0005", "verify", "verification", "VER-20240615-0001", "核验员甲", "客服支持", "active", "核验员甲", "AUD-BATCH-002", 1},
		{"AUD-20240615-0006", "archive", "flight_schedule", "MU1357", "调度员李四", "机场服务台", "active", "调度员李四", "AUD-BATCH-002", 7},
	}

	query := `INSERT INTO audit_logs (log_no, operation_type, entity_type, entity_id, entity_no, operator, operator_role, status, responsible_person, batch_no, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, l := range logs {
		now := time.Now()
		db.Exec(query, l.logNo, l.operationType, l.entityType, l.entityID, l.entityNo, l.operator, l.operatorRole, l.status, l.responsiblePerson, l.batchNo, now)
	}
}
