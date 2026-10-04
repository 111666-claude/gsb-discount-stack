# discount-stack

折扣结算：规则按作用域与类型分层，结算时按优先级依次作用在商品价格上，只用 Go 标准库。

```
go run ./cmd/discount-stack --sample=layer
go run ./cmd/discount-stack --sample=scope
go run ./cmd/discount-stack --sample=floor
go run ./cmd/discount-stack --sample=stale
go run ./cmd/discount-stack --sample=dup
go run ./cmd/discount-stack --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **作用域覆盖层**：规则的作用域是 `all` 或某个商品类目；对一件商品，同一个 `kind` 里若类目规则存在，
  全局规则就不参与（类目覆盖全局）。
- **同层取优**：每个 `kind` 在一件商品上只保留一条规则，取优先级最高的，同优先级按 id 升序取第一条。
- **结算顺序**：一件商品上生效的规则按优先级升序依次扣减，先低优先级。
- **价格下限**：每件商品与整单总额都不小于 0；总折扣等于原价总额减最终总额。
- **开关**：`coupon` 或 `vip` 为 false 时对应 kind 的规则整层不参与。
- **幂等与重算**：同一个订单号只结算一次，重复调用返回上一次的结果并说明不是本次结算；
  `Add` 与 `Remove` 之后旧结果一律作废，下次结算必须用最新的规则表。
- **代价**：结算不许全扫规则表，`Scanned` 不随规则数乘结算次数增长；10 万条规则每秒 10 万次结算。

## 输出契约（不改格式）

```
total=..
second_first=..
scanned=..
```

## 目录

```
discount.go            作用域覆盖、同层取优、价格下限与幂等
discount_test.go       单元测试
cmd/discount-stack/    命令行入口
```
