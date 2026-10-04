package discount

import "testing"

func TestPriceWithoutRules(t *testing.T) {
	price, discount := New().Price(1000, true, true)
	if price != 1000 || discount != 0 {
		t.Fatalf("没有规则时应该原价，得到 price=%d discount=%d", price, discount)
	}
}

func TestSingleCoupon(t *testing.T) {
	desk := New()
	desk.Add("c1", "coupon", 100, 1)
	price, discount := desk.Price(1000, true, false)
	if price != 900 || discount != 100 {
		t.Fatalf("一张券应该扣 100，得到 price=%d discount=%d", price, discount)
	}
}

func TestInactiveLayerIsIgnored(t *testing.T) {
	desk := New()
	desk.Add("v1", "vip", 100, 1)
	price, _ := desk.Price(1000, false, false)
	if price != 1000 {
		t.Fatalf("没开会员时不该扣会员折扣，得到 price=%d", price)
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New().Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}
