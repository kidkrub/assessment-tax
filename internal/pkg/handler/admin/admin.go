package admin

import (
	"database/sql"
	"net/http"

	"github.com/kidkrub/assessment-tax/internal/pkg/db"
	"github.com/labstack/echo/v4"
)

type SetDuctionRequestObject struct {
	Amount float64 `json:"amount"`
}

type SetDuctionResponseObject struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type handler struct {
	db *sql.DB
}

func New(db *sql.DB) *handler {
	return &handler{db}
}

func (h handler) SetDeductionValueHandler(c echo.Context) error {
	dType := c.Param("type")
	setDuctionRequestObject := SetDuctionRequestObject{}
	if err := c.Bind(&setDuctionRequestObject); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}
	if setDuctionRequestObject.Amount > 100000 || setDuctionRequestObject.Amount < 10000 {
		return echo.NewHTTPError(http.StatusBadRequest, "amount must between 10,000 - 100,000")
	}
	value := db.SetDeductionValue(h.db, dType, setDuctionRequestObject.Amount)
	return c.JSON(http.StatusOK, SetDuctionResponseObject{value})
}
