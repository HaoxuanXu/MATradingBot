package internal

import (
	"log"
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

func (broker *AlpacaBroker) refreshOrderStatus(orderID string) *alpaca.Order {
	newOrder, _ := broker.client.GetOrder(orderID)

	return newOrder
}

func (broker *AlpacaBroker) MonitorOrder(order *alpaca.Order) *alpaca.Order {
	finished := false
	orderID := order.ID
	order = broker.refreshOrderStatus(orderID)
	if order.Type == alpaca.Market {
		for !finished {
			switch order.Status {
			case "new", "accepted", "partially_filled":
				time.Sleep(time.Second)
				order = broker.refreshOrderStatus(orderID)
			case "filled":
				finished = true
			case "done_for_day", "canceled", "expired", "replaced":
				finished = true
			default:
				time.Sleep(time.Second)
				order = broker.refreshOrderStatus(orderID)
			}
		}
	} else if order.Type == alpaca.TrailingStop {
		for !finished {
			switch order.Status {
			case "accepted", "partially_filled":
				time.Sleep(time.Second)
				order = broker.refreshOrderStatus(orderID)
			case "new":
				finished = true
			case "done_for_day", "canceled", "expired", "replaced":
				finished = true
			default:
				time.Sleep(time.Second)
				order = broker.refreshOrderStatus(orderID)
			}
		}
	}

	return order
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

	finalOrder := broker.MonitorOrder(order)
	return finalOrder

}

func (broker *AlpacaBroker) SubmitTrailingStopOrder(qty, trail_percent float64, symbol, side string) *alpaca.Order {
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
			TimeInForce:  alpaca.Day,
		},
	)
	finalOrder := broker.MonitorOrder(order)
	return finalOrder
}

func (broker *AlpacaBroker) ChangeOrderTrail(order *alpaca.Order, newTrail float64) *alpaca.Order {
	ordeID := order.ID
	newTrailDecimal := decimal.NewFromFloat(newTrail)
	order, _ = broker.client.ReplaceOrder(
		ordeID,
		alpaca.ReplaceOrderRequest{
			Trail: &newTrailDecimal,
		},
	)
	finalOrder := broker.MonitorOrder(order)
	return finalOrder
}

func (broker *AlpacaBroker) ListPositions() []alpaca.Position {
	positions, err := broker.client.ListPositions()
	if err != nil {
		log.Panic(err)
	}
	return positions
}

func (broker *AlpacaBroker) GetPosition(symbol string) (*alpaca.Position, error) {
	position, err := broker.client.GetPosition(symbol)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return position, nil
}

func (broker *AlpacaBroker) ClosePosition(symbol string) error {
	// Check if the position has already been closed or not
	position, err := broker.GetPosition(symbol)
	if position == nil && err != nil {
		log.Printf("The %s position has already been closed.\n", symbol)
		return err
	}
	err = broker.client.ClosePosition(symbol)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
		return err
	}
	return nil
}
