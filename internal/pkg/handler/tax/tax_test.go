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
		inputData TaxRequestObject
		expected  float64
	}{
		{TaxRequestObject{60000.0, 0.0, nil}, 0},
		{TaxRequestObject{500000.0, 0.0, nil}, 29000.0},
		{TaxRequestObject{560000.0, 0.0, nil}, 35000.0},
		{TaxRequestObject{1060000.0, 0.0, nil}, 110000.0},
		{TaxRequestObject{2060000.0, 0.0, nil}, 310000.0},
		{TaxRequestObject{2060001.0, 0.0, nil}, 310000.35},
		{TaxRequestObject{150000.0, 1000.0, nil}, -1000.0},
		{TaxRequestObject{500000.0, 25000.0, nil}, 4000.0},
		{TaxRequestObject{500000.0, 29000.0, nil}, 0.0},
		{TaxRequestObject{500000.0, 30000.0, nil}, -1000.0},
		{TaxRequestObject{500000.0, 0.0, []Allowance{{"donation", 200000.0}}}, 19000.0},
		{TaxRequestObject{500000.0, 0.0, []Allowance{{"donation", 100000.0}}}, 19000.0},
		{TaxRequestObject{500000.0, 0.0, []Allowance{{"donation", 50000.0}}}, 24000.0},
	}

	// Act & Assert
	for _, tc := range testCases {
		// Act
		actualTax := taxCalculate(tc.inputData)
		// Assert
		assert.Equal(t, tc.expected, actualTax, "tax calculation is incorrect for %.2f case", tc.inputData.TotalIncome)
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
		{`{"totalIncome":150000.0,"wht":1000.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, `{"tax":0.0,"taxRefund":1000.0}`},
		{`{"totalIncome":500000.0,"wht":25000.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, `{"tax":4000.0}`},
		{`{"totalIncome":500000.0,"wht":29000.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, `{"tax":0.0}`},
		{`{"totalIncome":500000.0,"wht":30000.0,"allowances":[{"allowanceType":"donation","amount":0.0}]}`, `{"tax":0.0,"taxRefund":1000.0}`},
		{`{"totalIncome":500000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":200000.0}]}`, `{"tax":19000.0}`},
		{`{"totalIncome":500000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":100000.0}]}`, `{"tax":19000.0}`},
		{`{"totalIncome":500000.0,"wht":0.0,"allowances":[{"allowanceType":"donation","amount":50000.0}]}`, `{"tax":24000.0}`},
	}

	for _, tc := range testCases {
		// Act
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
