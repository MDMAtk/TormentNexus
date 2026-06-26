package skillregistry

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"time"
)

type SkillOutcome struct {
	SkillName    string    `json:"skillName"`
	SuccessCount int       `json:"successCount"`
	FailureCount int       `json:"failureCount"`
	WinRate      float64   `json:"winRate"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func RecordOutcome(ctx context.Context, db *sql.DB, skillName string, success bool) error {
	var successInc, failureInc int
	if success {
		successInc = 1
	} else {
		failureInc = 1
	}

	_, err := db.ExecContext(ctx, `
		INSERT INTO skill_outcomes (skill_name, success_count, failure_count, win_rate, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(skill_name) DO UPDATE SET
			success_count = skill_outcomes.success_count + ?,
			failure_count = skill_outcomes.failure_count + ?,
			win_rate = CAST(skill_outcomes.success_count + ? AS REAL) / CAST(skill_outcomes.success_count + skill_outcomes.failure_count + 1 AS REAL),
			updated_at = CURRENT_TIMESTAMP
	`, skillName, successInc, failureInc, 1.0, successInc, failureInc, successInc)
	return err
}

func GetLowPerformingSkills(ctx context.Context, db *sql.DB, threshold float64) ([]SkillOutcome, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT skill_name, success_count, failure_count, win_rate, updated_at
		FROM skill_outcomes
		WHERE win_rate < ? AND (success_count + failure_count) >= 5
	`, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SkillOutcome
	for rows.Next() {
		var s SkillOutcome
		var updateStr string
		if err := rows.Scan(&s.SkillName, &s.SuccessCount, &s.FailureCount, &s.WinRate, &updateStr); err == nil {
			s.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updateStr)
			results = append(results, s)
		}
	}
	return results, nil
}

func EvolveSkills(ctx context.Context, workspaceRoot string, db *sql.DB) (int, error) {
	lowSkills, err := GetLowPerformingSkills(ctx, db, 0.5)
	if err != nil {
		return 0, err
	}

	deactivatedCount := 0
	for _, skill := range lowSkills {
		// Move low-performing tools from go/internal/tools/ to go/internal/tools/_disabled/
		src := filepath.Join(workspaceRoot, "go", "internal", "tools", skill.SkillName+".go")
		dest := filepath.Join(workspaceRoot, "go", "internal", "tools", "_disabled", skill.SkillName+".go")
		
		if _, statErr := os.Stat(src); statErr == nil {
			// File exists, let's move it (deactivate it)
			if renameErr := os.Rename(src, dest); renameErr == nil {
				deactivatedCount++
				// Reset DB status
				_, _ = db.ExecContext(ctx, "UPDATE mcp_servers SET status='pending', notes=? WHERE name=?", "auto-deactivated due to low win-rate", skill.SkillName)
			}
		}
	}

	return deactivatedCount, nil
}
