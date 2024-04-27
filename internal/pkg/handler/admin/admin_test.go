package admin

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSetDeductionValueHandler(t *testing.T) {
	// Arrange
	testCases := []struct {
		ptype           string
		reqBody         string
		sqlFn           func() (*sql.DB, error)
		expectedResBody string
	}{
		{"personal", `{"amount":70000.0}`, func() (*sql.DB, error) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				return nil, err
			}
			row := sqlmock.NewRows([]string{"maxAmount"}).AddRow(70000.0)
			mock.ExpectQuery("INSERT INTO \"deductions\" (\"name\", maxAmount) VALUES ($1, $2) ON CONFLICT (\"name\") DO UPDATE SET maxAmount = EXCLUDED.maxAmount RETURNING maxAmount;").WithArgs("personal", 70000.0).WillReturnRows(row)
			return db, err
		}, `{"personalDeduction":70000.0}`},
		{"k-receipt", `{"amount":80000.0}`, func() (*sql.DB, error) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				return nil, err
			}
			row := sqlmock.NewRows([]string{"maxAmount"}).AddRow(80000.0)
			mock.ExpectQuery("INSERT INTO \"deductions\" (\"name\", maxAmount) VALUES ($1, $2) ON CONFLICT (\"name\") DO UPDATE SET maxAmount = EXCLUDED.maxAmount RETURNING maxAmount;").WithArgs("k-receipt", 80000.0).WillReturnRows(row)
			return db, err
		}, `{"kReceipt":80000.0}`},
	}

	for _, tc := range testCases {
		// Act
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:type")
		c.SetParamNames("type")
		c.SetParamValues(tc.ptype)

		db, err := tc.sqlFn()
		h := New(db)
		// Assertions
		assert.NoError(t, err)
		if assert.NoError(t, h.SetDeductionValueHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, tc.expectedResBody, rec.Body.String())
		}

	}
}

func TestErrorSetDeductionValueHandler(t *testing.T) {
	testCases := []struct {
		ptype       string
		reqBody     string
		sqlFn       func() (*sql.DB, error)
		expectedErr error
	}{
		{"personal", `{"amount":9999.0}`, func() (*sql.DB, error) {
			db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				return nil, err
			}
			return db, err
		}, echo.NewHTTPError(http.StatusBadRequest, "amount must between 10,000 - 100,000")},
		{"personal", `{"amount":100001.0}`, func() (*sql.DB, error) {
			db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				return nil, err
			}
			return db, err
		}, echo.NewHTTPError(http.StatusBadRequest, "amount must between 10,000 - 100,000")},
		{"k-receipt", `{"amount":-1.0}`, func() (*sql.DB, error) {
			db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				return nil, err
			}
			return db, err
		}, echo.NewHTTPError(http.StatusBadRequest, "amount must between 0 - 100,000")},
		{"k-receipt", `{"amount":100001.0}`, func() (*sql.DB, error) {
			db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				return nil, err
			}
			return db, err
		}, echo.NewHTTPError(http.StatusBadRequest, "amount must between 0 - 100,000")},
	}

	for _, tc := range testCases {
		// Act
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/:type")
		c.SetParamNames("type")
		c.SetParamValues(tc.ptype)

		db, err := tc.sqlFn()
		h := New(db)
		// Assertions
		assert.NoError(t, err)
		terr := h.SetDeductionValueHandler(c)

		// assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, terr, tc.expectedErr)

	}
}
