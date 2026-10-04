// Package discount 是折扣结算：规则按作用域与类型分层，结算时按优先级依次作用在商品价格上。
// 类目规则覆盖同 kind 的全局规则；首次结算把规则表编译成生效表（同 kind 同作用域只留最优一条），
// 之后 Add 追加、Remove 触发重编译；每件商品与整单都夹到 0；同一订单号只结算一次。
package discount

import "sort"

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

// outcome 是一笔订单的结算结果。
type outcome struct {
	total    int
	discount int
}

// layer 用 (kind, scope) 定位同层规则。
type layer struct {
	kind  string
	scope string
}

// Desk 是折扣台账。
type Desk struct {
	rules   []rule            // 首次结算前登记的原始规则
	ids     map[string]bool   // 已登记规则的 id 集合
	byKind  map[string][]rule // 编译后的生效表：kind -> 规则（按优先级升序、同优先级按 id 升序）
	seen    map[string]outcome
	scanned int
}

// New 建空台账。
func New() *Desk {
	return &Desk{ids: map[string]bool{}, seen: map[string]outcome{}}
}

// Add 登记一条折扣规则，返回是否成功。同一个 id 重复登记返回 false，且不影响已有规则与结算结果。
// 生效表编译后新登记的直接进入生效表；登记成功后旧订单的结算结果一律作废。
func (d *Desk) Add(id, scope, kind string, amount, priority int) bool {
	if d.ids[id] {
		return false
	}
	d.ids[id] = true
	entry := rule{id: id, scope: scope, kind: kind, amount: amount, priority: priority}
	d.rules = append(d.rules, entry)
	if d.byKind != nil {
		d.byKind[kind] = insertRule(d.byKind[kind], entry)
	}
	d.seen = map[string]outcome{}
	return true
}

// Remove 下架一条折扣规则，返回它原本是否存在。下架后生效表作废，下次结算从最新规则表重编译，
// 旧结算结果一律作废。
func (d *Desk) Remove(id string) bool {
	for index, item := range d.rules {
		if item.id == id {
			d.rules = append(d.rules[:index], d.rules[index+1:]...)
			delete(d.ids, id)
			d.byKind = nil
			d.seen = map[string]outcome{}
			return true
		}
	}
	return false
}

// compile 在首次结算时把当前规则表编译成生效表：
// 同一个 (kind, scope) 只保留优先级最高的一条，同优先级按 id 升序取第一条。
func (d *Desk) compile() {
	d.byKind = map[string][]rule{}
	best := map[layer]rule{}
	for _, entry := range d.rules {
		key := layer{kind: entry.kind, scope: entry.scope}
		current, ok := best[key]
		if !ok || betterRule(entry, current) {
			best[key] = entry
		}
	}
	for _, entry := range best {
		d.byKind[entry.kind] = insertRule(d.byKind[entry.kind], entry)
	}
}

// betterRule 判断 left 是否比 right 更优：优先级更高，或同优先级时 id 更小。
func betterRule(left, right rule) bool {
	if left.priority != right.priority {
		return left.priority > right.priority
	}
	return left.id < right.id
}

// insertRule 把规则按结算顺序（优先级升序、同优先级按 id 升序）插入有序列表。
func insertRule(list []rule, entry rule) []rule {
	at := sort.Search(len(list), func(index int) bool {
		item := list[index]
		if item.priority != entry.priority {
			return item.priority > entry.priority
		}
		return item.id > entry.id
	})
	list = append(list, rule{})
	copy(list[at+1:], list[at:])
	list[at] = entry
	return list
}

// Settle 结算整单，返回最终总额、总折扣，以及这次调用是否真的完成了一次结算。
// 同一个订单号只结算一次，重复调用返回上一次的结果且不算本次结算。
func (d *Desk) Settle(orderID string, items []Item, coupon, vip bool) (int, int, bool) {
	if result, ok := d.seen[orderID]; ok {
		return result.total, result.discount, false
	}
	if d.byKind == nil {
		d.compile()
	}
	base := 0
	total := 0
	for _, item := range items {
		base += item.Price
		total += d.settleItem(item, coupon, vip)
	}
	if total < 0 {
		total = 0
	}
	d.seen[orderID] = outcome{total: total, discount: base - total}
	return total, base - total, true
}

// settleItem 结算一件商品：类目规则覆盖同 kind 的全局规则，生效规则按优先级升序依次扣减，价格夹到 0。
func (d *Desk) settleItem(item Item, coupon, vip bool) int {
	price := item.Price
	for kind, list := range d.byKind {
		if kind == "coupon" && !coupon {
			continue
		}
		if kind == "vip" && !vip {
			continue
		}
		scoped := false
		for _, entry := range list {
			if entry.scope == item.Category {
				scoped = true
				break
			}
		}
		for _, entry := range list {
			if entry.scope != "all" && entry.scope != item.Category {
				continue
			}
			if scoped && entry.scope == "all" {
				continue
			}
			d.scanned++
			price -= entry.amount
			if price < 0 {
				price = 0
			}
		}
	}
	return price
}

// Scanned 是累计扫描的规则数（规模观测）。只统计结算时实际触达的生效规则，
// 不随规则总数乘结算次数增长。
func (d *Desk) Scanned() int { return d.scanned }
