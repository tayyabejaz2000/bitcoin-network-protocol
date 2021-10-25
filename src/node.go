package bitcoin_node

import (
	"fmt"
	"net"
)

type Node struct {
	Socket     net.Conn
	peerAddr   net.IP
	port       uint16
	subVersion string
}

func CreateNode(peer_address string, port uint16) Node {
	socket, _ := net.Dial("tcp4", fmt.Sprintf("%s:%d", peer_address, port))
	return Node{
		Socket:     socket,
		peerAddr:   net.ParseIP(peer_address),
		port:       port,
		subVersion: "Satoshi:0.7.2",
	}
}

func (n *Node) GetPeerAddress() string {
	return n.peerAddr.String()
}
func (n *Node) GetPort() uint16 {
	return n.port
}
func (n *Node) GetSubversion() string {
	return n.subVersion
}

func (n *Node) Close() {
	n.Socket.Close()
}
