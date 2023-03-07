module github.com/ncarlier/webhookd/tooling/httpsig

go 1.19

require (
	github.com/go-fed/httpsig v1.1.0
	github.com/ncarlier/webhookd v1.0.0-00010101000000-000000000000
)

require (
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
)

replace github.com/ncarlier/webhookd => ../..
