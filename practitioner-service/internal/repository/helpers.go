package repository

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/azuvicenna/selaras-services/practitioner-service/internal/domain"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return domain.ErrNotFound
	case errors.As(err, &pgErr) && pgErr.Code == "23505": // unique_violation
		return domain.ErrAlreadyExists
	}
	return err
}

func placeholders(n int) string {
	items := make([]string, n)
	for i := range items {
		items[i] = "$" + strconv.Itoa(i+1)
	}
	return strings.Join(items, ", ")
}

func assignments(columns []string) string {
	items := make([]string, len(columns))
	for i, c := range columns {
		items[i] = fmt.Sprintf("%s = $%d", c, i+1)
	}
	return strings.Join(items, ", ")
}
