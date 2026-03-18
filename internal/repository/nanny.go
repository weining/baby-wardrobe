package repository

import (
	"database/sql"
	"fmt"
)

type NannyConfig struct {
	MonthlySalary float64 `json:"monthly_salary"`
	FirstPayday   string  `json:"first_payday"` // "YYYY-MM-DD"
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
