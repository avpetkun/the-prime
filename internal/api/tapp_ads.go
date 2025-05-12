package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/avpetkun/the-prime/internal/common"
	"github.com/rs/zerolog"
)

const tappAdsBaseURL = "https://wallapi.tappads.io/v1"

func TappAdsGetTasks(ctx context.Context, apiKey string, ip, userAgent, lang string, isPremium bool, userID int64) ([]common.TappAdsTask, error) {
	params := url.Values{}
	params.Add("apikey", apiKey)
	params.Add("user_id", strconv.FormatInt(userID, 10))
	params.Add("ip", ip)
	params.Add("lang", lang)
	params.Add("is_premium", strconv.FormatBool(isPremium))
	params.Add("ua", userAgent)

	reqURL := fmt.Sprintf("%s/feed?%s", tappAdsBaseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode response
	var tasks common.TappAdsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return tasks, nil
}

func TappAdsSendClick(ctx context.Context, log zerolog.Logger, apiKey, ip, userAgent string, userID, subTaskID int64) error {
	params := url.Values{}
	params.Add("apikey", apiKey)
	params.Add("user_id", strconv.FormatInt(userID, 10))
	params.Add("ip", ip)
	params.Add("ua", userAgent)
	params.Add("offer_id", strconv.FormatInt(subTaskID, 10))

	reqURL := fmt.Sprintf("%s/click?%s", tappAdsBaseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w, url: %s", err, reqURL)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w, url: %s", err, reqURL)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, url: %s, response: %s", resp.StatusCode, reqURL, respData)
	}

	log.Info().
		Str("url", reqURL).
		Int64("user_id", userID).
		Int64("subtask_id", subTaskID).
		Msgf("[tapp_ads] click sent")

	return nil
}
