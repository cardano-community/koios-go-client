module github.com/cardano-community/koios-go-client/internal/testdatagen

go 1.16

replace github.com/cardano-community/koios-go-client v1.3.0 => ../../
// introduces breaking change since v1.3.0 (api change)

require (
	github.com/cardano-community/koios-go-client v1.3.0
	github.com/urfave/cli/v2 v2.3.0
)
