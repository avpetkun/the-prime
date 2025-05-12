package fragment

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/avpetkun/the-prime/pkg/tonu"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

const fragmentCookie = "stel_ssid=wefwefewfwfwf; stel_dt=-12312; stel_token=fewfwefewfewf; stel_ton_token=fewfwefewfwfewfewf"

func TestApi(t *testing.T) {
	t.SkipNow()
	phrase := `12-24 words`

	log := zerolog.New(os.Stdout)

	tonApi, err := tonu.ConnectGlobal(context.TODO(), log)
	require.NoError(t, err)

	fragment, err := NewAPI(tonApi, zerolog.New(os.Stdout), fragmentCookie, phrase)
	require.NoError(t, err)

	var tx *SentTx
	if true {
		tx, err = fragment.SendTelegramStars(context.TODO(), log, "avpetkun", 55, false)
	} else {
		tx, err = fragment.SendTelegramPremium(context.TODO(), log, "avpetkundev", 1, false)
	}
	require.NoError(t, err)

	fmt.Println("sent tx", tx)
}
