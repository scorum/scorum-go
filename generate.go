package scorumgo

//go:generate tools/bin/mockgen -source=caller/caller.go -destination=caller/caller_mock.go -self_package=github.com/scorum/scorum-go/caller -package=caller CallCloser
