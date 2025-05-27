package dao

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/melkdesousa/gamgo/dao/models"
)

type AccountDAO struct {
	connection *pgx.Conn
}

func NewAccountDAO(connection *pgx.Conn) *AccountDAO {
	return &AccountDAO{connection: connection}
}

func (dao *AccountDAO) GetUserByEmail(email string) (*models.Account, error) {
	account := &models.Account{}
	var (
		rawID           string
		rawName         string
		rawEmail        string
		rawPasswordHash string
		rawCreatedAt    time.Time
		rawDeletedAt    *time.Time
		rawIsActive     bool
	)
	stmt, err := dao.connection.Prepare(context.Background(), "get_user_by_email", `
		SELECT
			id, name, email, passwordHash, createdAt, deletedAt, isActive
		FROM accounts WHERE email = $1 AND deletedAt IS NULL AND isActive = true
	`)
	if err != nil {
		return nil, err
	}
	err = dao.connection.
		QueryRow(context.Background(), stmt.Name, email).
		Scan(
			&rawID,
			&rawName,
			&rawEmail,
			&rawPasswordHash,
			&rawCreatedAt,
			&rawDeletedAt,
			&rawIsActive,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	account.ID = uuid.MustParse(rawID)
	account.Name = rawName
	account.Email = rawEmail
	account.PasswordHash = rawPasswordHash
	account.CreatedAt = rawCreatedAt
	if rawDeletedAt != nil {
		account.DeletedAt = *rawDeletedAt
	} else {
		account.DeletedAt = time.Time{} // Set to zero time if deletedAt is NULL
	}
	account.IsActive = rawIsActive
	return account, nil
}
