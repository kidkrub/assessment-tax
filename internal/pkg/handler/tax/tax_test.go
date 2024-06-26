package tax

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

func TestTaxCalculate(t *testing.T) {
	// Arrange
	maxDeductions := map[string]float64{"personal": 60000.0, "k-receipt": 50000.0}
	testCases := []struct {
		inputData     TaxRequestObject
		maxDeductions map[string]float64
		expected      struct {
			tax       float64
			taxlevels []TaxLevel
		}
	}{
		{TaxRequestObject{60000.0, 0.0, nil}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{0.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 0.0},
			{"500,001-1,000,000", 0.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{500000.0, 0.0, nil}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{29000.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 29000.0},
			{"500,001-1,000,000", 0.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{560000.0, 0.0, nil}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{35000.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 35000.0},
			{"500,001-1,000,000", 0.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{1060000.0, 0.0, nil}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{110000.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 35000.0},
			{"500,001-1,000,000", 75000.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{2060000.0, 0.0, nil}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{310000.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 35000.0},
			{"500,001-1,000,000", 75000.0},
			{"1,000,001-2,000,000", 200000.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{2060001.0, 0.0, nil}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{310000.35, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 35000.0},
			{"500,001-1,000,000", 75000.0},
			{"1,000,001-2,000,000", 200000.0},
			{"2,000,001 ขึ้นไป", 0.35},
		}}},
		{TaxRequestObject{150000.0, 1000.0, nil}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{-1000.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 0.0},
			{"500,001-1,000,000", 0.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{500000.0, 25000.0, nil}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{4000.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 29000.0},
			{"500,001-1,000,000", 0.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{500000.0, 29000.0, nil}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{0.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 29000.0},
			{"500,001-1,000,000", 0.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{500000.0, 30000.0, nil}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{-1000.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 29000.0},
			{"500,001-1,000,000", 0.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{500000.0, 0.0, []Allowance{{"donation", 200000.0}}}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{19000.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 19000.0},
			{"500,001-1,000,000", 0.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{500000.0, 0.0, []Allowance{{"donation", 100000.0}}}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{19000.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 19000.0},
			{"500,001-1,000,000", 0.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{500000.0, 0.0, []Allowance{{"donation", 50000.0}}}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{24000.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 24000.0},
			{"500,001-1,000,000", 0.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{500000.0, 0.0, []Allowance{{"k-receipt", 200000.0}, {"donation", 100000.0}}}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{14000.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 14000.0},
			{"500,001-1,000,000", 0.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{500000.0, 0.0, []Allowance{{"k-receipt", 50000.0}, {"donation", 100000.0}}}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{14000.0, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 14000.0},
			{"500,001-1,000,000", 0.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
		{TaxRequestObject{500000.0, 0.0, []Allowance{{"k-receipt", 49999.0}, {"donation", 100000.0}}}, maxDeductions, struct {
			tax       float64
			taxlevels []TaxLevel
		}{14000.1, []TaxLevel{
			{"0-150,000", 0.0},
			{"150,001-500,000", 14000.1},
			{"500,001-1,000,000", 0.0},
			{"1,000,001-2,000,000", 0.0},
			{"2,000,001 ขึ้นไป", 0.0},
		}}},
	}

	// Act & Assert
	for _, tc := range testCases {
		// Act
		actualTax, actualLevels := taxCalculate(tc.inputData, tc.maxDeductions)
		// Assert
		assert.Equal(t, tc.expected.tax, actualTax, "tax calculation is incorrect for %.2f case", tc.inputData.TotalIncome)
		assert.Equal(t, tc.expected.taxlevels, actualLevels, "levels calculation is incorrect for %.2f case", tc.inputData.TotalIncome)
	}
}

func TestTaxCalculateHandler(t *testing.T) {
	// Arrange
	sqlFunc := func() (*sql.DB, error) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			return nil, err
		}
		pRow := sqlmock.NewRows([]string{"maxAmount"}).AddRow(60000.0)
		kRow := sqlmock.NewRows([]string{"maxAmount"}).AddRow(50000.0)
		mock.ExpectQuery("SELECT 'maxAmount' FROM \"dedictions\" WHERE \"name\" = $1;").WithArgs("personal").WillReturnRows(pRow)
		mock.ExpectQuery("SELECT 'maxAmount' FROM \"dedictions\" WHERE \"name\" = $1;").WithArgs("k-receipt").WillReturnRows(kRow)
		return db, err
	}
	testCases := []struct {
		reqBody         string
		sqlFn           func() (*sql.DB, error)
		expectedResBody string
	}{
		{`{"totalIncome":60000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, sqlFunc, `{"tax":0.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":0.0},{"level":"500,001-1,000,000","tax":0.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":500000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, sqlFunc, `{"tax":29000.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":29000.0},{"level":"500,001-1,000,000","tax":0.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":560000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, sqlFunc, `{"tax":35000.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":35000.0},{"level":"500,001-1,000,000","tax":0.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":1060000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, sqlFunc, `{"tax":110000.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":35000.0},{"level":"500,001-1,000,000","tax":75000.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":2060000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, sqlFunc, `{"tax":310000.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":35000.0},{"level":"500,001-1,000,000","tax":75000.0},{"level":"1,000,001-2,000,000","tax":200000.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":2060001.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, sqlFunc, `{"tax":310000.35,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":35000.0},{"level":"500,001-1,000,000","tax":75000.0},{"level":"1,000,001-2,000,000","tax":200000.0},{"level":"2,000,001 ขึ้นไป","tax":0.35}]}`},
		{`{"totalIncome":150000.0,"wht":1000.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, sqlFunc, `{"tax":0.0,"taxRefund":1000.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":0.0},{"level":"500,001-1,000,000","tax":0.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":500000.0,"wht":25000.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, sqlFunc, `{"tax":4000.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":29000.0},{"level":"500,001-1,000,000","tax":0.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":500000.0,"wht":29000.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, sqlFunc, `{"tax":0.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":29000.0},{"level":"500,001-1,000,000","tax":0.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":500000.0,"wht":30000.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, sqlFunc, `{"tax":0.0,"taxRefund":1000.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":29000.0},{"level":"500,001-1,000,000","tax":0.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":500000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":200000.0}]}`, sqlFunc, `{"tax":19000.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":19000.0},{"level":"500,001-1,000,000","tax":0.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":500000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":100000.0}]}`, sqlFunc, `{"tax":19000.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":19000.0},{"level":"500,001-1,000,000","tax":0.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":500000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":50000.0}]}`, sqlFunc, `{"tax":24000.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":24000.0},{"level":"500,001-1,000,000","tax":0.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":500000.0,"wht":0.0,"allowances":[{"allowanceType":"k-receipt","amount":200000.0},{"allowanceType":"donation","amount":100000.0}]}`, sqlFunc, `{"tax":14000.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":14000.0},{"level":"500,001-1,000,000","tax":0.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":500000.0,"wht":0.0,"allowances":[{"allowanceType":"k-receipt","amount":50000.0},{"allowanceType":"donation","amount":100000.0}]}`, sqlFunc, `{"tax":14000.0,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":14000.0},{"level":"500,001-1,000,000","tax":0.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
		{`{"totalIncome":500000.0,"wht":0.0,"allowances":[{"allowanceType":"k-receipt","amount":49999.0},{"allowanceType":"donation","amount":100000.0}]}`, sqlFunc, `{"tax":14000.1,"taxLevel":[{"level":"0-150,000","tax":0.0},{"level":"150,001-500,000","tax":14000.1},{"level":"500,001-1,000,000","tax":0.0},{"level":"1,000,001-2,000,000","tax":0.0},{"level":"2,000,001 ขึ้นไป","tax":0.0}]}`},
	}

	for _, tc := range testCases {
		// Act
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, err := tc.sqlFn()
		h := New(db)
		// Assertions
		assert.NoError(t, err)
		if assert.NoError(t, h.TaxCalculateHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, tc.expectedResBody, rec.Body.String())
		}

	}
}

// TODO : write upload test
