package money

import "fmt"

type amount int64

type Money struct {
	amount
}

func New(amount amount) Money {
	return Money{amount: amount}
}

func NewFromFloat(value float64) Money {
	return Money{amount: amount(value * 100)}
}

func (m Money) Display() string {
	return fmt.Sprintf("%d.%02d", m.amount/100, m.amount%100)
}

func (m1 Money) IsGreaterAs(m2 Money) bool {
	return m1.amount > m2.amount
}

func (m1 Money) IsLessAs(m2 Money) bool {
	return m1.amount > m2.amount
}
