package kiteconnect

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

// Holding is an individual holdings response.
type Holding struct {
	Tradingsymbol   string `json:"tradingsymbol"`
	Exchange        string `json:"exchange"`
	InstrumentToken uint32 `json:"instrument_token"`
	ISIN            string `json:"isin"`
	Product         string `json:"product"`

	Price              float64 `json:"price"`
	Quantity           int     `json:"quantity"`
	T1Quantity         int     `json:"t1_quantity"`
	RealisedQuantity   int     `json:"realised_quantity"`
	CollateralQuantity int     `json:"collateral_quantity"`
	CollateralType     string  `json:"collateral_type"`

	AveragePrice        float64 `json:"average_price"`
	LastPrice           float64 `json:"last_price"`
	ClosePrice          float64 `json:"close_price"`
	PnL                 float64 `json:"pnl"`
	DayChange           float64 `json:"day_change"`
	DayChangePercentage float64 `json:"day_change_percentage"`
}

func (h Holding) String() string {
	return fmt.Sprintf(
		"Tradingsymbol: %v\nExchange: %v\nInstrumentToken: %v\nISIN: %v\nProduct: %v\nPrice: %v\nQuantity: %v\nT1Quantity: %v\nRealisedQuantity: %v\nCollateralQuantity: %v\nCollateralType: %v\nAveragePrice: %v\nLastPrice: %v\nClosePrice: %v\nPnL: %v\nDayChange: %v\nDayChangePercentage: %v\n",
		h.Tradingsymbol, h.Exchange, h.InstrumentToken, h.ISIN, h.Product, h.Price, h.Quantity, h.T1Quantity, h.RealisedQuantity, h.CollateralQuantity, h.CollateralType, h.AveragePrice, h.LastPrice, h.ClosePrice, h.PnL, h.DayChange, h.DayChangePercentage,
	)
}

// Holdings is a list of holdings
type Holdings []Holding

func (h Holdings) String() string {
	var w = bytes.NewBufferString("")
	for _, holding := range h {
		fmt.Fprintf(w, "%v\n", holding)
	}
	return w.String()
}

// Position represents an individual position response.
type Position struct {
	Tradingsymbol   string `json:"tradingsymbol"`
	Exchange        string `json:"exchange"`
	InstrumentToken uint32 `json:"instrument_token"`
	Product         string `json:"product"`

	Quantity          int     `json:"quantity"`
	OvernightQuantity int     `json:"overnight_quantity"`
	Multiplier        float64 `json:"multiplier"`

	AveragePrice float64 `json:"average_price"`
	ClosePrice   float64 `json:"close_price"`
	LastPrice    float64 `json:"last_price"`
	Value        float64 `json:"value"`
	PnL          float64 `json:"pnl"`
	M2M          float64 `json:"m2m"`
	Unrealised   float64 `json:"unrealised"`
	Realised     float64 `json:"realised"`

	BuyQuantity int     `json:"buy_quantity"`
	BuyPrice    float64 `json:"buy_price"`
	BuyValue    float64 `json:"buy_value"`
	BuyM2MValue float64 `json:"buy_m2m"`

	SellQuantity int     `json:"sell_quantity"`
	SellPrice    float64 `json:"sell_price"`
	SellValue    float64 `json:"sell_value"`
	SellM2MValue float64 `json:"sell_m2m"`

	DayBuyQuantity int     `json:"day_buy_quantity"`
	DayBuyPrice    float64 `json:"day_buy_price"`
	DayBuyValue    float64 `json:"day_buy_value"`

	DaySellQuantity int     `json:"day_sell_quantity"`
	DaySellPrice    float64 `json:"day_sell_price"`
	DaySellValue    float64 `json:"day_sell_value"`
}

func (p Position) String() string {
	return fmt.Sprintf(
		"Tradingsymbol:%v\nExchange:%v\nInstrumentToken:%v\nProduct:%v\nQuantity:%v\nOvernightQuantity:%v\nMultiplier:%v\nAveragePrice:%v\nClosePrice:%v\nLastPrice:%v\nValue:%v\nPnL:%v\nM2M:%v\nUnrealised:%v\nRealised:%v\nBuyQuantity:%v\nBuyPrice:%v\nBuyValue:%v\nBuyM2MValue:%v\nSellQuantity:%v\nSellPrice:%v\nSellValue:%v\nSellM2MValue:%v\nDayBuyQuantity:%v\nDayBuyPrice:%v\nDayBuyValue:%v\nDaySellQuantity:%v\nDaySellPrice:%v\nDaySellValue:%v\n",
		p.Tradingsymbol, p.Exchange, p.InstrumentToken, p.Product, p.Quantity, p.OvernightQuantity, p.Multiplier, p.AveragePrice, p.ClosePrice, p.LastPrice, p.Value, p.PnL, p.M2M, p.Unrealised, p.Realised, p.BuyQuantity, p.BuyPrice, p.BuyValue, p.BuyM2MValue, p.SellQuantity, p.SellPrice, p.SellValue, p.SellM2MValue, p.DayBuyQuantity, p.DayBuyPrice, p.DayBuyValue, p.DaySellQuantity, p.DaySellPrice, p.DaySellValue,
	)
}

// Positions represents a list of net and day positions.
type Positions struct {
	Net []Position `json:"net"`
	Day []Position `json:"day"`
}

func (p Positions) String() string {
	var w = bytes.NewBufferString("")
	fmt.Fprintf(w, "Net:\n")
	for _, position := range p.Net {
		fmt.Fprintf(w, "%v\n", position)
	}
	fmt.Fprintf(w, "Day:\n")
	for _, position := range p.Net {
		fmt.Fprintf(w, "%v\n", position)
	}
	return w.String()
}

// ConvertPositionParams represents the input params for a position conversion.
type ConvertPositionParams struct {
	Exchange        string `url:"exchange"`
	TradingSymbol   string `url:"tradingsymbol"`
	OldProduct      string `url:"old_product"`
	NewProduct      string `url:"new_product"`
	PositionType    string `url:"position_type"`
	TransactionType string `url:"transaction_type"`
	Quantity        string `url:"quantity"`
}

// GetHoldings gets a list of holdings.
func (c *Client) GetHoldings() (Holdings, error) {
	var holdings Holdings
	err := c.doEnvelope(http.MethodGet, URIGetHoldings, nil, nil, &holdings)
	return holdings, err
}

// GetPositions gets user positions.
func (c *Client) GetPositions() (Positions, error) {
	var positions Positions
	err := c.doEnvelope(http.MethodGet, URIGetPositions, nil, nil, &positions)
	return positions, err
}

// ConvertPosition converts postion's product type.
func (c *Client) ConvertPosition(positionParams ConvertPositionParams) (bool, error) {
	var (
		b      bool
		err    error
		params url.Values
	)

	if params, err = query.Values(positionParams); err != nil {
		return false, NewError(InputError, fmt.Sprintf("Error decoding order params: %v", err), nil)
	}

	if err = c.doEnvelope(http.MethodPut, URIConvertPosition, params, nil, nil); err == nil {
		b = true
	}

	return b, err
}
