package common

import (
	"errors"
	"strconv"
)

// PriceLevel is a common structure for bids and asks in the
// order book.
type PriceLevel struct {
	Price    string
	Quantity string
}

// Parse parses this PriceLevel's Price and Quantity and
// returns them both.  It also returns an error if either
// fails to parse.
func (p *PriceLevel) Parse() (float64, float64, error) {
	price, err := strconv.ParseFloat(p.Price, 64)
	if err != nil {
		return 0, 0, err
	}
	quantity, err := strconv.ParseFloat(p.Quantity, 64)
	if err != nil {
		return price, 0, err
	}
	return price, quantity, nil
}

type PriceLevelArray []string

// Parse parses this PriceLevelArray Price and Quantity and
// returns them both.  It also returns an error if either
// fails to parse.
func (p PriceLevelArray) Parse() (float64, float64, error) {
	if len(p) != 2 {
		return 0, 0, errors.New("empty price level array")
	}
	price, err := strconv.ParseFloat(p[0], 64)
	if err != nil {
		return 0, 0, err
	}
	quantity, err := strconv.ParseFloat(p[1], 64)
	if err != nil {
		return price, 0, err
	}
	return price, quantity, nil
}
