package main

import (
	bitcoin_node "bitcoin_protocol/src"
	"encoding/hex"
	"fmt"
	"log"
	"net"
)

func main() {
	var magic uint32 = 0xd9b4bef9
	var tx_id = "fc57704eff327aecfadb2cf3774edc919ba69aba624b836461ce2be9c00a0c20"

	var peer_ip_address = "82.64.194.26"
	var peer_tcp_port = 8333

	var payload_version = bitcoin_node.CreatePayloadVersion(peer_ip_address, "Satoshi:0.7.2")

	var command [12]byte = [12]byte{}
	copy(command[:], "version")
	var version_message = bitcoin_node.CreateMessage(magic, command, payload_version.Blob())

	var verack_message = []byte(bitcoin_node.VERACK)

	var get_data_payload = bitcoin_node.CreateGetData([]string{tx_id})

	command = [12]byte{}
	copy(command[:], "getdata")
	var get_data_message = bitcoin_node.CreateMessage(magic, command, get_data_payload.Blob())

	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%d", peer_ip_address, peer_tcp_port))
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	var buffer = make([]byte, 1024)
	var request = version_message.Blob()
	conn.Write(request)
	n, _ := conn.Read(buffer)
	fmt.Printf("Command: version\n")
	fmt.Printf("Request: %v\n", hex.EncodeToString(request))
	fmt.Printf("Response: %v\n", hex.EncodeToString(buffer[:n]))

	buffer = make([]byte, 1024)
	request = verack_message
	conn.Write(request)
	n, _ = conn.Read(buffer)
	fmt.Printf("Command: verack\n")
	fmt.Printf("Request: %v\n", hex.EncodeToString(request))
	fmt.Printf("Response: %v\n", hex.EncodeToString(buffer[:n]))

	buffer = make([]byte, 1024)
	request = get_data_message.Blob()
	conn.Write(request)
	n, _ = conn.Read(buffer)
	fmt.Printf("Command: getdata\n")
	fmt.Printf("Request: %v\n", hex.EncodeToString(request))
	fmt.Printf("Response: %v\n", hex.EncodeToString(buffer[:n]))
}
