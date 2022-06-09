package api

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/HaoxuanXu/MATradingBot/config"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/shopspring/decimal"
)

type AlpacaBroker struct {
	Client  alpaca.Client
	Account *alpaca.Account
	Clock   *alpaca.Clock
	Cash    float64
}

var lock sync.Mutex

func GetBroker(accountType, serverType string) *AlpacaBroker {

	generatedBroker := &AlpacaBroker{}
	generatedBroker.initialize(accountType, serverType)

	return generatedBroker
}

func (broker *AlpacaBroker) initialize(accountType, serverType string) {
	cred := config.GetCredentials(accountType, serverType)
	broker.Client = alpaca.NewClient(
		alpaca.ClientOpts{
			ApiKey:    cred.API_KEY,
			ApiSecret: cred.API_SECRET,
			BaseURL:   cred.BASE_URL,
		},
	)
	broker.Account, _ = broker.Client.GetAccount()
	clock, _ := broker.Client.GetClock()
	broker.Clock = clock
	broker.Cash = broker.Account.Equity.Abs().InexactFloat64()
}

func (broker *AlpacaBroker) GetClock() (*alpaca.Clock, error) {
	return broker.Client.GetClock()
}

func (broker *AlpacaBroker) RefreshOrderStatus(orderID string) (*alpaca.Order, error) {
	newOrder, err := broker.Client.GetOrder(orderID)

	return newOrder, err
}

func (broker *AlpacaBroker) MonitorOrder(order *alpaca.Order) *alpaca.Order {
	finished := false
	orderID := order.ID
	order, _ = broker.RefreshOrderStatus(orderID)
	if order.Type == alpaca.Market {
		for !finished {
			switch order.Status {
			case "new", "accepted", "partially_filled":
				time.Sleep(time.Second)
				order, _ = broker.RefreshOrderStatus(orderID)
			case "filled":
				finished = true
			case "done_for_day", "canceled", "expired", "replaced":
				finished = true
			default:
				time.Sleep(time.Second)
				order, _ = broker.RefreshOrderStatus(orderID)
			}
		}
	} else if order.Type == alpaca.TrailingStop {
		for !finished {
			switch order.Status {
			case "accepted", "partially_filled":
				time.Sleep(time.Second)
				order, _ = broker.RefreshOrderStatus(orderID)
			case "new":
				finished = true
			case "done_for_day", "canceled", "expired", "replaced":
				finished = true
			default:
				time.Sleep(time.Second)
				order, _ = broker.RefreshOrderStatus(orderID)
			}
		}
	}

	return order
}

func (broker *AlpacaBroker) SubmitMarketOrder(qty float64, symbol, side, timeInForce string) *alpaca.Order {
	defer lock.Unlock()
	lock.Lock()
	quantity := decimal.NewFromFloat(qty)
	order, err := broker.Client.PlaceOrder(
		alpaca.PlaceOrderRequest{
			AssetKey:    &symbol,
			AccountID:   broker.Account.ID,
			Qty:         &quantity,
			Side:        alpaca.Side(side),
			Type:        alpaca.OrderType(alpaca.Market),
			TimeInForce: alpaca.TimeInForce(timeInForce),
		},
	)
	if err != nil {
		log.Println(err)
	}

	finalOrder := broker.MonitorOrder(order)
	return finalOrder

}

func (broker *AlpacaBroker) SubmitTrailingStopOrder(qty, trail_price float64, symbol, side string) *alpaca.Order {
	defer lock.Unlock()
	lock.Lock()
	quantity := decimal.NewFromFloat(qty)
	trail := decimal.NewFromFloat(trail_price)
	order, err := broker.Client.PlaceOrder(
		alpaca.PlaceOrderRequest{
			AssetKey:    &symbol,
			AccountID:   broker.Account.ID,
			Qty:         &quantity,
			Side:        alpaca.Side(side),
			Type:        alpaca.OrderType(alpaca.TrailingStop),
			TrailPrice:  &trail,
			TimeInForce: alpaca.GTC,
		},
	)
	if err != nil {
		log.Println(err)
	}
	finalOrder := broker.MonitorOrder(order)
	return finalOrder
}

func (broker *AlpacaBroker) SubmitBracketOrder(qty, take_profit, take_loss float64, symbol, side string) *alpaca.Order {
	defer lock.Unlock()
	lock.Lock()
	quantity := decimal.NewFromFloat(qty)
	takeProfitLimitPrice := decimal.NewFromFloat(math.Floor(take_profit*100) / 100)
	takeLossLimitPrice := decimal.NewFromFloat(math.Floor(take_loss*100) / 100)
	order, err := broker.Client.PlaceOrder(
		alpaca.PlaceOrderRequest{
			AssetKey:   &symbol,
			AccountID:  broker.Account.ID,
			Qty:        &quantity,
			Side:       alpaca.Side(side),
			Type:       alpaca.Market,
			OrderClass: alpaca.Bracket,
			TakeProfit: &alpaca.TakeProfit{
				LimitPrice: &takeProfitLimitPrice,
			},
			StopLoss: &alpaca.StopLoss{
				LimitPrice: &takeLossLimitPrice,
				StopPrice:  &takeLossLimitPrice,
			},
			TimeInForce: alpaca.GTC,
		},
	)
	if err != nil {
		log.Printf("%s: %v", symbol, err)
	}
	finalOrder := broker.MonitorOrder(order)
	return finalOrder

}

func (broker *AlpacaBroker) ChangeOrderTrail(order *alpaca.Order, newTrail float64) *alpaca.Order {
	ordeID := order.ID
	newTrailDecimal := decimal.NewFromFloat(newTrail)
	order, _ = broker.Client.ReplaceOrder(
		ordeID,
		alpaca.ReplaceOrderRequest{
			Trail: &newTrailDecimal,
		},
	)
	finalOrder := broker.MonitorOrder(order)
	return finalOrder
}

func (broker *AlpacaBroker) RetrieveOrderIfExists(symbol, status, orderType string) (*alpaca.Order, error) {
	defer lock.Unlock()
	lock.Lock()
	limit := 10
	nested := false
	until := time.Now()
	orderList, err := broker.Client.ListOrders(&status, &until, &limit, &nested)

	for _, order := range orderList {
		if order.Symbol == symbol && order.Type == alpaca.OrderType(orderType) {
			return &order, err
		}
	}
	return nil, err
}

func (broker *AlpacaBroker) ListPositions() []alpaca.Position {
	positions, err := broker.Client.ListPositions()
	if err != nil {
		log.Fatal(err)
	}
	return positions
}

func (broker *AlpacaBroker) GetPosition(symbol string) (*alpaca.Position, error) {
	lock.Lock()
	defer lock.Unlock()
	position, err := broker.Client.GetPosition(symbol)
	if err != nil {
		return nil, err
	}
	return position, nil
}

func (broker *AlpacaBroker) ClosePosition(symbol string) error {
	lock.Lock()
	defer lock.Unlock()
	// Check if the position has already been closed or not
	position, err := broker.GetPosition(symbol)
	if position == nil && err != nil {
		log.Printf("The %s position has already been closed.\n", symbol)
		return err
	} else {
		err = broker.Client.ClosePosition(symbol)
		if err != nil {
			log.Printf("An error has occurred: %s", err)
			return err
		}
		// close open trailing order
		trailingStopOrder, _ := broker.RetrieveOrderIfExists(symbol, "new", "trailing_stop")
		err = broker.Client.CancelOrder(trailingStopOrder.ID)
		if err != nil {
			log.Printf("An error has occurred: %s", err)
			return err
		}
	}

	return nil
}
