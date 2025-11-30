package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sahared/llm-observability/internal/models"
)

// CreateUser creates a new user in the database
func (r *ClickHouseRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO llm_observability.users (
			id, organization_id, email, name, role, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.OrganizationID,
		user.Email,
		user.Name,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByEmail retrieves a user by email
func (r *ClickHouseRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, organization_id, email, name, role, created_at, updated_at
		FROM llm_observability.users
		WHERE email = ?
		LIMIT 1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %s", email)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *ClickHouseRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	query := `
		SELECT id, organization_id, email, name, role, created_at, updated_at
		FROM llm_observability.users
		WHERE id = ?
		LIMIT 1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %s", userID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// CreateOrganization creates a new organization
func (r *ClickHouseRepository) CreateOrganization(ctx context.Context, org *models.Organization) error {
	query := `
		INSERT INTO llm_observability.organizations (
			id, name, plan, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		org.ID,
		org.Name,
		org.Plan,
		org.CreatedAt,
		org.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}

	return nil
}

// GetOrganization retrieves an organization by ID
func (r *ClickHouseRepository) GetOrganization(ctx context.Context, orgID string) (*models.Organization, error) {
	query := `
		SELECT id, name, plan, created_at, updated_at
		FROM llm_observability.organizations
		WHERE id = ?
		LIMIT 1
	`

	var org models.Organization
	err := r.db.QueryRowContext(ctx, query, orgID).Scan(
		&org.ID,
		&org.Name,
		&org.Plan,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("organization not found: %s", orgID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query organization: %w", err)
	}

	return &org, nil
}

// CreateProject creates a new project
func (r *ClickHouseRepository) CreateProject(ctx context.Context, project *models.Project) error {
	query := `
		INSERT INTO llm_observability.projects (
			id, organization_id, name, description, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		project.ID,
		project.OrganizationID,
		project.Name,
		project.Description,
		project.CreatedAt,
		project.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}

// GetProject retrieves a project by ID
func (r *ClickHouseRepository) GetProject(ctx context.Context, projectID string) (*models.Project, error) {
	query := `
		SELECT id, organization_id, name, description, created_at, updated_at
		FROM llm_observability.projects
		WHERE id = ?
		LIMIT 1
	`

	var project models.Project
	err := r.db.QueryRowContext(ctx, query, projectID).Scan(
		&project.ID,
		&project.OrganizationID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found: %s", projectID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query project: %w", err)
	}

	return &project, nil
}

// GetProjectsByOrg retrieves all projects for an organization
func (r *ClickHouseRepository) GetProjectsByOrg(ctx context.Context, orgID string) ([]*models.Project, error) {
	query := `
		SELECT id, organization_id, name, description, created_at, updated_at
		FROM llm_observability.projects
		WHERE organization_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(
			&project.ID,
			&project.OrganizationID,
			&project.Name,
			&project.Description,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, &project)
	}

	return projects, nil
}
