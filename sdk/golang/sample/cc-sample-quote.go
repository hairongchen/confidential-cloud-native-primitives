package main

import (
	"encoding/base64"
	"encoding/binary"
	"log"
	"math"
	"math/rand"

	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base"
	"github.com/hairongchen/confidential-cloud-native-primitives/sdk/golang/ccnp/"
)

func main() {

	sdk := ccnp.SDK{}

	nonce := makeNonce()
	userData := makeUserData()
	FlagFormat := "raw"

	report, err := sdk.GetCCReport(nonce, userData, nil)
	if err != nil {
		log.Fatalf("can not get cc report with error: %v", err)
	}

	report.Dump(cctrusted_base.QuoteDumpFormat(FlagFormat))
	return
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
