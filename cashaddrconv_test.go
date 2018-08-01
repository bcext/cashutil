// Copyright (c) 2018 The bcext developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package cashutil

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/bcext/gcash/chaincfg"
)

var networks = []*chaincfg.Params{
	&chaincfg.MainNetParams,
	&chaincfg.TestNet3Params,
	&chaincfg.RegressionNetParams,
}

var validSizes = map[int]int{0: 20, 1: 24, 2: 28, 3: 32, 4: 40, 5: 48, 6: 56, 7: 64}

func InsecureGetRandUint160() []byte {
	ret := make([]byte, 20)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Read(ret)
	return ret
}

func InsecureGetRandomByteArray(size int) []byte {
	ret := make([]byte, size)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Read(ret)
	return ret
}

func PackCashaddrContent(content *addrContent) []byte {
	return packAddrData(content.hash, uint8(content.t))
}

func TestEncodeDecodeAllSizes(t *testing.T) {
	param := &chaincfg.MainNetParams
	for size1, size2 := range validSizes {
		data := InsecureGetRandomByteArray(size2)
		content := addrContent{pubKeyType, data}
		packedData := PackCashaddrContent(&content)

		// Check that the packed size is correct
		if int(packedData[1]>>2) != size1 {
			t.Error("the packed size is incorrect")
		}

		addr := encode(param.CashAddrPrefix, packedData)
		// Check that the address decodes properly
		decode := decodeCashaddrContent(addr, param)
		if !bytes.Equal(content.hash, decode.hash) {
			t.Error("cashaddr encode or decode error")
		}
	}
}

func TestCheckPackAddr(t *testing.T) {
	for _, size2 := range validSizes {
		data := InsecureGetRandomByteArray(size2 - 1)
		content := addrContent{pubKeyType, data}

		defer func() {
			if r := recover(); r != nil {
				// do nothing
			}
		}()

		PackCashaddrContent(&content) // error length causes panic
		t.Error("addr pack error")
	}
}

func Hash160Padding(origin []byte) []byte {
	if len(origin) >= 20 {
		return origin[:20]
	}

	buf := bytes.NewBuffer(make([]byte, 0, 20))
	buf.Write(origin)
	buf.Write(bytes.Repeat([]byte{0}, 20-len(origin)))

	return buf.Bytes()
}

func TestEncodeDecode(t *testing.T) {
	pubKeyHash, err := hex.DecodeString("0badf00d")
	if err != nil {
		panic(err)
	}

	scriptHash, err := hex.DecodeString("0f00dbad")
	if err != nil {
		panic(err)
	}

	for _, net := range networks {
		dst, err := NewAddressScriptHashFromHash(Hash160Padding(pubKeyHash), net)
		if err != nil {
			panic(err)
		}
		encode := dst.EncodeAddress(true)
		decode, err := decodeCashAddr(encode, net)
		if err != nil || decode.String() != dst.String() {
			t.Error("encode or decode error")
		}

	}

	for _, net := range networks {
		dst, err := NewAddressScriptHashFromHash(Hash160Padding(scriptHash), net)
		if err != nil {
			panic(err)
		}
		encode := dst.EncodeAddress(true)
		decode, err := decodeCashAddr(encode, net)
		if err != nil || decode.String() != dst.String() {
			t.Error("encode or decode error")
		}
	}
}

// Check that an encoded cash address is not valid on another network.
func TestInvalidOnWrongNetwork(t *testing.T) {
	b, err := hex.DecodeString("c0ffee")
	if err != nil {
		panic(err)
	}
	b20 := Hash160Padding(b)
	for _, net := range networks {
		for _, otherNet := range networks {
			if net == otherNet {
				continue
			}
			pubKeyHashAddr, err := NewAddressPubKeyHash(b20, net)
			if err != nil {
				panic(err)
			}
			encoded := pubKeyHashAddr.EncodeAddress(true)
			decoded, _ := decodeCashAddr(encoded, otherNet)

			if decoded != nil {
				t.Error("the decoded address should be nil as to incorrect network")
			}
		}
	}
}

func TestRandomDst(t *testing.T) {
	param := &chaincfg.MainNetParams

	for i := 0; i < 5000; i++ {
		hash := InsecureGetRandUint160()

		pubKeyDst, err := NewAddressPubKeyHash(hash, param)
		if err != nil {
			panic(err)
		}
		encodedKey := pubKeyDst.EncodeAddress(true)
		decodedKey, err1 := decodeCashAddr(encodedKey, param)

		scriptHashDst, err := NewAddressScriptHashFromHash(hash, param)
		if err != nil {
			panic(err)
		}
		encodedSrc := scriptHashDst.EncodeAddress(true)
		decodedSrc, err2 := decodeCashAddr(encodedSrc, param)

		errString := fmt.Sprintf("cashaddr failed for hash: %s", hex.EncodeToString(hash))
		if err1 != nil || pubKeyDst.String() != decodedKey.String() {
			t.Error(errString)
		}
		if err2 != nil || scriptHashDst.String() != decodedSrc.String() {
			t.Error(errString)
		}
	}
}

