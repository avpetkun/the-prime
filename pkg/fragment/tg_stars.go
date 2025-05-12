package fragment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
)

func (api *API) SendTelegramStars(ctx context.Context, log zerolog.Logger, username string, quantity int, showSender bool) (*SentTx, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}
	if quantity < 50 {
		return nil, fmt.Errorf("quantity (%d) must be >= 50", quantity)
	}

	log = log.With().
		Str("username", username).
		Bool("show_sender", showSender).
		Int("stars_quantity", quantity).
		Logger()
	apiHash, err := api.getStarsApiHash(ctx)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("[fragment:tg_stars] got api hash %s", apiHash)
	recipient, found, err := api.searchStarsRecipient(ctx, apiHash, username, quantity)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("[fragment:tg_stars] got recipient %s", recipient)
	if !found {
		return nil, errors.Join(ErrUserNotFound, errors.New(username))
	}
	reqID, price, err := api.initBuyStarsRequest(ctx, apiHash, recipient, quantity)
	if err != nil {
		return nil, err
	}
	log.Info().Str("price", price).Msgf("[fragment:tg_stars] got request id %s", reqID)
	tx, err := api.getBuyStarsLink(ctx, apiHash, reqID, recipient, quantity, showSender)
	if err != nil {
		return nil, err
	}
	txMessage, err := tx.GetMessage()
	if err != nil {
		return nil, err
	}
	return api.sendTx(ctx, log, txMessage)
}

func (api *API) getStarsApiHash(ctx context.Context) (apiHash string, err error) {
	resp, err := api.doStarsRequest(ctx, "", "", 0, "")
	if err != nil {
		return "", err
	}
	defer resp.Close()
	body, err := io.ReadAll(resp)
	if err != nil {
		return "", err
	}
	return getApiHash(body)
}

func (api *API) searchStarsRecipient(ctx context.Context, apiHash, username string, quantity int) (recipient string, found bool, err error) {
	reqBody := fmt.Sprintf(`query=%s&quantity=%d&method=searchStarsRecipient`, username, quantity)
	resp, err := api.doStarsRequest(ctx, apiHash, reqBody, quantity, "")
	if err != nil {
		return "", false, err
	}
	defer resp.Close()

	var response SearchRecipientResponse
	err = json.NewDecoder(resp).Decode(&response)
	if err != nil {
		return "", false, err
	}
	if response.Error != "" {
		if response.Error == "No Telegram users found." {
			return "", false, nil
		}
		return "", false, errors.New(response.Error)
	}
	return response.Found.Recipient, response.OK, nil
}

func (api *API) initBuyStarsRequest(ctx context.Context, apiHash, recipient string, quantity int) (reqID, price string, err error) {
	reqBody := fmt.Sprintf(`recipient=%s&quantity=%d&method=initBuyStarsRequest`, recipient, quantity)
	resp, err := api.doStarsRequest(ctx, apiHash, reqBody, quantity, recipient)
	if err != nil {
		return "", "0.0", err
	}
	defer resp.Close()

	var response InitBuyResponse
	err = json.NewDecoder(resp).Decode(&response)
	if err != nil {
		return "", "0.0", err
	}
	if response.Error != "" {
		return "", "0.0", errors.New(response.Error)
	}
	return response.ReqID, response.Amount, nil
}

