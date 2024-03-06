package main

import (
	"encoding/base64"
	"encoding/binary"
	"math"
	"math/rand"

	"github.com/hairongchen/confidential-cloud-native-primitives/sdk/golang/ccnp"
)

func main() {
	nonce := makeNonce()
	userData := makeUserData()

	report, err := ccnp.ccnpsdk.GetCCReport(nonce, userData, nil)
	if err != nil {
		return err
	}

	report.Dump(cctrusted_base.QuoteDumpFormat(FlagFormat))
	return nil
}

func makeNonce() string {
	num := uint64(rand.Int63n(math.MaxInt64))
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, num)
	return base64.StdEncoding.EncodeToString(b)
}

func makeUserData() string {
	b := []byte("demo user data")
	return base64.StdEncoding.EncodeToString(b)
}
