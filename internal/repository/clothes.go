package repository

import (
	"database/sql"
	"fmt"
	"time"

	"clothes-manager/internal/model"
	_ "modernc.org/sqlite"
)

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS clothes (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT NOT NULL,
			category   TEXT NOT NULL,
			size       TEXT NOT NULL,
			season     TEXT NOT NULL,
			status     TEXT NOT NULL DEFAULT 'wearing',
			color      TEXT,
			brand      TEXT,
			photo_path TEXT,
			thumb_path TEXT,
			notes      TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS categories (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT NOT NULL UNIQUE,
			sort_order INTEGER DEFAULT 0
		);
		CREATE TABLE IF NOT EXISTS statuses (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			value      TEXT NOT NULL UNIQUE,
			label      TEXT NOT NULL,
			color      TEXT NOT NULL DEFAULT 'gray',
			sort_order INTEGER DEFAULT 0
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("create tables: %w", err)
	}

	if err := seedDefaults(db); err != nil {
		return nil, fmt.Errorf("seed defaults: %w", err)
	}

	return db, nil
}

// seedDefaults 首次运行时写入默认类别和状态
func seedDefaults(db *sql.DB) error {
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM categories`).Scan(&count)
	if count == 0 {
		for i, name := range model.DefaultCategories {
			_, err := db.Exec(`INSERT INTO categories (name, sort_order) VALUES (?, ?)`, name, i)
			if err != nil {
				return fmt.Errorf("seed category %s: %w", name, err)
			}
		}
	}

	db.QueryRow(`SELECT COUNT(*) FROM statuses`).Scan(&count)
	if count == 0 {
		for i, s := range model.DefaultStatuses {
			_, err := db.Exec(`INSERT INTO statuses (value, label, color, sort_order) VALUES (?, ?, ?, ?)`,
				s.Value, s.Label, s.Color, i)
			if err != nil {
				return fmt.Errorf("seed status %s: %w", s.Value, err)
			}
		}
	}
	return nil
}

// --- Category CRUD ---

func ListCategories(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`SELECT name FROM categories ORDER BY sort_order, id`)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()
	var list []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		list = append(list, name)
	}
	return list, rows.Err()
}

func AddCategory(db *sql.DB, name string) error {
	_, err := db.Exec(`INSERT INTO categories (name) VALUES (?)`, name)
	if err != nil {
		return fmt.Errorf("add category: %w", err)
	}
	return nil
}

func DeleteCategory(db *sql.DB, name string) error {
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM clothes WHERE category = ?`, name).Scan(&count)
	if count > 0 {
		return fmt.Errorf("类别「%s」下还有 %d 件衣物，无法删除", name, count)
	}
	_, err := db.Exec(`DELETE FROM categories WHERE name = ?`, name)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}

// --- Status CRUD ---

