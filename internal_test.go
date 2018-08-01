// Copyright (c) 2013-2017 The btcsuite developers
// Copyright (c) 2018 The bcext developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

/*
This test file is part of the cashutil package rather than than the
cashutil_test package so it can bridge access to the internals to properly test
cases which are either not possible or can't reliably be tested via the public
interface. The functions are only exported while the tests are being run.
*/

package cashutil

import (
	"strings"

	"github.com/bcext/cashutil/base58"
	"github.com/bcext/gcash/btcec"
	"github.com/bcext/gcash/chaincfg"
	"golang.org/x/crypto/ripemd160"
)

// SetBlockBytes sets the internal serialized block byte buffer to the passed
// buffer.  It is used to inject errors and is only available to the test
// package.
func (b *Block) SetBlockBytes(buf []byte) {
	b.serializedBlock = buf
}

// TstAppDataDir makes the internal appDataDir function available to the test
// package.
func TstAppDataDir(goos, appName string, roaming bool) string {
	return appDataDir(goos, appName, roaming)
}

// TstAddressPubKeyHash makes an AddressPubKeyHash, setting the
// unexported fields with the parameters hash and netID.
func TstAddressPubKeyHash(hash [ripemd160.Size]byte,
	param *chaincfg.Params) *AddressPubKeyHash {

	return &AddressPubKeyHash{
		hash: hash,
		net:  param,
	}
}

// TstAddressScriptHash makes an AddressScriptHash, setting the
// unexported fields with the parameters hash and netID.
func TstAddressScriptHash(hash [ripemd160.Size]byte,
	param *chaincfg.Params) *AddressScriptHash {

	return &AddressScriptHash{
		hash: hash,
		net:  param,
	}
}

// TstAddressPubKey makes an AddressPubKey, setting the unexported fields with
// the parameters.
func TstAddressPubKey(serializedPubKey []byte, pubKeyFormat PubKeyFormat,
	param *chaincfg.Params) *AddressPubKey {

	pubKey, _ := btcec.ParsePubKey(serializedPubKey, btcec.S256())
	return &AddressPubKey{
		pubKeyFormat: pubKeyFormat,
		pubKey:       (*btcec.PublicKey)(pubKey),
		net:          param,
	}
}

var netAddressMapping = map[string]*chaincfg.Params{
	"bitcoincash": &chaincfg.MainNetParams,
	"bchtest":     &chaincfg.TestNet3Params,
	"bchreg":      &chaincfg.RegressionNetParams,
}

// TstAddressSAddr returns the expected script address bytes for
// P2PKH and P2SH bitcoin addresses.
func TstAddressSAddr(addr string) []byte {
	if strings.Contains(addr, ":") {
		ret := strings.Split(addr, ":")
		net := netAddressMapping[ret[0]]
		address, err := DecodeAddress(addr, net)
		if err != nil {
			return nil
		}

		return address.ScriptAddress()
	}

	decoded := base58.Decode(addr)
	return decoded[1 : 1+ripemd160.Size]
}
