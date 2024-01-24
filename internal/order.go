package internal

import "github.com/sandmannmax/go-trade/internal/money"

type order struct {
	volume float64
	isBid  bool
}

func (o order) IsFilled() bool {
	return o.volume == 0
}

func (o order) Type() string {
	if o.isBid {
		return "BID"
	}
	return "ASK"
}

func getMinMax(o1 *order, o2 *order) (*order, *order) {
	if o2.volume < o1.volume {
		return o2, o1
	}

	return o1, o2
}

type limitOrder struct {
	order
	price  money.Money
	orders []*order
}

type limitOrderListElement struct {
	order *limitOrder
	next  *limitOrderListElement
	prev  *limitOrderListElement
}

type limitOrderList struct {
	first     *limitOrderListElement
	last      *limitOrderListElement
	isBidList bool
}

func (list *limitOrderList) Add(price money.Money, order *limitOrder) {
	newElem := &limitOrderListElement{order: order}
	if list.first == nil {
		list.first = newElem
		list.last = newElem
		return
	}

	elem := list.first

	for {
		if elem == nil {
			newElem.prev = list.last
			newElem.prev.next = newElem
			list.last = newElem
			return
		}

		if elem.order.price.IsGreaterAs(price) {
			newElem.prev = elem.prev
			if newElem.prev != nil {
				newElem.prev.next = newElem
			} else {
				list.first = newElem
			}
			newElem.next = elem
			newElem.next.prev = newElem
			return
		}

		if elem.order.price == price {
			panic("ELEMENT ALREADY THERE")
		}

		elem = elem.next
	}

}

func (list limitOrderList) Find(price money.Money) (*limitOrder, bool) {
	elem := list.first

	for {
		if elem == nil || elem.order.price.IsGreaterAs(price) {
			return nil, false
		}

		if elem.order.price == price {
			return elem.order, true
		}

		elem = elem.next
	}
}

func (list limitOrderList) Delete(price money.Money) {
	var (
		elem = list.first
	)

	for {
		if elem == nil || elem.order.price.IsGreaterAs(price) {
			return
		}

		if elem.order.price == price {
			if elem.next != nil {
				elem.next.prev = elem.prev
			} else {
				list.last = elem.prev
			}
			if elem.prev != nil {
				elem.prev.next = elem.next
			} else {
				list.first = elem.next
			}
			return
		}

		elem = elem.next
	}
}

func (list limitOrderList) Iterator() []*limitOrder {
	orders := []*limitOrder{}
	elem := list.first

	if list.isBidList {
		elem = list.last
	}

	for {
		if elem == nil {
			return orders
		}

		orders = append(orders, elem.order)

		if list.isBidList {
			elem = elem.prev
		} else {
			elem = elem.next
		}
	}
}
