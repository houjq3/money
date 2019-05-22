# go money [![Build Status](https://travis-ci.com/exfly/money.svg?branch=master)](https://travis-ci.com/exfly/money)

simple decimal money

```go
type Money struct {
	Value    bsond.Decimal128 `bson:"value"`
	Currency string           `bson:"currency"`
}
```
