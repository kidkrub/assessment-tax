package tax

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestTaxCalculate(t *testing.T) {
	// Arrange
	testCases := []struct {
		income   float64
		expected float64
	}{
		{60000.0, 0},
		{500000.0, 29000.0},
		{560000.0, 35000.0},
		{1060000.0, 110000.0},
		{2060000.0, 310000.0},
		{2060001.0, 310000.35}}

	// Act & Assert
	for _, tc := range testCases {
		// Act
		actualTax := taxCalculate(tc.income)
		// Assert
		assert.Equal(t, tc.expected, actualTax, "tax calculation is incorrect for %.2f case", tc.income)
	}
}

func TestTaxCalculateHandler(t *testing.T) {
	// Arrange
	testCases := []struct {
		reqBody         string
		expectedResBody string
	}{
		{`{"totalIncome":60000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, `{"tax":0.0}`},
		{`{"totalIncome":500000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, `{"tax":29000.0}`},
		{`{"totalIncome":560000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, `{"tax":35000.0}`},
		{`{"totalIncome":1060000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, `{"tax":110000.0}`},
		{`{"totalIncome":2060000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, `{"tax":310000.0}`},
		{`{"totalIncome":2060001.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, `{"tax":310000.35}`},
	}

	for _, tc := range testCases {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := New()
		// Assertions
		if assert.NoError(t, h.TaxCalculateHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, tc.expectedResBody, rec.Body.String())
		}

	}
}