func ListStatuses(db *sql.DB) ([]model.Status, error) {
	rows, err := db.Query(`SELECT id, value, label, color FROM statuses ORDER BY sort_order, id`)
	if err != nil {
		return nil, fmt.Errorf("list statuses: %w", err)
	}
	defer rows.Close()
	var list []model.Status
	for rows.Next() {
		var s model.Status
		if err := rows.Scan(&s.ID, &s.Value, &s.Label, &s.Color); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func AddStatus(db *sql.DB, s model.Status) error {
	_, err := db.Exec(`INSERT INTO statuses (value, label, color) VALUES (?, ?, ?)`,
		s.Value, s.Label, s.Color)
	if err != nil {
		return fmt.Errorf("add status: %w", err)
	}
	return nil
}

func UpdateStatus(db *sql.DB, s model.Status) error {
	_, err := db.Exec(`UPDATE statuses SET label=?, color=? WHERE id=?`,
		s.Label, s.Color, s.ID)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	return nil
}

func DeleteStatus(db *sql.DB, id int64) error {
	var value string
	db.QueryRow(`SELECT value FROM statuses WHERE id=?`, id).Scan(&value)
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM clothes WHERE status = ?`, value).Scan(&count)
	if count > 0 {
		return fmt.Errorf("状态下还有 %d 件衣物，无法删除", count)
	}
	_, err := db.Exec(`DELETE FROM statuses WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete status: %w", err)
	}
	return nil
}

type FilterParams struct {
	Category string
	Season   string
	Status   string
}

func ListClothes(db *sql.DB, f FilterParams) ([]model.Cloth, error) {
	query := `SELECT id, name, category, size, season, status, color, brand,
	          photo_path, thumb_path, notes,
	          strftime('%Y-%m-%d %H:%M:%S', created_at),
	          strftime('%Y-%m-%d %H:%M:%S', updated_at)
	          FROM clothes WHERE 1=1`
	args := []any{}

	if f.Category != "" {
		query += " AND category = ?"
		args = append(args, f.Category)
	}
	if f.Season != "" {
		query += " AND season = ?"
		args = append(args, f.Season)
	}
	if f.Status != "" {
		query += " AND status = ?"
		args = append(args, f.Status)
	}
	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list clothes: %w", err)
	}
	defer rows.Close()

	var list []model.Cloth
	for rows.Next() {
		c, err := scanCloth(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func GetCloth(db *sql.DB, id int64) (model.Cloth, error) {
	row := db.QueryRow(`SELECT id, name, category, size, season, status, color, brand,
	    photo_path, thumb_path, notes,
	    strftime('%Y-%m-%d %H:%M:%S', created_at),
	    strftime('%Y-%m-%d %H:%M:%S', updated_at)
	    FROM clothes WHERE id = ?`, id)
	return scanCloth(row)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanCloth(s scanner) (model.Cloth, error) {
	var c model.Cloth
	var color, brand, photoPath, thumbPath, notes sql.NullString
	var createdAt, updatedAt string

	err := s.Scan(&c.ID, &c.Name, &c.Category, &c.Size, &c.Season, &c.Status,
		&color, &brand, &photoPath, &thumbPath, &notes, &createdAt, &updatedAt)
	if err != nil {
		return c, fmt.Errorf("scan cloth: %w", err)
	}

	c.Color = color.String
	c.Brand = brand.String
	c.PhotoPath = photoPath.String
	c.ThumbPath = thumbPath.String
	c.Notes = notes.String

	c.CreatedAt, _ = parseTime(createdAt)
	c.UpdatedAt, _ = parseTime(updatedAt)
	return c, nil
}

func CreateCloth(db *sql.DB, c model.Cloth) (int64, error) {
	res, err := db.Exec(`INSERT INTO clothes
	    (name, category, size, season, status, color, brand, photo_path, thumb_path, notes)
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Name, c.Category, c.Size, c.Season, c.Status,
		c.Color, c.Brand, c.PhotoPath, c.ThumbPath, c.Notes)
	if err != nil {
		return 0, fmt.Errorf("create cloth: %w", err)
	}
	return res.LastInsertId()
}

func UpdateCloth(db *sql.DB, c model.Cloth) error {
	_, err := db.Exec(`UPDATE clothes SET
	    name=?, category=?, size=?, season=?, status=?, color=?, brand=?,
	    photo_path=?, thumb_path=?, notes=?, updated_at=CURRENT_TIMESTAMP
	    WHERE id=?`,
		c.Name, c.Category, c.Size, c.Season, c.Status,
		c.Color, c.Brand, c.PhotoPath, c.ThumbPath, c.Notes, c.ID)
	if err != nil {
		return fmt.Errorf("update cloth: %w", err)
	}
	return nil
}

func DeleteCloth(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM clothes WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete cloth: %w", err)
	}
	return nil
}

type Stats struct {
	Total    int
	Wearing  int
	Stored   int
	Outgrown int
	Donated  int
}

func GetStats(db *sql.DB) (Stats, error) {
	rows, err := db.Query(`SELECT status, COUNT(*) FROM clothes GROUP BY status`)
	if err != nil {
		return Stats{}, fmt.Errorf("get stats: %w", err)
	}
	defer rows.Close()

	var s Stats
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return s, err
		}
		s.Total += count
		switch status {
		case model.StatusWearing:
			s.Wearing = count
		case model.StatusStored:
			s.Stored = count
		case model.StatusOutgrown:
			s.Outgrown = count
		case model.StatusDonated:
			s.Donated = count
		}
	}
	return s, rows.Err()
}

var timeFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05",
	time.RFC3339,
}

func parseTime(s string) (time.Time, error) {
	for _, f := range timeFormats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse time: %q", s)
}
