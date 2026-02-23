package repository

import (
	"context"
	"database/sql"

	"github.com/krauzx/gitright/internal/models"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	query := `
		INSERT INTO projects (user_id, github_id, full_name, priority, focus_tag, custom_summary, generated_summary, include_in_profile)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		project.UserID, project.GitHubID, project.FullName, project.Priority, project.FocusTag,
		project.CustomSummary, project.GeneratedSummary, project.IncludeInProfile,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
}

func (r *ProjectRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Project, error) {
	query := `
		SELECT id, user_id, github_id, full_name, priority, focus_tag, custom_summary,
		       generated_summary, include_in_profile, created_at, updated_at
		FROM projects
		WHERE user_id = $1
		ORDER BY priority ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err := rows.Scan(
			&project.ID, &project.UserID, &project.GitHubID, &project.FullName, &project.Priority,
			&project.FocusTag, &project.CustomSummary, &project.GeneratedSummary,
			&project.IncludeInProfile, &project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, rows.Err()
}

func (r *ProjectRepository) Update(ctx context.Context, project *models.Project) error {
	query := `
		UPDATE projects
		SET priority = $1, focus_tag = $2, custom_summary = $3, generated_summary = $4,
		    include_in_profile = $5, updated_at = NOW()
		WHERE id = $6
	`
	_, err := r.db.ExecContext(ctx, query,
		project.Priority, project.FocusTag, project.CustomSummary,
		project.GeneratedSummary, project.IncludeInProfile, project.ID,
	)
	return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM projects WHERE id = $1`, id)
	return err
}
