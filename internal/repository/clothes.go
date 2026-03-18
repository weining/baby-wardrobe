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
	`)
	if err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	return db, nil
}

type FilterParams struct {
	Category string
	Season   string
	Status   string
}

func ListClothes(db *sql.DB, f FilterParams) ([]model.Cloth, error) {
	query := `SELECT id, name, category, size, season, status, color, brand,
	          photo_path, thumb_path, notes, created_at, updated_at
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
	    photo_path, thumb_path, notes, created_at, updated_at
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

	c.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	c.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
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
