package tonu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseAnyAddressString(t *testing.T) {
	textAddr, err := ParseAnyAddressString("0:e2233f689a19b86d85757140df8f3dd6471fc838139353aa36f474d41a0b3d46")
	require.NoError(t, err)
	require.Equal(t, "EQDiIz9omhm4bYV1cUDfjz3WRx_IOBOTU6o29HTUGgs9RsJi", textAddr)
	textAddr, err = ParseAnyAddressString("EQDiIz9omhm4bYV1cUDfjz3WRx_IOBOTU6o29HTUGgs9RsJi")
	require.NoError(t, err)
	require.Equal(t, "EQDiIz9omhm4bYV1cUDfjz3WRx_IOBOTU6o29HTUGgs9RsJi", textAddr)

	addr, err := ParseAnyAddress("0:e2233f689a19b86d85757140df8f3dd6471fc838139353aa36f474d41a0b3d46")
	require.NoError(t, err)
	require.Equal(t, textAddr, addr.String())

	addr, err = ParseAnyAddress("EQDiIz9omhm4bYV1cUDfjz3WRx_IOBOTU6o29HTUGgs9RsJi")
	require.NoError(t, err)
	require.Equal(t, textAddr, addr.String())
}
