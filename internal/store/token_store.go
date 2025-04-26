package store

import (
	"database/sql"
	"github.com/helmigandi/go-workout-api/internal/tokens"
	"time"
)

type PostgresTokenStore struct {
	db *sql.DB
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

func (p *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) VALUES ($1, $2, $3, $4)`

	_, err := p.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

func (p *PostgresTokenStore) CreateToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = p.Insert(token)
	return token, err
}

func (p *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	query := `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`
	_, err := p.db.Exec(query, userID, scope)
	return err
}
