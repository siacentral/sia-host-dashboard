package types

import (
	"bytes"
	"fmt"
	"math/big"

	siatypes "gitlab.com/NebulousLabs/Sia/types"
)

//BigNumber BigNumber
type BigNumber struct {
	i big.Int
}

// MarshalJSON implements the json.Marshaler interface.
func (b BigNumber) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, b.i.String())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface. An error is
// returned if a negative number is provided.
func (b *BigNumber) UnmarshalJSON(buf []byte) error {
	// UnmarshalJSON does not expect quotes
	buf = bytes.Trim(buf, `"`)
	err := b.i.UnmarshalJSON(buf)
	if err != nil {
		return err
	}
	return nil
}

//AddCurrency adds the currency value to the BigNumber and returns the result
func (b BigNumber) AddCurrency(a siatypes.Currency) (c BigNumber) {
	c.i.Add(&b.i, a.Big())

	return c
}

//SubCurrency subtracts the currency value from the BigNumber and returns the result
func (b BigNumber) SubCurrency(a siatypes.Currency) (c BigNumber) {
	c.i.Sub(&b.i, a.Big())

	return c
}

//Add adds the BigNumber value to the BigNumber and returns the result
func (b BigNumber) Add(a BigNumber) (c BigNumber) {
	c.i.Add(&b.i, &a.i)

	return c
}

//Sub subtracts the BigNumber value from the BigNumber and returns the result
func (b BigNumber) Sub(a BigNumber) (c BigNumber) {
	c.i.Sub(&b.i, &a.i)

	return c
}

//Div64 divides the BigNumber by the uint64 value and returns the result
func (b BigNumber) Div64(a uint64) (c BigNumber) {
	aB := new(big.Int)
	aB.SetUint64(a)
	c.i.Div(&b.i, aB)

	return c
}
