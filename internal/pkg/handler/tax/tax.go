package tax

import (
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxRequestObject struct {
	TotalIncome float64     `json:"totalIncome"`
	Wht         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type TaxResponseObject struct {
	Tax       float64 `json:"tax"`
	TaxRefund float64 `json:"taxRefund,omitempty"`
}

type handler struct {
}

func New() *handler {
	return &handler{}
}

func (h handler) TaxCalculateHandler(c echo.Context) error {
	taxRequestObject := TaxRequestObject{}
	if err := c.Bind(&taxRequestObject); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}
	tax := taxCalculate(taxRequestObject)
	res := TaxResponseObject{}
	if tax < 0 {
		res.Tax = 0
		res.TaxRefund = math.Abs(tax)
	} else {
		res.Tax = tax
	}
	return c.JSON(http.StatusOK, res)
}

func taxCalculate(inputData TaxRequestObject) float64 {
	taxable := inputData.TotalIncome - 60000

	taxLevels := []struct {
		tierDiff   float64
		multiplier float64
	}{
		{150000, 0},
		{350000, 0.1},
		{500000, 0.15},
		{1000000, 0.2},
		{-1, 0.35},
	}
	tax := 0.0
	for _, taxLevel := range taxLevels {
		if taxable < 0 {
			break
		}
		if taxable > taxLevel.tierDiff && taxLevel.tierDiff != -1 {
			tax += taxLevel.tierDiff * taxLevel.multiplier
			taxable -= taxLevel.tierDiff
			continue
		}
		tax += taxable * taxLevel.multiplier
		taxable = 0
	}
	return tax - inputData.Wht
}
