package portfolio

import (
	"fmt"
	"log"
	"time"
)

type Portfolio struct {
	balance   float64
	positions map[string]float64
	history   []Transaction
}

type Transaction struct {
	Timestamp time.Time
	Type      string
	Symbol    string
	Amount    float64
	Price     float64
	Total     float64
}

func NewPortfolio(initialBalance float64) *Portfolio {
	return &Portfolio{
		balance:   initialBalance,
		positions: make(map[string]float64),
		history:   make([]Transaction, 0),
	}
}

func (p *Portfolio) GetBalance() float64 {
	return p.balance
}

func (p *Portfolio) GetPosition(symbol string) float64 {
	return p.positions[symbol]
}

func (p *Portfolio) GetPositions() map[string]float64 {
	positions := make(map[string]float64)
	for symbol, amount := range p.positions {
		positions[symbol] = amount
	}
	return positions
}

func (p *Portfolio) GetHistory() []Transaction {
	return p.history
}

func (p *Portfolio) Buy(symbol string, dollarAmount float64, price float64) error {
	if dollarAmount > p.balance {
		return fmt.Errorf("insufficient balance: have %.2f, need %.2f", p.balance, dollarAmount)
	}

	quantity := dollarAmount / price
	p.balance -= dollarAmount
	p.positions[symbol] += quantity

	transaction := Transaction{
		Timestamp: time.Now(),
		Type:      "BUY",
		Symbol:    symbol,
		Amount:    quantity,
		Price:     price,
		Total:     dollarAmount,
	}
	p.history = append(p.history, transaction)

	log.Printf("BUY: %.6f %s at $%.2f (Total: $%.2f)", quantity, symbol, price, dollarAmount)
	return nil
}

func (p *Portfolio) Sell(symbol string, quantity float64, price float64) error {
	if position, exists := p.positions[symbol]; !exists || position < quantity {
		return fmt.Errorf("insufficient %s position: have %.6f, trying to sell %.6f", symbol, p.positions[symbol], quantity)
	}

	dollarAmount := quantity * price
	p.positions[symbol] -= quantity
	p.balance += dollarAmount

	if p.positions[symbol] <= 0 {
		delete(p.positions, symbol)
	}

	transaction := Transaction{
		Timestamp: time.Now(),
		Type:      "SELL",
		Symbol:    symbol,
		Amount:    quantity,
		Price:     price,
		Total:     dollarAmount,
	}
	p.history = append(p.history, transaction)

	log.Printf("SELL: %.6f %s at $%.2f (Total: $%.2f)", quantity, symbol, price, dollarAmount)
	return nil
}

func (p *Portfolio) GetTotalValue(currentPrices map[string]float64) float64 {
	totalValue := p.balance

	for symbol, quantity := range p.positions {
		if price, exists := currentPrices[symbol]; exists {
			totalValue += quantity * price
		}
	}

	return totalValue
}

func (p *Portfolio) PrintSummary(currentPrices map[string]float64) {
	fmt.Println("\n=== Portfolio Summary ===")
	fmt.Printf("Cash Balance: $%.2f\n", p.balance)

	if len(p.positions) > 0 {
		fmt.Println("Holdings:")
		for symbol, quantity := range p.positions {
			if price, exists := currentPrices[symbol]; exists {
				value := quantity * price
				fmt.Printf("  %s: %.6f (Value: $%.2f at $%.2f)\n", symbol, quantity, value, price)
			} else {
				fmt.Printf("  %s: %.6f (Price unknown)\n", symbol, quantity)
			}
		}
	}

	totalValue := p.GetTotalValue(currentPrices)
	fmt.Printf("Total Portfolio Value: $%.2f\n", totalValue)
	fmt.Println("========================")
}

func (p *Portfolio) GetRecentTransactions(count int) []Transaction {
	start := len(p.history) - count
	if start < 0 {
		start = 0
	}
	return p.history[start:]
}
