package discount

import (
	"fmt"
	"testing"
)

func TestSettleWithoutRules(t *testing.T) {
	total, discount, _ := New().Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, true)
	if total != 1000 || discount != 0 {
		t.Fatalf("没有规则时应该原价，得到 total=%d discount=%d", total, discount)
	}
}

func TestSingleCoupon(t *testing.T) {
	desk := New()
	desk.Add("c1", "all", "coupon", 100, 1)
	total, discount, _ := desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, false)
	if total != 900 || discount != 100 {
		t.Fatalf("一张券应该扣 100，得到 total=%d discount=%d", total, discount)
	}
}

func TestInactiveKindIsIgnored(t *testing.T) {
	desk := New()
	desk.Add("v1", "all", "vip", 100, 1)
	total, _, _ := desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, false)
	if total != 1000 {
		t.Fatalf("没开会员时不该扣会员折扣，得到 total=%d", total)
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New().Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}

func TestScopeCategoryOverridesAll(t *testing.T) {
	desk := New()
	desk.Add("g1", "all", "coupon", 200, 1)
	desk.Add("s1", "weapon", "coupon", 100, 1)
	total, discount, _ := desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, false)
	if total != 900 || discount != 100 {
		t.Fatalf("类目规则应该覆盖全局规则，得到 total=%d discount=%d", total, discount)
	}
	total, _, _ = desk.Settle("o2", []Item{{Category: "armor", Price: 1000}}, true, false)
	if total != 800 {
		t.Fatalf("没有类目规则时全局规则应该生效，得到 total=%d", total)
	}
}

func TestSameKindKeepsBestRule(t *testing.T) {
	desk := New()
	desk.Add("v1", "all", "vip", 600, 1)
	desk.Add("v2", "all", "vip", 500, 5)
	total, discount, _ := desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, false, true)
	if total != 500 || discount != 500 {
		t.Fatalf("同层只该留优先级最高的一条，得到 total=%d discount=%d", total, discount)
	}
}

func TestSamePriorityKeepsSmallerID(t *testing.T) {
	desk := New()
	desk.Add("c2", "all", "coupon", 50, 1)
	desk.Add("c1", "all", "coupon", 100, 1)
	total, _, _ := desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, false)
	if total != 900 {
		t.Fatalf("同优先级应该按 id 升序取第一条，得到 total=%d", total)
	}
}

func TestDifferentKindsStack(t *testing.T) {
	desk := New()
	desk.Add("c1", "all", "coupon", 100, 1)
	desk.Add("v1", "all", "vip", 50, 2)
	total, discount, _ := desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, true)
	if total != 850 || discount != 150 {
		t.Fatalf("不同 kind 应该各留一条一起扣，得到 total=%d discount=%d", total, discount)
	}
}

func TestPriceFloorPerItemAndOrder(t *testing.T) {
	desk := New()
	desk.Add("c1", "all", "coupon", 600, 1)
	desk.Add("v1", "all", "vip", 500, 5)
	total, discount, _ := desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, true)
	if total != 0 || discount != 1000 {
		t.Fatalf("整单总额不该小于 0，得到 total=%d discount=%d", total, discount)
	}

	perItem := New()
	perItem.Add("c1", "all", "coupon", 300, 1)
	total, discount, _ = perItem.Settle("o1", []Item{
		{Category: "weapon", Price: 100},
		{Category: "armor", Price: 500},
	}, true, false)
	if total != 200 || discount != 400 {
		t.Fatalf("每件商品应该各自夹到 0，得到 total=%d discount=%d", total, discount)
	}
}

func TestDuplicateOrderSettlesOnce(t *testing.T) {
	desk := New()
	desk.Add("c1", "all", "coupon", 100, 1)
	total, discount, first := desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, false)
	if !first || total != 900 || discount != 100 {
		t.Fatalf("第一次结算应该生效，得到 total=%d discount=%d first=%v", total, discount, first)
	}
	scanned := desk.Scanned()
	total, discount, first = desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, false)
	if first {
		t.Fatal("同一订单号重复结算不该算本次结算")
	}
	if total != 900 || discount != 100 {
		t.Fatalf("重复结算应该返回上一次的结果，得到 total=%d discount=%d", total, discount)
	}
	if desk.Scanned() != scanned {
		t.Fatalf("重复结算不该再扫描规则，scanned 从 %d 变成 %d", scanned, desk.Scanned())
	}
}

func TestAddInvalidatesAndRecomputes(t *testing.T) {
	desk := New()
	desk.Add("c1", "all", "coupon", 100, 1)
	desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, false)
	desk.Add("c2", "all", "coupon", 50, 2)
	total, _, first := desk.Settle("o2", []Item{{Category: "weapon", Price: 1000}}, true, false)
	if !first || total != 850 {
		t.Fatalf("新登记的规则应该参与结算，得到 total=%d first=%v", total, first)
	}
	total, _, first = desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, false)
	if !first || total != 850 {
		t.Fatalf("Add 之后旧结果应该作废重算，得到 total=%d first=%v", total, first)
	}
}

func TestRemoveInvalidatesAndRecomputes(t *testing.T) {
	desk := New()
	desk.Add("c1", "all", "coupon", 100, 1)
	desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, false)
	if !desk.Remove("c1") {
		t.Fatal("c1 原本存在，Remove 应该返回 true")
	}
	if desk.Remove("c1") {
		t.Fatal("c1 已经下架，Remove 应该返回 false")
	}
	total, _, first := desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, false)
	if !first || total != 1000 {
		t.Fatalf("Remove 之后旧结果应该作废重算，得到 total=%d first=%v", total, first)
	}
}

func TestRemoveWinnerResurfacesNextRule(t *testing.T) {
	desk := New()
	desk.Add("v1", "all", "vip", 600, 1)
	desk.Add("v2", "all", "vip", 500, 5)
	total, _, _ := desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, false, true)
	if total != 500 {
		t.Fatalf("取优后应该只留 v2，得到 total=%d", total)
	}
	desk.Remove("v2")
	total, _, first := desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, false, true)
	if !first || total != 400 {
		t.Fatalf("移除胜出规则后应按最新规则表重算，v1 重新生效，得到 total=%d first=%v", total, first)
	}
}

func TestAddDuplicateIDRejected(t *testing.T) {
	desk := New()
	if !desk.Add("c1", "all", "coupon", 100, 1) {
		t.Fatal("首次登记应该成功")
	}
	if desk.Add("c1", "all", "coupon", 200, 2) {
		t.Fatal("同一个 id 重复登记应该返回 false")
	}
	total, _, _ := desk.Settle("o1", []Item{{Category: "weapon", Price: 1000}}, true, false)
	if total != 900 {
		t.Fatalf("重复登记不该影响已有规则，得到 total=%d", total)
	}
}

func TestScannedDoesNotScaleWithRules(t *testing.T) {
	desk := New()
	for index := 0; index < 100; index++ {
		desk.Add(fmt.Sprintf("r-%d", index), "all", "coupon", 1, index)
	}
	for index := 0; index < 100; index++ {
		desk.Settle(fmt.Sprintf("o-%d", index), []Item{{Category: "weapon", Price: 100000}}, true, false)
	}
	if desk.Scanned() >= 100*100 {
		t.Fatalf("Scanned 不该随规则数乘结算次数增长，得到 %d", desk.Scanned())
	}
}
