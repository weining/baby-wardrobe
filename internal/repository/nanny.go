package repository

import (
	"database/sql"
	"fmt"
)

type NannyConfig struct {
	MonthlySalary float64 `json:"monthly_salary"`
	FirstPayday   string  `json:"first_payday"` // "YYYY-MM-DD"
}

type NannyLeaveRange struct {
	ID        int64  `json:"id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// GetNannyConfig 读取阿姨工资配置，未配置时返回零值
func GetNannyConfig(db *sql.DB) (NannyConfig, error) {
	var c NannyConfig
	err := db.QueryRow(`SELECT monthly_salary, first_payday FROM nanny_config WHERE id = 1`).
		Scan(&c.MonthlySalary, &c.FirstPayday)
	if err == sql.ErrNoRows {
		return NannyConfig{}, nil
	}
	if err != nil {
		return NannyConfig{}, fmt.Errorf("get nanny config: %w", err)
	}
	return c, nil
}

// SaveNannyConfig 保存（upsert）阿姨工资配置
func SaveNannyConfig(db *sql.DB, c NannyConfig) error {
	_, err := db.Exec(`
		INSERT INTO nanny_config (id, monthly_salary, first_payday, updated_at)
		VALUES (1, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(id) DO UPDATE SET
			monthly_salary = excluded.monthly_salary,
			first_payday   = excluded.first_payday,
			updated_at     = excluded.updated_at
	`, c.MonthlySalary, c.FirstPayday)
	if err != nil {
		return fmt.Errorf("save nanny config: %w", err)
	}
	return nil
}

func ListNannyLeaveRanges(db *sql.DB) ([]NannyLeaveRange, error) {
	rows, err := db.Query(`
		SELECT id, start_date, end_date
		FROM nanny_leave_ranges
		ORDER BY start_date, end_date, id
	`)
	if err != nil {
		return nil, fmt.Errorf("list nanny leave ranges: %w", err)
	}
	defer rows.Close()

	var list []NannyLeaveRange
	for rows.Next() {
		var r NannyLeaveRange
		if err := rows.Scan(&r.ID, &r.StartDate, &r.EndDate); err != nil {
			return nil, fmt.Errorf("scan nanny leave range: %w", err)
		}
		list = append(list, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate nanny leave ranges: %w", err)
	}
	return list, nil
}

func AddNannyLeaveRange(db *sql.DB, r NannyLeaveRange) (int64, error) {
	res, err := db.Exec(`
		INSERT INTO nanny_leave_ranges (start_date, end_date)
		VALUES (?, ?)
	`, r.StartDate, r.EndDate)
	if err != nil {
		return 0, fmt.Errorf("add nanny leave range: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get nanny leave range id: %w", err)
	}
	return id, nil
}

func DeleteNannyLeaveRange(db *sql.DB, id int64) error {
	res, err := db.Exec(`DELETE FROM nanny_leave_ranges WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete nanny leave range: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get deleted nanny leave rows: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
