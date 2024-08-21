package signal

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"
)

// IndexOf returns the index of the Client with the given connectionId, or -1 if not found
func IndexOf(connectionId string, roomClients []Client) int {
	for i, client := range roomClients {
		if client.connectionId == connectionId {
			return i
		}
	}
	return -1
}

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
)

func randomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		result[i] = alphabet[num.Int64()]
	}
	return string(result)
}

func CreateConnectionId() string {
	// CUID parts
	cuidPrefix := "c"
	timestamp := time.Now().UnixNano()
	counter := randomString(4)
	clientFingerprint := randomString(4)
	randomBlock := randomString(8)

	// Construct the CUID
	return fmt.Sprintf("%s%x%s%s%s", cuidPrefix, timestamp, counter, clientFingerprint, randomBlock)
}

func DecodeQueryData(queryData string) (map[string]string, error) {
	output := make(map[string]string)
	if queryData == "" {
		return output, nil
	}
	// URL-decode the 'queryData' parameter
	decodedQueryData, err := url.QueryUnescape(queryData)
	if err != nil {
		return nil, fmt.Errorf("error decoding queryData: %s", err)
	}

	// Parse the decoded query data into key-value pairs
	params := strings.Split(decodedQueryData, "&")
	for _, param := range params {
		kv := strings.SplitN(param, "=", 2)
		if len(kv) == 2 {
			output[kv[0]] = kv[1]
		}
	}
	return output, nil
}
