package p2p

type Peer interface {
	Close() error
}

type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
