# discount-stack

折扣结算：优惠券与会员各是一层，同层只取优先级最高的一条，结算价不为负，只用 Go 标准库。

```
go run ./cmd/discount-stack --sample=layer
go run ./cmd/discount-stack --sample=floor
go run ./cmd/discount-stack --sample=dup
go run ./cmd/discount-stack --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **分层**：`coupon` 与 `vip` 各是一层；同一层里只有优先级最高的一条生效，同优先级按 id 升序取第一条。
- **结算顺序**：生效的规则按优先级升序依次从剩余价格里扣，先低优先级。
- **价格下限**：最终价不小于 0；总折扣等于原价减最终价。
- **幂等**：同一个 id 重复登记返回 false，且不影响已有规则。
- **不变量**：没有对应层开关的规则不参与结算；同一串登记与结算重复执行结果相同。
- **代价**：结算不许全扫折扣表，`Scanned` 不随规则数乘结算次数增长；10 万条规则每秒 10 万次结算。

## 输出契约（不改格式）

```
price=..
second=..
scanned=..
```

## 目录

```
discount.go            分层、同层取优、价格下限与幂等
discount_test.go       单元测试
cmd/discount-stack/    命令行入口
```
