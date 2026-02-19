// ============================================================
// RESERVATION SERVICE — internal/services/reservation.go
// ============================================================
// THIS IS THE MOST IMPORTANT FILE IN THE PROJECT
// Interviewers WILL ask you about everything in here.
//
// THE CORE PROBLEM:
//   Two users click seat A1 at the same time.
//   Both check the DB → both see it's AVAILABLE.
//   Both write → both think they got it.
//   = DOUBLE BOOKING (race condition)
//
// THE SOLUTION (3-step flow):
//   1. LOCK IN REDIS    → atomic SET NX with TTL
//      Redis processes commands one at a time → no race condition
//      NX = "only set if Not eXists" → atomic check-and-set
//
//   2. CONFIRM IN DB    → Postgres transaction
//      A transaction is all-or-nothing (ACID)
//      If anything fails mid-way, everything rolls back
//
//   3. BACKGROUND CLEANUP → Asynq job auto-cancels unpaid seats after TTL
//
// GO CONCEPTS TO LEARN:
//   - context with timeout: ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
//   - error wrapping: fmt.Errorf("lock seats: %w", err)
//   - Redis SetNX: atomic "set if not exists" command
//   - DB transactions in Go: tx, err := db.BeginTxx(ctx, nil)
//
// TODO ORDER:
//   1. LockSeats()
//   2. ConfirmReservation()
//   3. GetAvailableSeats()
// ============================================================

package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"github.com/yourusername/movie-reservation/internal/models"
	"github.com/yourusername/movie-reservation/internal/repository"
)

var (
	ErrSeatUnavailable = errors.New("one or more seats are unavailable")
	ErrReservationNotFound = errors.New("reservation not found")
	ErrNotYourReservation = errors.New("reservation belongs to another user")
)

const seatLockTTL = 10 * time.Minute

type ReservationService interface {
	LockSeats(ctx context.Context, userID string, req models.LockSeatsRequest) (string, error)
	ConfirmReservation(ctx context.Context, reservationID, userID string) error
	GetAvailableSeats(ctx context.Context, showtimeID string) ([]models.Seat, error)
	CancelReservation(ctx context.Context, reservationID string) error // called by background job
}

type reservationService struct {
	db          *sqlx.DB
	rdb         *redis.Client
	seatRepo    repository.SeatRepository
	resRepo     repository.ReservationRepository
}

func NewReservationService(db *sqlx.DB, rdb *redis.Client, seatRepo repository.SeatRepository, resRepo repository.ReservationRepository) ReservationService {
	return &reservationService{db: db, rdb: rdb, seatRepo: seatRepo, resRepo: resRepo}
}

// ── LOCK SEATS ────────────────────────────────────────────────
// TODO (Phase 2 - Step 1): IMPLEMENT THIS FIRST
//
// Called when user clicks "Reserve" (before payment).
// Returns reservationID so client can use it to confirm later.
func (s *reservationService) LockSeats(ctx context.Context, userID string, req models.LockSeatsRequest) (string, error) {
	// Step 1: Check seats are AVAILABLE in the DB
	// TODO: seats, err := s.seatRepo.FindByIDs(ctx, req.SeatIDs)
	// TODO: for _, seat := range seats {
	//   if seat.Status != "AVAILABLE" {
	//     return "", ErrSeatUnavailable
	//   }
	// }

	// Step 2: Try to lock each seat in Redis using SET NX
	// NX = "only set if key does Not eXist" = atomic check + set
	// LEARN: Why Redis and not just the DB for locking?
	//   Redis is in-memory = microseconds. DB = milliseconds.
	//   For seat selection, users expect instant feedback.
	//
	// TODO: for _, seatID := range req.SeatIDs {
	//   lockKey := fmt.Sprintf("seat_lock:%s", seatID)
	//   // SetNX returns true if the key was set (lock acquired)
	//   // Returns false if key already exists (someone else has the lock)
	//   set, err := s.rdb.SetNX(ctx, lockKey, userID, seatLockTTL).Result()
	//   if err != nil {
	//     return "", fmt.Errorf("redis lock: %w", err)
	//   }
	//   if !set {
	//     // Seat is locked by someone else — release any locks we just set
	//     // TODO: release previously acquired locks (rollback)
	//     return "", ErrSeatUnavailable
	//   }
	// }

	// Step 3: Create a PENDING reservation in DB
	// TODO: reservationID, err := s.resRepo.Create(ctx, repository.CreateReservationParams{
	//   UserID:    userID,
	//   SeatIDs:   req.SeatIDs,
	//   ExpiresAt: time.Now().Add(seatLockTTL),
	// })

	// Step 4: Schedule background job to auto-cancel if not paid
	// (done in the handler or a job scheduler)

	return "", fmt.Errorf("TODO: implement LockSeats")
}

