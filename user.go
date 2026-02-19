// ============================================================
// USER REPOSITORY — internal/repository/user.go
// ============================================================
// This layer ONLY talks to the database. No business logic here.
// Services call this. This runs SQL.
//
// GO CONCEPTS TO LEARN:
//   - interfaces: UserRepository defines the contract
//   - sqlx: wraps standard sql, lets you scan rows into structs
//   - sql.ErrNoRows: what you get when SELECT finds nothing
//   - context: always pass ctx to DB calls (enables timeouts)
//
// ARCHITECTURE: Handler → Service → Repository → Database
//   Each layer only knows about the layer directly below it.
//
// TODO (Phase 1): Implement FindByEmail and Create first.
//   These are needed by the auth service.
// ============================================================

package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/yourusername/movie-reservation/internal/models"
)

// UserRepository defines all DB operations for users
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	Create(ctx context.Context, params CreateUserParams) (*models.User, error)
	UpdateRole(ctx context.Context, userID, role string) error // for admin promotion
}

type CreateUserParams struct {
	Email        string
	PasswordHash string
	Role         string
}

// userRepository is the concrete implementation
type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

// ── FIND BY EMAIL ─────────────────────────────────────────────
// TODO (Phase 1): Implement this — used by Login
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	// TODO: Write the SQL query
	// query := `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`
	// err := r.db.GetContext(ctx, &user, query, email)
	// if err == sql.ErrNoRows {
	//   return nil, nil // not found, return nil without error
	// }
	// if err != nil {
	//   return nil, fmt.Errorf("finding user by email: %w", err)
	// }
	// return &user, nil

	_ = user
	return nil, fmt.Errorf("TODO: implement FindByEmail")
}

// ── FIND BY ID ────────────────────────────────────────────────
func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	// TODO: Similar to FindByEmail but WHERE id = $1
	// query := `SELECT id, email, password_hash, role, created_at FROM users WHERE id = $1`

	_ = user
	_ = sql.ErrNoRows
	return nil, fmt.Errorf("TODO: implement FindByID")
}

// ── CREATE ────────────────────────────────────────────────────
// TODO (Phase 1): Implement this — used by Signup
func (r *userRepository) Create(ctx context.Context, params CreateUserParams) (*models.User, error) {
	var user models.User

	// LEARN: RETURNING * gets the created row back in one query
	// TODO:
	// query := `
	//   INSERT INTO users (id, email, password_hash, role, created_at)
	//   VALUES (gen_random_uuid(), $1, $2, $3, NOW())
	//   RETURNING id, email, password_hash, role, created_at
	// `
	// err := r.db.QueryRowxContext(ctx, query, params.Email, params.PasswordHash, params.Role).
	//   StructScan(&user)
	// if err != nil {
	//   return nil, fmt.Errorf("creating user: %w", err)
	// }
	// return &user, nil

	_ = user
	return nil, fmt.Errorf("TODO: implement Create")
}

// ── UPDATE ROLE ───────────────────────────────────────────────
// TODO (Phase 1 - RBAC): Admin can promote users
func (r *userRepository) UpdateRole(ctx context.Context, userID, role string) error {
	// TODO:
	// query := `UPDATE users SET role = $1 WHERE id = $2`
	// _, err := r.db.ExecContext(ctx, query, role, userID)
	// return err
	return fmt.Errorf("TODO: implement UpdateRole")
}
