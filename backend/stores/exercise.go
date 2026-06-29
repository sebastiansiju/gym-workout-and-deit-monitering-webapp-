package stores

import (
	"database/sql"
	"encoding/json"

	"github.com/Cawlumm/lyftr-backend/models"
	"github.com/Cawlumm/lyftr-backend/seed"
)

// ExerciseStore owns all SQL for the (global, read-only) exercises catalog, and
// wraps the seed subsystem for the admin sync/reset endpoints.
type ExerciseStore struct{ db *sql.DB }

func NewExerciseStore(db *sql.DB) *ExerciseStore { return &ExerciseStore{db: db} }

// ExerciseFilter holds the optional list filters (empty string = no filter).
type ExerciseFilter struct {
	Query, MuscleGroup, Category, Equipment string
	Limit                                   int
}

const exerciseSelect = `SELECT id, name, muscle_group, secondary_muscles, category, equipment, description, image_url FROM exercises`

type scanner interface{ Scan(dest ...any) error }

func scanExercise(row scanner) models.Exercise {
	var e models.Exercise
	var secondaryRaw string
	row.Scan(&e.ID, &e.Name, &e.MuscleGroup, &secondaryRaw, &e.Category, &e.Equipment, &e.Description, &e.ImageURL)
	json.Unmarshal([]byte(secondaryRaw), &e.SecondaryMuscles)
	if e.SecondaryMuscles == nil {
		e.SecondaryMuscles = []string{}
	}
	return e
}

func (s *ExerciseStore) List(f ExerciseFilter) ([]models.Exercise, error) {
	q := exerciseSelect + ` WHERE 1=1`
	args := []any{}
	if f.Query != "" {
		q += " AND name LIKE ?"
		args = append(args, "%"+f.Query+"%")
	}
	if f.MuscleGroup != "" {
		q += " AND muscle_group = ?"
		args = append(args, f.MuscleGroup)
	}
	if f.Category != "" {
		q += " AND category = ?"
		args = append(args, f.Category)
	}
	if f.Equipment != "" {
		q += " AND equipment = ?"
		args = append(args, f.Equipment)
	}
	q += " ORDER BY name LIMIT ?"
	args = append(args, f.Limit)

	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	exercises := []models.Exercise{}
	for rows.Next() {
		exercises = append(exercises, scanExercise(rows))
	}
	return exercises, rows.Err()
}

// Get returns one exercise, or sql.ErrNoRows if not found.
func (s *ExerciseStore) Get(id int64) (models.Exercise, error) {
	e := scanExercise(s.db.QueryRow(exerciseSelect+` WHERE id = ?`, id))
	if e.ID == 0 {
		return models.Exercise{}, sql.ErrNoRows
	}
	return e, nil
}

func (s *ExerciseStore) Count() (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM exercises`).Scan(&n)
	return n, err
}

// Sync / SeedStatus / Reset delegate to the seed subsystem (which owns the
// ExerciseDB fetch + the concurrency guard) so the controller stays SQL-free.
func (s *ExerciseStore) Sync() error                 { return seed.SyncExercises(s.db) }
func (s *ExerciseStore) SeedStatus() seed.SeedStatus { return seed.GetSeedStatus(s.db) }
func (s *ExerciseStore) Reset() error                { return seed.WipeAndReseed(s.db) }