// ── CONFIRM RESERVATION ───────────────────────────────────────
// TODO (Phase 2 - Step 2): Implement after LockSeats works
//
// Called after user pays. Moves from "locked in Redis" to "reserved in DB".
// MUST use a DB transaction — all-or-nothing.
//
// LEARN: What is a DB transaction?
//   A set of operations that either ALL succeed or ALL fail (rollback).
//   ACID: Atomic, Consistent, Isolated, Durable
func (s *reservationService) ConfirmReservation(ctx context.Context, reservationID, userID string) error {
	// Step 1: Start a DB transaction
	// LEARN: tx.Rollback() is safe to call even after tx.Commit() — it's a no-op
	// TODO: tx, err := s.db.BeginTxx(ctx, nil)
	// TODO: if err != nil { return fmt.Errorf("begin tx: %w", err) }
	// TODO: defer tx.Rollback() // rolls back if we return early with an error

	// Step 2: Get the reservation and verify ownership
	// TODO: res, err := s.resRepo.FindByIDTx(ctx, tx, reservationID)
	// TODO: if err != nil { return ErrReservationNotFound }
	// TODO: if res.UserID != userID { return ErrNotYourReservation }
	// TODO: if res.Status != "PENDING" { return errors.New("reservation already processed") }

	// Step 3: Update seats to RESERVED in DB
	// TODO: seatIDs := extract seatIDs from reservation
	// TODO: err = s.seatRepo.UpdateStatusTx(ctx, tx, seatIDs, "RESERVED")
	// TODO: if err != nil { return fmt.Errorf("updating seats: %w", err) }

	// Step 4: Update reservation status to CONFIRMED
	// TODO: err = s.resRepo.UpdateStatusTx(ctx, tx, reservationID, "CONFIRMED")
	// TODO: if err != nil { return fmt.Errorf("confirming reservation: %w", err) }

	// Step 5: Commit the transaction (makes all changes permanent)
	// TODO: if err = tx.Commit(); err != nil { return fmt.Errorf("commit: %w", err) }

	// Step 6: Clean up Redis locks (no longer needed, DB is source of truth now)
	// TODO: for _, seatID := range seatIDs {
	//   s.rdb.Del(ctx, fmt.Sprintf("seat_lock:%s", seatID))
	// }

	return fmt.Errorf("TODO: implement ConfirmReservation")
}

// ── GET AVAILABLE SEATS ───────────────────────────────────────
// TODO (Phase 2): Returns seats that are available (not in DB as RESERVED, not locked in Redis)
func (s *reservationService) GetAvailableSeats(ctx context.Context, showtimeID string) ([]models.Seat, error) {
	// Step 1: Get all seats for this showtime from DB
	// TODO: seats, err := s.seatRepo.FindByShowtime(ctx, showtimeID)

	// Step 2: For each AVAILABLE seat, check if it's locked in Redis
	// If Redis lock exists → treat as unavailable (temporarily held by another user)
	// TODO: var available []models.Seat
	// TODO: for _, seat := range seats {
	//   if seat.Status != "AVAILABLE" {
	//     continue
	//   }
	//   lockKey := fmt.Sprintf("seat_lock:%s", seat.ID)
	//   exists, _ := s.rdb.Exists(ctx, lockKey).Result()
	//   if exists == 0 {
	//     available = append(available, seat)
	//   }
	// }
	// return available, nil

	return nil, fmt.Errorf("TODO: implement GetAvailableSeats")
}

// ── CANCEL RESERVATION ────────────────────────────────────────
// Called by background job when reservation expires unpaid
func (s *reservationService) CancelReservation(ctx context.Context, reservationID string) error {
	// TODO: Same pattern as ConfirmReservation but set status to CANCELLED
	// TODO: and set seats back to AVAILABLE

	return fmt.Errorf("TODO: implement CancelReservation")
}

// keep compiler happy
var _ = fmt.Sprintf
var _ = time.Now
