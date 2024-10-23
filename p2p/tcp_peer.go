package p2p

import "net"

type TCPPeer struct {
	net.Conn
	outboundPeer bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:         conn,
		outboundPeer: outbound,
	}
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}
