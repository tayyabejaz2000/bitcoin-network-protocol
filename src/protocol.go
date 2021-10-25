package bitcoin_node

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"net"
	"time"
)

const VERACK = "\xf9\xbe\xb4\xd9\x76\x65\x72\x61\x63\x6b\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x5d\xf6\xe0\xe2"

const (
	SERVICE_NODE_NETWORK         = uint64(1)
	SERVICE_NODE_GETUTXO         = uint64(2)
	SERVICE_NODE_BLOOM           = uint64(4)
	SERVICE_NODE_WITNESS         = uint64(8)
	SERVICE_NODE_XTHIN           = uint64(16)
	SERVICE_NODE_COMPACT_FILTERS = uint64(64)
	SERVICE_NODE_NETWORK_FILTERS = uint64(1024)
)

const (
	MSG_ERROR                  = uint32(0)
	MSG_TX                     = uint32(1)
	MSG_BLOCK                  = uint32(2)
	MSG_FILTERED_BLOCK         = uint32(3)
	MSG_CMPCT_BLOCK            = uint32(4)
	MSG_WITNESS_TX             = uint32(0x40000001)
	MSG_WITNESS_BLOCK          = uint32(0x40000002)
	MSG_FILTERED_WITNESS_BLOCK = uint32(0x40000003)
)

type NetAddr struct {
	Services uint64
	Addr     net.IP
	Port     uint16
}

type Message struct {
	Magic    uint32
	Command  [12]byte
	Length   uint32
	Checksum uint32
	Payload  []byte
}

type PayloadVersion struct {
	Version     int32
	Services    uint64
	Timestamp   int64
	AddrRecv    NetAddr
	AddrFrom    NetAddr
	Nonce       uint64
	UserAgent   []byte
	StartHeight int32
}

type InventoryVector struct {
	Type uint32
	Hash [32]byte
}

type GetData struct {
	Count     uint8
	Inventory []InventoryVector
}

func CreateMessage(magic uint32, command [12]byte, payload []byte) Message {
	hash_1 := sha256.Sum256(payload)
	hash_2 := sha256.Sum256(hash_1[:])
	checksum := binary.LittleEndian.Uint32(hash_2[0:4])
	return Message{
		Magic:    magic,
		Command:  command,
		Length:   uint32(len(payload)),
		Checksum: checksum,
		Payload:  payload,
	}
}

func CreateNetworkAddress(ip_address string, port uint16) NetAddr {
	return NetAddr{
		Services: SERVICE_NODE_NETWORK,
		Addr:     net.ParseIP(ip_address),
		Port:     port,
	}
}

func CreateSubversion(sub_version string) []byte {
	var subversion = "/" + sub_version + "/"
	return bytes.Join([][]byte{{byte(len(subversion))}, []byte(subversion)}, []byte{})
}

func CreatePayloadVersion(peer_ip_address string, sub_version string) PayloadVersion {
	var version int32 = 60002
	var services = SERVICE_NODE_NETWORK
	var timestamp = time.Now().Unix()
	var addr_local = CreateNetworkAddress("127.0.0.1", 8333)
	var addr_peer = CreateNetworkAddress(peer_ip_address, 8333)
	var nonce = rand.Uint64()
	var start_height = 0
	return PayloadVersion{
		Version:     version,
		Services:    services,
		Timestamp:   timestamp,
		AddrFrom:    addr_local,
		AddrRecv:    addr_peer,
		Nonce:       nonce,
		UserAgent:   CreateSubversion(sub_version),
		StartHeight: int32(start_height),
	}
}

func CreateGetData(tx_ids []string) GetData {
	var getData = GetData{
		Count:     uint8(len(tx_ids)),
		Inventory: make([]InventoryVector, len(tx_ids)),
	}
	for i, tx_id := range tx_ids {
		decoded, _ := hex.DecodeString(tx_id)
		var hash [32]byte
		copy(hash[:], decoded)
		getData.Inventory[i] = InventoryVector{
			Type: MSG_TX,
			Hash: hash,
		}
	}
	return getData
}
