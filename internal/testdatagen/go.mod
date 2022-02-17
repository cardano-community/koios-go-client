module github.com/cardano-community/koios-go-client/internal/testdatagen

go 1.18

replace (
	github.com/cardano-community/koios-go-client v0.0.0 => ../..
	github.com/shopspring/decimal v1.3.1 => github.com/howijd/decimal v1.3.1
)

require (
	github.com/cardano-community/koios-go-client v0.0.0
	github.com/urfave/cli/v2 v2.3.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0-20190314233015-f79a8a8ca69d // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/text v0.3.7 // indirect
)
