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
	Tax       float64    `json:"tax"`
	TaxRefund float64    `json:"taxRefund,omitempty"`
	TaxLevels []TaxLevel `json:"taxLevel"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
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
	tax, taxLevels := taxCalculate(taxRequestObject)
	res := TaxResponseObject{}
	res.TaxLevels = taxLevels
	if tax < 0 {
		res.Tax = 0
		res.TaxRefund = math.Abs(tax)
	} else {
		res.Tax = tax
	}
	return c.JSON(http.StatusOK, res)
}

func taxCalculate(inputData TaxRequestObject) (tax float64, taxLevelsObject []TaxLevel) {
	taxable := inputData.TotalIncome - 60000

	if len(inputData.Allowances) > 0 {
		donationAmount := 0.0
		for _, allowance := range inputData.Allowances {
			if allowance.AllowanceType == "donation" {
				donationAmount += allowance.Amount
			}
		}
		if donationAmount > 100000 {
			donationAmount = 100000
		}
		taxable -= donationAmount
	}

	taxLevels := []struct {
		level      string
		tierDiff   float64
		multiplier float64
	}{
		{"0-150,000", 150000, 0},
		{"150,001-500,000", 350000, 0.1},
		{"500,001-1,000,000", 500000, 0.15},
		{"1,000,001-2,000,000", 1000000, 0.2},
		{"2,000,001 ขึ้นไป", -1, 0.35},
	}
	for _, taxLevel := range taxLevels {
		if taxable < 0 {
			taxLevelsObject = []TaxLevel{
				{"0-150,000", 0.0},
				{"150,001-500,000", 0.0},
				{"500,001-1,000,000", 0.0},
				{"1,000,001-2,000,000", 0.0},
				{"2,000,001 ขึ้นไป", 0.0},
			}
			break
		}
		if taxable > taxLevel.tierDiff && taxLevel.tierDiff != -1 {
			tierTax := taxLevel.tierDiff * taxLevel.multiplier
			tax += tierTax
			taxable -= taxLevel.tierDiff
			taxLevelObject := TaxLevel{taxLevel.level, tierTax}
			taxLevelsObject = append(taxLevelsObject, taxLevelObject)
			continue
		}
		if taxable == 0 {
			taxLevelObject := TaxLevel{taxLevel.level, 0}
			taxLevelsObject = append(taxLevelsObject, taxLevelObject)
			continue
		}
		tierTax := taxable * taxLevel.multiplier
		tax += tierTax
		taxable = 0
		taxLevelObject := TaxLevel{taxLevel.level, tierTax}
		taxLevelsObject = append(taxLevelsObject, taxLevelObject)
	}
	return tax - inputData.Wht, taxLevelsObject
}