// Cashaddr payload made of 5-bit nibbles. The last one is padded. When
// converting back to bytes, this extra padding is truncated. In order to ensure
// cashaddr are cannonicals, we check that the data we truncate is zeroed.
func TestCheckPadding(t *testing.T) {
	var version uint8
	data := bytes.NewBuffer(make([]byte, 0, 34))
	data.WriteByte(version)
	data.Write(bytes.Repeat([]byte{1}, 33))

	if data.Len() != 34 {
		t.Errorf("origin data length should be: %d", data.Len())
	}

	param := &chaincfg.MainNetParams
	originBytes := data.Bytes()
	for i := 0; i < 32; i++ {
		originBytes[len(originBytes)-1] = byte(i)
		fake := encode(param.CashAddrPrefix, originBytes)
		dst, _ := decodeCashAddr(fake, param)

		// We have 168 bits of payload encoded as 170 bits in 5 bits nimbles. As
		// a result, we must have 2 zeros.
		if i&0x03 != 0 {
			if dst != nil {
				t.Error("check error")
			}
		} else {
			if dst == nil {
				t.Error("check error")
			}
		}
	}
}

// We ensure type is extracted properly from the version.
func TestCheckType(t *testing.T) {
	data := bytes.NewBuffer(make([]byte, 0, 34))
	data.Write(bytes.Repeat([]byte{0}, 34))
	param := &chaincfg.MainNetParams

	originBytes := data.Bytes()
	for i := byte(0); i < 16; i++ {
		originBytes[0] = i
		content := decodeCashaddrContent(encode(param.CashAddrPrefix, originBytes), param)
		if content.t != addrType(i) {
			t.Error("type error")
		}
		if len(content.hash) != 20 {
			t.Error("hash length error")
		}

		// Check that using the reserved bit result in a failure.
		originBytes[0] |= 0x10
		content = decodeCashaddrContent(encode(param.CashAddrPrefix, originBytes), param)
		if content != nil && content.t != 0 {
			t.Error("type error")
		}
		if content != nil && len(content.hash) != 0 {
			t.Error("hash length error")
		}
	}
}

// We ensure size is extracted and checked properly.
func TestCheckSize(t *testing.T) {
	param := &chaincfg.MainNetParams

	for size1, size2 := range validSizes {
		// Number of bytes required for a 5-bit packed version of a hash, with
		// version byte.  Add half a byte(4) so integer math provides the next
		// multiple-of-5 that would fit all the data.
		expectedSize := (8*(1+size2) + 4) / 5
		data := bytes.Repeat([]byte{0}, expectedSize)
		// After conversion from 8 bit packing to 5 bit packing, the size will
		// be in the second 5-bit group, shifted left twice.
		data[1] = byte(size1 << 2)
		content := decodeCashaddrContent(encode(param.CashAddrPrefix, data), param)
		if content != nil {
			if content.t != addrType(0) {
				t.Error("error")
			}

			if len(content.hash) != size2 {
				t.Error("error")
			}
		}

		data = bytes.Repeat([]byte{0}, expectedSize+1)
		content = decodeCashaddrContent(encode(param.CashAddrPrefix, data), param)
		if content != nil {
			if content.t != addrType(0) {
				t.Error("error")
			}
			if len(content.hash) != 0 {
				t.Error("error")
			}
		}

		data = data[:len(data)-2]
		content = decodeCashaddrContent(encode(param.CashAddrPrefix, data), param)
		if content != nil {
			if content.t != addrType(0) {
				t.Error("error")
			}
			if len(content.hash) != 0 {
				t.Error("error")
			}
		}
	}
}

func TestCashAddresses(t *testing.T) {
	param := &chaincfg.MainNetParams

	hash := [][]byte{
		{118, 160, 64, 83, 189, 160, 168, 139, 218, 81, 119, 184, 106, 21, 195, 178, 159, 85, 152, 115},
		{203, 72, 18, 50, 41, 156, 213, 116, 49, 81, 172, 75, 45, 99, 174, 25, 142, 123, 176, 169},
		{1, 31, 40, 228, 115, 201, 95, 64, 19, 215, 213, 62, 197, 251, 195, 180, 45, 248, 237, 16},
	}

	pubkey := []string{
		"bitcoincash:qpm2qsznhks23z7629mms6s4cwef74vcwvy22gdx6a",
		"bitcoincash:qr95sy3j9xwd2ap32xkykttr4cvcu7as4y0qverfuy",
		"bitcoincash:qqq3728yw0y47sqn6l2na30mcw6zm78dzqre909m2r",
	}

	script := []string{
		"bitcoincash:ppm2qsznhks23z7629mms6s4cwef74vcwvn0h829pq",
		"bitcoincash:pr95sy3j9xwd2ap32xkykttr4cvcu7as4yc93ky28e",
		"bitcoincash:pqq3728yw0y47sqn6l2na30mcw6zm78dzq5ucqzc37",
	}

	for index, item := range hash {
		// pubKey hash address
		dstKey, err := NewAddressPubKeyHash(item, param)
		if err != nil {
			panic(err)
		}
		if dstKey.EncodeAddress(true) != pubkey[index] {
			t.Error("encode bitcoin pubKeyHash address error")
		}

		// script hash address
		dstScript, err := NewAddressScriptHashFromHash(item, param)
		if err != nil {
			panic(err)
		}
		if dstScript.EncodeAddress(true) != script[index] {
			t.Error("encode bitcoin scriptHash address error")
		}
	}
}
