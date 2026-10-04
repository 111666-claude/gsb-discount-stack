// Command discount-stack 跑折扣结算样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/discount"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("discount-stack", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "layer", "layer / floor / dup / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "layer":
		desk := discount.New()
		desk.Add("v1", "vip", 600, 1)
		desk.Add("v2", "vip", 500, 5)
		price, _ := desk.Price(1000, false, true)
		fmt.Fprintf(stdout, "price=%d\n", price)
	case "floor":
		desk := discount.New()
		desk.Add("c1", "coupon", 600, 1)
		desk.Add("v1", "vip", 500, 5)
		price, _ := desk.Price(1000, true, true)
		fmt.Fprintf(stdout, "price=%d\n", price)
	case "dup":
		desk := discount.New()
		desk.Add("c1", "coupon", 100, 1)
		second := desk.Add("c1", "coupon", 200, 2)
		fmt.Fprintf(stdout, "second=%v\n", second)
	case "scan":
		desk := discount.New()
		for index := 0; index < 3000; index++ {
			desk.Add(fmt.Sprintf("r-%d", index), "coupon", 1, index)
		}
		for index := 0; index < 3000; index++ {
			desk.Price(100000, true, false)
		}
		fmt.Fprintf(stdout, "scanned=%d\n", desk.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample layer|floor|dup|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
