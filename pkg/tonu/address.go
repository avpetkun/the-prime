package tonu

import "github.com/xssnick/tonutils-go/address"

func ParseAnyAddress(addr string) (*address.Address, error) {
	adr, err := address.ParseAddr(addr)
	if err != nil {
		adr, err = address.ParseRawAddr(addr)
	}
	return adr, err
}

func ParseAnyAddressString(addr string) (string, error) {
	a, err := ParseAnyAddress(addr)
	if err != nil {
		return "", err
	}
	// return a.Bounce(false).String(), nil
	return a.String(), nil
}
