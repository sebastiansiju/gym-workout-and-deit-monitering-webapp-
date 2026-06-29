package stores

import "database/sql"

// WorkoutStore owns all SQL for the workout entity.
type WorkoutStore struct{ db *sql.DB }

func NewWorkoutStore(db *sql.DB) *WorkoutStore { return &WorkoutStore{db: db} }

// CountOn returns how many workouts the user started on the given calendar day
// (YYYY-MM-DD). Used by the food daily-stats view (cross-entity composition).
func (s *WorkoutStore) CountOn(uid int64, date string) (int, error) {
	var n int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM workouts WHERE user_id = ? AND substr(started_at, 1, 10) = ?`,
		uid, date,
	).Scan(&n)
	return n, err
}
