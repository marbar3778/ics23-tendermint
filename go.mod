module github.com/confio/ics23-tendermint

go 1.14

require (
	github.com/confio/ics23/go v0.0.0-20200604202538-6e2c36a74465 // v0.6.0
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200619181135-735a2e312b46
	github.com/etcd-io/bbolt v1.3.3 // indirect
	github.com/tendermint/tendermint v0.33.5
)

replace github.com/tendermint/tendermint => github.com/tendermint/tendermint v0.32.2-0.20200612160744-206c814a8e64
