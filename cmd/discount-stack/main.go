// Command discount-stack 跑折扣结算样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/discount"
)

func weapon(price int) []discount.Item {
	return []discount.Item{{Category: "weapon", Price: price}}
}

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("discount-stack", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "layer", "layer / scope / floor / stale / dup / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "layer":
		desk := discount.New()
		desk.Add("v1", "all", "vip", 600, 1)
		desk.Add("v2", "all", "vip", 500, 5)
		total, _, _ := desk.Settle("o1", weapon(1000), false, true)
		fmt.Fprintf(stdout, "total=%d\n", total)
	case "scope":
		desk := discount.New()
		desk.Add("g1", "all", "coupon", 200, 1)
		desk.Add("s1", "weapon", "coupon", 100, 1)
		total, _, _ := desk.Settle("o1", weapon(1000), true, false)
		fmt.Fprintf(stdout, "total=%d\n", total)
	case "floor":
		desk := discount.New()
		desk.Add("c1", "all", "coupon", 600, 1)
		desk.Add("v1", "all", "vip", 500, 5)
		total, _, _ := desk.Settle("o1", weapon(1000), true, true)
		fmt.Fprintf(stdout, "total=%d\n", total)
	case "stale":
		desk := discount.New()
		desk.Add("c1", "all", "coupon", 100, 1)
		desk.Settle("o1", weapon(1000), true, false)
		desk.Add("c2", "all", "coupon", 50, 2)
		total, _, _ := desk.Settle("o2", weapon(1000), true, false)
		fmt.Fprintf(stdout, "total=%d\n", total)
	case "dup":
		desk := discount.New()
		desk.Add("c1", "all", "coupon", 100, 1)
		desk.Settle("o1", weapon(1000), true, false)
		_, _, first := desk.Settle("o1", weapon(1000), true, false)
		fmt.Fprintf(stdout, "second_first=%v\n", first)
	case "scan":
		desk := discount.New()
		for index := 0; index < 3000; index++ {
			desk.Add(fmt.Sprintf("r-%d", index), "all", "coupon", 1, index)
		}
		for index := 0; index < 3000; index++ {
			desk.Settle(fmt.Sprintf("o-%d", index), weapon(100000), true, false)
		}
		fmt.Fprintf(stdout, "scanned=%d\n", desk.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample layer|scope|floor|stale|dup|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
