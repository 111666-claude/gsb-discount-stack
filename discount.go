// Package discount 是折扣结算：规则按作用域与类型分层，结算时按优先级依次作用在商品价格上。
// 缺陷：作用域不生效（类目规则与全局规则一起扣）、同层不取优先级最高的那条、
// 价格没有夹到 0、规则上下线后仍用第一次结算时冻结的表、同一订单重复结算、每次结算全扫规则表。
package discount

type rule struct {
	id       string
	scope    string
	kind     string
	amount   int
	priority int
}

// Item 是订单里的一件商品。
type Item struct {
	Category string
	Price    int
}

// Desk 是折扣台账。
type Desk struct {
	rules   []rule
	frozen  []rule
	seen    map[string]bool
	scanned int
}

// New 建空台账。
func New() *Desk { return &Desk{seen: map[string]bool{}} }

// Add 登记一条折扣规则，返回是否成功。缺陷：同一个 id 重复登记也照收。
func (d *Desk) Add(id, scope, kind string, amount, priority int) bool {
	d.rules = append(d.rules, rule{id: id, scope: scope, kind: kind, amount: amount, priority: priority})
	return true
}

// Remove 下架一条折扣规则，返回它原本是否存在。
func (d *Desk) Remove(id string) bool {
	for index, item := range d.rules {
		if item.id == id {
			d.rules = append(d.rules[:index], d.rules[index+1:]...)
			return true
		}
	}
	return false
}

// Settle 结算整单，返回最终总额、总折扣，以及这次调用是否真的完成了一次结算。
// 缺陷：作用域不生效、同层不取优、价格不夹 0、规则表只冻结一次、没有幂等、全扫规则表。
func (d *Desk) Settle(orderID string, items []Item, coupon, vip bool) (int, int, bool) {
	d.scanned += len(d.rules)
	if d.frozen == nil {
		d.frozen = append([]rule{}, d.rules...)
	}
	base := 0
	total := 0
	for _, item := range items {
		base += item.Price
		price := item.Price
		for _, entry := range d.frozen {
			if entry.kind == "coupon" && !coupon {
				continue
			}
			if entry.kind == "vip" && !vip {
				continue
			}
			price -= entry.amount
		}
		total += price
	}
	return total, base - total, true
}

// Scanned 是累计扫描的规则数（规模观测）。
func (d *Desk) Scanned() int { return d.scanned }
