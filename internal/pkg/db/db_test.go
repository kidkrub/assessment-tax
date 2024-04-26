package db

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetDeductionValue(t *testing.T) {
	// Arrange
	testCases := []struct {
		key      string
		sqlFn    func() (*sql.DB, error)
		expected float64
	}{{"personal", func() (*sql.DB, error) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			return nil, err
		}
		row := sqlmock.NewRows([]string{"maxAmount"}).AddRow(60000.0)
		mock.ExpectQuery("SELECT 'maxAmount' FROM \"dedictions\" WHERE \"name\" = $1;").WithArgs("personal").WillReturnRows(row)
		return db, err
	}, 60000.0}, {"k-receipt", func() (*sql.DB, error) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			return nil, err
		}
		row := sqlmock.NewRows([]string{"maxAmount"}).AddRow(50000.0)
		mock.ExpectQuery("SELECT 'maxAmount' FROM \"dedictions\" WHERE \"name\" = $1;").WithArgs("k-receipt").WillReturnRows(row)
		return db, err
	}, 50000.0}}

	// Act & Assert
	for _, tc := range testCases {
		db, err := tc.sqlFn()
		actualMaxAmount := GetDeductionValue(db, tc.key)

		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actualMaxAmount)
	}

}
