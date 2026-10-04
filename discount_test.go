package discount

import "testing"

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
