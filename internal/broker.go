package internal

import (
	"log"
	"sync"
	"time"

	"github.com/HaoxuanXu/MATradingBot/config"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/shopspring/decimal"
)

type AlpacaBroker struct {
	client  alpaca.Client
	account *alpaca.Account
	Clock   alpaca.Clock
}

var lock = &sync.Mutex{}

var (
	generatedBroker *AlpacaBroker
)

func GetBroker(accountType, serverType string) *AlpacaBroker {
	lock.Lock()
	defer lock.Unlock()

	if generatedBroker == nil {
		generatedBroker = &AlpacaBroker{}
		generatedBroker.initialize(accountType, serverType)
	}
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
	return updatedOrder, success
}

func (broker *AlpacaBroker) SubmitOrder(qty float64, symbol, side, orderType, timeInForce string) *alpaca.Order {
	quantity := decimal.NewFromFloat(qty)
	order, _ := broker.client.PlaceOrder(
		alpaca.PlaceOrderRequest{
			AssetKey:    &symbol,
			AccountID:   broker.account.ID,
			Qty:         &quantity,
			Side:        alpaca.Side(side),
			Type:        alpaca.OrderType(orderType),
			TimeInForce: alpaca.TimeInForce(timeInForce),
		},
	)

	finalOrder, _ := broker.MonitorOrder(order)
	return finalOrder

}

func (broker *AlpacaBroker) ListPositions() []alpaca.Position {
	positions, err := broker.client.ListPositions()
	if err != nil {
		log.Panic(err)
	}
	return positions
}
