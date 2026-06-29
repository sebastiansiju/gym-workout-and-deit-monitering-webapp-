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

// ExercisePR is the best single set for an exercise (heaviest, then most reps).
type ExercisePR struct {
	Weight    float64
	Reps      int
	Date      string
	WorkoutID int64
}

// ExerciseHistoryPoint is a per-workout rollup for an exercise's progress chart.
type ExerciseHistoryPoint struct {
	Date        string  `json:"date"`
	MaxWeight   float64 `json:"max_weight"`
	TotalVolume float64 `json:"total_volume"`
	SetsCount   int     `json:"sets_count"`
}

// PRForExercise returns the user's PR set for an exercise, or sql.ErrNoRows if
// they have no qualifying (non-warmup, weighted) set. Reads workout tables — this
// is workout-derived analytics keyed by exercise, so it lives on WorkoutStore.
func (s *WorkoutStore) PRForExercise(uid, exerciseID int64) (ExercisePR, error) {
	var pr ExercisePR
	err := s.db.QueryRow(`
		SELECT s.weight, s.reps, w.started_at, w.id
		FROM sets s
		JOIN workout_exercises we ON we.id = s.workout_exercise_id
		JOIN workouts w ON w.id = we.workout_id
		WHERE w.user_id = ? AND we.exercise_id = ? AND s.is_warmup = 0 AND s.weight > 0
		ORDER BY s.weight DESC, s.reps DESC
		LIMIT 1
	`, uid, exerciseID).Scan(&pr.Weight, &pr.Reps, &pr.Date, &pr.WorkoutID)
	return pr, err
}

func (s *WorkoutStore) HistoryForExercise(uid, exerciseID int64, limit int) ([]ExerciseHistoryPoint, error) {
	rows, err := s.db.Query(`
		SELECT w.started_at, MAX(s.weight), SUM(s.reps * s.weight), COUNT(s.id)
		FROM sets s
		JOIN workout_exercises we ON we.id = s.workout_exercise_id
		JOIN workouts w ON w.id = we.workout_id
		WHERE w.user_id = ? AND we.exercise_id = ? AND s.is_warmup = 0
		GROUP BY w.id
		ORDER BY w.started_at DESC
		LIMIT ?
	`, uid, exerciseID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	history := []ExerciseHistoryPoint{}
	for rows.Next() {
		var p ExerciseHistoryPoint
		if err := rows.Scan(&p.Date, &p.MaxWeight, &p.TotalVolume, &p.SetsCount); err != nil {
			return nil, err
		}
		history = append(history, p)
	}
	return history, rows.Err()
}
