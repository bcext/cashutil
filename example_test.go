package btcutil_test

import (
	"fmt"
	"math"

	"github.com/btcsuite/btcutil"
)

func ExampleAmount() {

	a := btcutil.Amount(0)
	fmt.Println("Zero Satoshi:", a)

	a = btcutil.Amount(1e8)
	fmt.Println("100,000,000 Satoshis:", a)

	a = btcutil.Amount(1e5)
	fmt.Println("100,000 Satoshis:", a)
	// Output:
	// Zero Satoshi: 0 BCH
	// 100,000,000 Satoshis: 1 BCH
	// 100,000 Satoshis: 0.001 BCH
}

func ExampleNewAmount() {
	amountOne, err := btcutil.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := btcutil.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := btcutil.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := btcutil.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 BCH
	// 0.01234567 BCH
	// 0 BCH
	// invalid bitcoin amount
}

func ExampleAmount_unitConversions() {
	amount := btcutil.Amount(44433322211100)

	fmt.Println("Satoshi to kBCH:", amount.Format(btcutil.AmountKiloBCH))
	fmt.Println("Satoshi to BCH:", amount)
	fmt.Println("Satoshi to MilliBCH:", amount.Format(btcutil.AmountMilliBCH))
	fmt.Println("Satoshi to MicroBCH:", amount.Format(btcutil.AmountMicroBCH))
	fmt.Println("Satoshi to Satoshi:", amount.Format(btcutil.AmountSatoshi))

	// Output:
	// Satoshi to kBCH: 444.333222111 kBCH
	// Satoshi to BCH: 444333.222111 BCH
	// Satoshi to MilliBCH: 444333222.111 mBCH
	// Satoshi to MicroBCH: 444333222111 μBCH
	// Satoshi to Satoshi: 44433322211100 Satoshi
}
