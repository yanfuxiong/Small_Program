package utils

import (
	"crypto/md5"
	"encoding/pem"
	"fmt"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"log"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p/core/crypto"
	ma "github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"
)

func ContentEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func ConcatIP(ip string, port string) string {
	publicIP := fmt.Sprintf("%s:%s", ip, port)
	return publicIP
}

func SplitIP(ip string) (string, string) {
	parts := strings.Split(ip, ":")
	return parts[0], parts[1]
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func WriteNodeID(ID string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write([]byte(ID))
	if err != nil {
		log.Println(err)
	}
}

func WriteErrJson(name string, strContent []byte) {
	fileName := fmt.Sprintf("/storage/emulated/0/Android/data/com.rtk.myapplication/files/%s.log", name)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(strContent)
	if err != nil {
		log.Println(err)
	}
}

func CreateMD5Hash(data []byte) (mh.Multihash, error) {
	hash := md5.Sum(data)

	multihash, err := mh.Encode(hash[:], mh.MD5)
	if err != nil {
		return nil, err
	}

	return multihash, nil
}

func MarshalPrivateKeyToPEM(key crypto.PrivKey) ([]byte, error) {
	encoded, err := crypto.MarshalPrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %v", err)
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: encoded,
	})
	return pemEncoded, nil
}

func UnmarshalPrivateKeyFromPEM(pemData []byte) (crypto.PrivKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}
	return crypto.UnmarshalPrivateKey(block.Bytes)
}

func GenKey(privKeyFile string) crypto.PrivKey {
	var priv crypto.PrivKey
	var err error
	var content []byte
	content, err = os.ReadFile(privKeyFile)
	if err != nil {
		priv, _, err = crypto.GenerateKeyPair(crypto.RSA, 2048)
		if err != nil {
			log.Fatal(err)
		}

		jsonData, err := MarshalPrivateKeyToPEM(priv)
		err = os.WriteFile(privKeyFile, jsonData, 0644)
		if err != nil {
			log.Fatal(err)
		}
		return priv
	}

	priv, err = UnmarshalPrivateKeyFromPEM(content)

	return priv
}

func GetLocalPort(Addrs []ma.Multiaddr) string {
	var localPort string
	for _, maddr := range Addrs {
		protocols := maddr.Protocols()
		hasTCP := false
		hasIP4 := false
		for _, protocol := range protocols {
			if protocol.Code == ma.P_TCP {
				hasTCP = true
			}
			if protocol.Code == ma.P_IP4 {
				hasIP4 = true
			}
		}
		if hasTCP && hasIP4 {
			port, err := maddr.ValueForProtocol(ma.P_TCP)
			if err != nil {
				return ""
			}
			localPort = port
			break
		}
	}

	//log.Println("Local port: " + localPort)
	return localPort
}

func ExtractTCPIPandPort(maddr ma.Multiaddr) (string, string) {
	ip, err := maddr.ValueForProtocol(ma.P_IP4)
	if err != nil {
		log.Println("Failed to get IP: %v", err)
	}

	port, err := maddr.ValueForProtocol(ma.P_TCP)
	if err != nil {
		log.Println("Failed to get port: %v", err)
	}
	return ip, port
}

func GetRemoteAddrFromStream(stream network.Stream) string {
	ip, port := ExtractTCPIPandPort(stream.Conn().RemoteMultiaddr())
	return ConcatIP(ip, port)
}

func RemoveMySelfID(slice []string, s string) []string {
	i := 0
	for _, v := range slice {
		if v != s {
			slice[i] = v
			i++
		}
	}
	return slice[:i]
}

func IsInPeerList(peerID string, list []peer.AddrInfo) bool {
	for _, item := range list {
		if strings.EqualFold(item.ID.String(), peerID) {
			return true
		}
	}
	return false
}

func DistinctMdnsID(slice []string, list []peer.AddrInfo) []string {

	for _, peer := range list {
		slice = RemoveMySelfID(slice, peer.ID.String())
	}
	return slice
}