func (api *API) getBuyStarsLink(ctx context.Context, apiHash, reqID, recipient string, quantity int, showSender bool) (*TonTx, error) {
	showSenderFlag := 0
	if showSender {
		showSenderFlag = 1
	}
	reqBody := fmt.Sprintf(
		`account={"address":"%s","chain":"-239","walletStateInit":"te6cckECFgEAArEAAgE0ARUBFP8A9KQT9LzyyAsCAgEgAw4CAUgEBQLc0CDXScEgkVuPYyDXCx8gghBleHRuvSGCEHNpbnS9sJJfA+CCEGV4dG66jrSAINchAdB01yH6QDD6RPgo+kQwWL2RW+DtRNCBAUHXIfQFgwf0Dm+hMZEw4YBA1yFwf9s84DEg10mBAoC5kTDgcOIREAIBIAYNAgEgBwoCAW4ICQAZrc52omhAIOuQ64X/wAAZrx32omhAEOuQ64WPwAIBSAsMABezJftRNBx1yHXCx+AAEbJi+1E0NcKAIAAZvl8PaiaECAoOuQ+gLAEC8g8BHiDXCx+CEHNpZ2668uCKfxAB5o7w7aLt+yGDCNciAoMI1yMggCDXIdMf0x/TH+1E0NIA0x8g0x/T/9cKAAr5AUDM+RCaKJRfCtsx4fLAh98Cs1AHsPLQhFEluvLghVA2uvLghvgju/LQiCKS+ADeAaR/yMoAyx8BzxbJ7VQgkvgP3nDbPNgRA/btou37AvQEIW6SbCGOTAIh1zkwcJQhxwCzji0B1yggdh5DbCDXScAI8uCTINdKwALy4JMg1x0GxxLCAFIwsPLQiddM1zkwAaTobBKEB7vy4JPXSsAA8uCT7VXi0gABwACRW+Dr1ywIFCCRcJYB1ywIHBLiUhCx4w8g10oSExQAlgH6QAH6RPgo+kQwWLry4JHtRNCBAUHXGPQFBJ1/yMoAQASDB/RT8uCLjhQDgwf0W/LgjCLXCgAhbgGzsPLQkOLIUAPPFhL0AMntVAByMNcsCCSOLSHy4JLSAO1E0NIAURO68tCPVFAwkTGcAYEBQNch1woA8uCO4sjKAFjPFsntVJPywI3iABCTW9sx4ddM0ABRgAAAAD///4jUCrzR35U3Jh89K9UBCv58vz8UDChbN4GLMM07ZSEM4SARJZN3","publicKey":"%s"}&device={"platform":"browser","appName":"telegram-wallet","appVersion":"1","maxProtocolVersion":2,"features":["SendTransaction",{"name":"SendTransaction","maxMessages":255,"extraCurrencySupported":true}]}&transaction=1&id=%s&show_sender=%d&method=getBuyStarsLink`,
		api.walRawAddr, api.publicKey, reqID, showSenderFlag,
	)
	resp, err := api.doStarsRequest(ctx, apiHash, reqBody, quantity, recipient)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	var response GetLinkResponse
	err = json.NewDecoder(resp).Decode(&response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	if !response.OK {
		return nil, errors.New("getBuyStarsLink not ok")
	}
	return &response.Transaction, nil
}

func (api *API) doStarsRequest(ctx context.Context, apiHash, body string, quantity int, recipient string) (respBody io.ReadCloser, err error) {
	var req *http.Request
	if apiHash == "" {
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, "https://fragment.com/stars", nil)
	} else {
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, "https://fragment.com/api?hash="+apiHash, strings.NewReader(body))
	}
	if err != nil {
		return
	}

	h := req.Header

	h.Set("accept-language", "ru-RU,ru;q=0.9")
	h.Set("cache-control", "no-cache")
	h.Set("pragma", "no-cache")
	h.Set("sec-ch-ua-mobile", "?0")
	h.Set("sec-ch-ua-platform", `"macOS"`)
	h.Set("sec-fetch-site", "same-origin")
	h.Set("sec-ch-ua", `"Google Chrome";v="135", "Not-A.Brand";v="8", "Chromium";v="135"`)
	h.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
	h.Set("cookie", api.auth)

	if body == "" {
		h.Set("priority", "u=0, i")
		h.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		h.Set("sec-fetch-dest", "document")
		h.Set("sec-fetch-mode", "navigate")
		h.Set("sec-fetch-user", "?1")
		h.Set("upgrade-insecure-requests", "1")
	} else {
		h.Set("origin", "https://fragment.com")
		h.Set("priority", "u=1, i")
		h.Set("accept", "application/json, text/javascript, */*; q=0.01")
		h.Set("sec-fetch-dest", "empty")
		h.Set("sec-fetch-mode", "cors")
		h.Set("x-requested-with", "XMLHttpRequest")

		if recipient == "" {
			h.Set("referer", fmt.Sprintf("https://fragment.com/stars/buy?quantity=%d", quantity))
		} else {
			h.Set("referer", fmt.Sprintf("https://fragment.com/stars/buy?recipient=%s&quantity=%d", recipient, quantity))
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
