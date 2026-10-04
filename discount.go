// Package discount 是折扣结算：优惠券与会员各是一层，同层只取优先级最高的一条，结算价不为负。
// 缺陷：同层不互斥（多条同时扣）、价格没夹到 0、重复登记照收、每次结算全扫折扣表。
package discount

type rule struct {
	id       string
	kind     string
	amount   int
	priority int
}

// Desk 是折扣台账。
type Desk struct {
	rules   []rule
	scanned int
}

// New 建空台账。
func New() *Desk { return &Desk{} }

// Add 登记一条折扣规则，返回是否成功。缺陷：同一个 id 重复登记也照收。
func (d *Desk) Add(id, kind string, amount, priority int) bool {
	d.rules = append(d.rules, rule{id: id, kind: kind, amount: amount, priority: priority})
	return true
}

// Price 计算最终价与总折扣。缺陷：同层多条一起扣、价格可以扣成负数、全扫折扣表。
func (d *Desk) Price(base int, coupon, vip bool) (int, int) {
	d.scanned += len(d.rules)
	price := base
	for _, item := range d.rules {
		if item.kind == "coupon" && !coupon {
			continue
		}
		if item.kind == "vip" && !vip {
			continue
		}
		price -= item.amount
	}
	return price, base - price
}

// Scanned 是累计扫描的规则数（规模观测）。
func (d *Desk) Scanned() int { return d.scanned }
