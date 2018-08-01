// Copyright (c) 2018 The bcext developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package cashutil

import (
	"bytes"
	"testing"
)

func CashAddrDecode(str string) (string, []byte) {
	return decode(str, "")
}

func CaseInsensitiveEqual(s1 string, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		c1 := int8(s1[i])
		if c1 >= 'A' && c1 <= 'Z' {
			c1 -= 'A' - 'a'
		}
		c2 := int8(s2[i])
		if c2 >= 'A' && c2 <= 'Z' {
			c2 -= 'A' - 'a'
		}
		if c1 != c2 {
			return false
		}
	}

	return true
}

func TestVectorsValid(t *testing.T) {
	cases := [...]string{
		"prefix:x64nx6hz",
		"PREFIX:X64NX6HZ",
		"p:gpf8m4h7",
		"bitcoincash:qpzry9x8gf2tvdw0s3jn54khce6mua7lcw20ayyn",
		"bchtest:testnetaddress4d6njnut",
		"bchreg:555555555555555555555555555555555555555555555udxmlmrz",
	}

	for _, str := range cases {
		prefix, payload := CashAddrDecode(str)
		if len(prefix) == 0 {
			t.Errorf("cashaddr: %s should not be decoded  empty", str)
		}
		recode := encode(prefix, payload)
		if len(recode) == 0 {
			t.Error("encode cashaddr should not be empty")
		}

		if !CaseInsensitiveEqual(str, recode) {
			t.Error("Decoded string should be equal to the orgin string")
		}
	}
}

func TestVectorInvalid(t *testing.T) {
	cases := [...]string{
		"prefix:x32nx6hz",
		"prEfix:x64nx6hz",
		"prefix:x64nx6Hz",
		"pref1x:6m8cxv73",
		"prefix:",
		":u9wsx07j",
		"bchreg:555555555555555555x55555555555555555555555555udxmlmrz",
		"bchreg:555555555555555555555555555555551555555555555udxmlmrz",
		"pre:fix:x32nx6hz",
		"prefixx64nx6hz",
	}

	for _, str := range cases {
		prefix, _ := CashAddrDecode(str)
		if len(prefix) != 0 {
			t.Errorf("The cashaddr: %s should be invalid!", str)
		}
	}
}

func TestRawEncode(t *testing.T) {
	prefix := "helloworld"
	payload := []byte{0x1f, 0x0d}

	encode := encode(prefix, payload)
	decodedPrefix, decodedPayload := CashAddrDecode(encode)

	if prefix != decodedPrefix {
		t.Error("encode() and Decode() are not matched!")
	}

	if !bytes.Equal(payload, decodedPayload) {
		t.Error("encode() and Decode() are not matched!")
	}
}

type cashaddr struct {
	prefix  string
	payload string
}

func TestVectorsNoPrefix(t *testing.T) {
	cases := []cashaddr{
		{"bitcoincash", "qpzry9x8gf2tvdw0s3jn54khce6mua7lcw20ayyn"},
		{"prefix", "x64nx6hz"},
		{"PREFIX", "X64NX6HZ"},
		{"p", "gpf8m4h7"},
		{"bitcoincash", "qpzry9x8gf2tvdw0s3jn54khce6mua7lcw20ayyn"},
		{"bchtest", "testnetaddress4d6njnut"},
		{"bchreg", "555555555555555555555555555555555555555555555udxmlmrz"},
	}

	for _, item := range cases {
		addr := item.prefix + ":" + item.payload
		prefix, payload := decode(item.payload, item.prefix)
		if !CaseInsensitiveEqual(prefix, item.prefix) {
			t.Errorf("cashaddr prefix: %s decode error", item.prefix)
		}

		recode := encode(prefix, payload)
		if len(recode) == 0 {
			t.Errorf("cashaddr: %s encode error", addr)
		}

		if !CaseInsensitiveEqual(addr, recode) {
			t.Errorf("cashaddr: %s encode error", addr)
		}
	}
}
