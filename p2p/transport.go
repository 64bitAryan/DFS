package p2p

type Peer struct{}

type Transport interface {
	ListenAndAccept() error
}
