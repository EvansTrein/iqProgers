package postgres

import (
	"context"
	"fmt"
)

func (s *PostgresDB) TestDB() (int, error) {
	s.log.Debug("Postgres test method")

	var result int
	err := s.db.QueryRow(context.Background(), "SELECT 1").Scan(&result)
	if err != nil {
		return 0, fmt.Errorf("failed to execute test query: %w", err)
	}
	return result, nil
}
