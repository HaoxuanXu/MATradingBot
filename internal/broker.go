package internal

import (
	"log"
	"time"

	"github.com/HaoxuanXu/MATradingBot/config"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/shopspring/decimal"
)

type AlpacaBroker struct {
	client         alpaca.Client
	account        *alpaca.Account
	Clock          alpaca.Clock
	FilledQuantity float64
}

func GetBroker(accountType, serverType string) *AlpacaBroker {
	generatedBroker := &AlpacaBroker{}
	generatedBroker.initialize(accountType, serverType)

	return generatedBroker
}

func (broker *AlpacaBroker) initialize(accountType, serverType string) {
	cred := config.GetCredentials(accountType, serverType)
	broker.client = alpaca.NewClient(
		alpaca.ClientOpts{
			ApiKey:    cred.API_KEY,
			ApiSecret: cred.API_SECRET,
			BaseURL:   cred.BASE_URL,
		},
	)
	broker.account, _ = broker.client.GetAccount()
	clock, _ := broker.client.GetClock()
	broker.Clock = *clock
}

func (broker *AlpacaBroker) refreshOrderStatus(orderID string) (string, *alpaca.Order) {
	newOrder, _ := broker.client.GetOrder(orderID)
	orderStatus := newOrder.Status

	return orderStatus, newOrder
}

func (broker *AlpacaBroker) MonitorOrder(order *alpaca.Order) (*alpaca.Order, bool) {
	success := false
	orderID := order.ID
	status, updatedOrder := broker.refreshOrderStatus(orderID)
	if order.Type == alpaca.Market {
		for !success {
			switch status {
			case "new", "accepted", "partially_filled":
				time.Sleep(time.Second)
				status, updatedOrder = broker.refreshOrderStatus(orderID)
			case "filled":
				success = true
			case "done_for_day", "canceled", "expired", "replaced":
				success = false
			default:
				time.Sleep(time.Second)
				status, updatedOrder = broker.refreshOrderStatus(orderID)
			}
		}
	}

	return updatedOrder, success
}

func (broker *AlpacaBroker) SubmitMarketOrder(qty float64, symbol, side, timeInForce string) *alpaca.Order {
	quantity := decimal.NewFromFloat(qty)
	order, _ := broker.client.PlaceOrder(
		alpaca.PlaceOrderRequest{
			AssetKey:    &symbol,
			AccountID:   broker.account.ID,
			Qty:         &quantity,
			Side:        alpaca.Side(side),
			Type:        alpaca.OrderType(alpaca.Market),
			TimeInForce: alpaca.TimeInForce(timeInForce),
		},
	)

	finalOrder, _ := broker.MonitorOrder(order)
	broker.FilledQuantity = finalOrder.FilledQty.InexactFloat64()
	return finalOrder

}

func (broker *AlpacaBroker) SubmitTrailingStopOrder(qty, trail_percent float64, symbol, side, timeInForce string) *alpaca.Order {
	quantity := decimal.NewFromFloat(qty)
	trail := decimal.NewFromFloat(trail_percent)
	order, _ := broker.client.PlaceOrder(
		alpaca.PlaceOrderRequest{
			AssetKey:     &symbol,
			AccountID:    broker.account.ID,
			Qty:          &quantity,
			Side:         alpaca.Side(side),
			Type:         alpaca.OrderType(alpaca.TrailingStop),
			TrailPercent: &trail,
		},
	)
	for order.Status != "new" {
		_, order = broker.refreshOrderStatus(order.ID)
		time.Sleep(5 * time.Second)
	}
	return order
}

func (broker *AlpacaBroker) ListPositions() []alpaca.Position {
	positions, err := broker.client.ListPositions()
	if err != nil {
		log.Panic(err)
	}
	return positions
}

func (broker *AlpacaBroker) GetPosition(symbol string) *alpaca.Position {
	position, err := broker.client.GetPosition(symbol)
	if err != nil {
		log.Println(err)
		return nil
	}
	return position
}
