package repositories

import (
	"airport-vip-service/src/models"
	"database/sql"
	"fmt"
	"time"
)

type UsageVoucherRepository struct {
	db *sql.DB
}

func NewUsageVoucherRepository(db *sql.DB) *UsageVoucherRepository {
	return &UsageVoucherRepository{db: db}
}

func (r *UsageVoucherRepository) Create(voucher *models.UsageVoucher) (int64, error) {
	query := `INSERT INTO usage_vouchers 
		(voucher_no, reservation_id, member_benefit_id, voucher_type, qr_code, valid_from, valid_to,
		used_at, used_location, verification_status, status, responsible_person, batch_no, remarks, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query,
		voucher.VoucherNo, voucher.ReservationID, voucher.MemberBenefitID, voucher.VoucherType,
		voucher.QRCode, voucher.ValidFrom, voucher.ValidTo, voucher.UsedAt,
		voucher.UsedLocation, voucher.VerificationStatus, voucher.Status,
		voucher.ResponsiblePerson, voucher.BatchNo, voucher.Remarks, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create usage voucher: %w", err)
	}

	return result.LastInsertId()
}

func (r *UsageVoucherRepository) GetByID(id int64) (*models.UsageVoucher, error) {
	query := `SELECT id, voucher_no, reservation_id, member_benefit_id, voucher_type, qr_code, valid_from, valid_to,
		used_at, used_location, verification_status, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM usage_vouchers WHERE id = ?`

	voucher := &models.UsageVoucher{}
	err := r.db.QueryRow(query, id).Scan(
		&voucher.ID, &voucher.VoucherNo, &voucher.ReservationID, &voucher.MemberBenefitID,
		&voucher.VoucherType, &voucher.QRCode, &voucher.ValidFrom, &voucher.ValidTo,
		&voucher.UsedAt, &voucher.UsedLocation, &voucher.VerificationStatus, &voucher.Status,
		&voucher.ResponsiblePerson, &voucher.BatchNo, &voucher.Remarks,
		&voucher.CreatedAt, &voucher.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get usage voucher: %w", err)
	}

	return voucher, nil
}

func (r *UsageVoucherRepository) GetByVoucherNo(voucherNo string) (*models.UsageVoucher, error) {
	query := `SELECT id, voucher_no, reservation_id, member_benefit_id, voucher_type, qr_code, valid_from, valid_to,
		used_at, used_location, verification_status, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM usage_vouchers WHERE voucher_no = ?`

	voucher := &models.UsageVoucher{}
	err := r.db.QueryRow(query, voucherNo).Scan(
		&voucher.ID, &voucher.VoucherNo, &voucher.ReservationID, &voucher.MemberBenefitID,
		&voucher.VoucherType, &voucher.QRCode, &voucher.ValidFrom, &voucher.ValidTo,
		&voucher.UsedAt, &voucher.UsedLocation, &voucher.VerificationStatus, &voucher.Status,
		&voucher.ResponsiblePerson, &voucher.BatchNo, &voucher.Remarks,
		&voucher.CreatedAt, &voucher.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get usage voucher: %w", err)
	}

	return voucher, nil
}

func (r *UsageVoucherRepository) List(offset, limit int, reservationID int64, status string) ([]models.UsageVoucher, int64, error) {
	var args []interface{}
	whereClause := "WHERE 1=1"

	if reservationID > 0 {
		whereClause += " AND reservation_id = ?"
		args = append(args, reservationID)
	}
	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	countQuery := "SELECT COUNT(*) FROM usage_vouchers " + whereClause
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count usage vouchers: %w", err)
	}

	query := `SELECT id, voucher_no, reservation_id, member_benefit_id, voucher_type, qr_code, valid_from, valid_to,
		used_at, used_location, verification_status, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM usage_vouchers ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list usage vouchers: %w", err)
	}
	defer rows.Close()

	var vouchers []models.UsageVoucher
	for rows.Next() {
		var voucher models.UsageVoucher
		err := rows.Scan(
			&voucher.ID, &voucher.VoucherNo, &voucher.ReservationID, &voucher.MemberBenefitID,
			&voucher.VoucherType, &voucher.QRCode, &voucher.ValidFrom, &voucher.ValidTo,
			&voucher.UsedAt, &voucher.UsedLocation, &voucher.VerificationStatus, &voucher.Status,
			&voucher.ResponsiblePerson, &voucher.BatchNo, &voucher.Remarks,
			&voucher.CreatedAt, &voucher.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan usage voucher: %w", err)
		}
		vouchers = append(vouchers, voucher)
	}

	return vouchers, total, nil
}

func (r *UsageVoucherRepository) Update(id int64, voucher *models.UsageVoucher) error {
	query := `UPDATE usage_vouchers SET 
		voucher_type = ?, qr_code = ?, valid_from = ?, valid_to = ?, used_at = ?,
		used_location = ?, verification_status = ?, status = ?, responsible_person = ?,
		batch_no = ?, remarks = ?, updated_at = ?
		WHERE id = ?`

	_, err := r.db.Exec(query,
		voucher.VoucherType, voucher.QRCode, voucher.ValidFrom, voucher.ValidTo, voucher.UsedAt,
		voucher.UsedLocation, voucher.VerificationStatus, voucher.Status, voucher.ResponsiblePerson,
		voucher.BatchNo, voucher.Remarks, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update usage voucher: %w", err)
	}

	return nil
}

func (r *UsageVoucherRepository) Delete(id int64) error {
	query := `DELETE FROM usage_vouchers WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete usage voucher: %w", err)
	}
	return nil
}

func (r *UsageVoucherRepository) MarkAsUsed(id int64, location string) error {
	query := `UPDATE usage_vouchers SET used_at = ?, used_location = ?, status = 'used', verification_status = 'used', updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, time.Now(), location, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to mark voucher as used: %w", err)
	}
	return nil
}

func (r *UsageVoucherRepository) GetByReservationID(reservationID int64) ([]models.UsageVoucher, error) {
	query := `SELECT id, voucher_no, reservation_id, member_benefit_id, voucher_type, qr_code, valid_from, valid_to,
		used_at, used_location, verification_status, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM usage_vouchers WHERE reservation_id = ? ORDER BY created_at`

	rows, err := r.db.Query(query, reservationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vouchers by reservation: %w", err)
	}
	defer rows.Close()

	var vouchers []models.UsageVoucher
	for rows.Next() {
		var voucher models.UsageVoucher
		err := rows.Scan(
			&voucher.ID, &voucher.VoucherNo, &voucher.ReservationID, &voucher.MemberBenefitID,
			&voucher.VoucherType, &voucher.QRCode, &voucher.ValidFrom, &voucher.ValidTo,
			&voucher.UsedAt, &voucher.UsedLocation, &voucher.VerificationStatus, &voucher.Status,
			&voucher.ResponsiblePerson, &voucher.BatchNo, &voucher.Remarks,
			&voucher.CreatedAt, &voucher.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan usage voucher: %w", err)
		}
		vouchers = append(vouchers, voucher)
	}

	return vouchers, nil
}
