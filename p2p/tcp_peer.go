package p2p

import "net"

type TCPPeer struct {
	conn         net.Conn
	outboundPeer bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:         conn,
		outboundPeer: outbound,
	}
}

// close implements the peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}
