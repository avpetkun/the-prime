package common

import "strings"

// TappAdsRequest represents the request parameters for TapAds tasks
type TappAdsRequest struct {
	UserID    int64  `json:"user_id"`
	IP        string `json:"ip"`
	Lang      string `json:"lang"`
	IsPremium bool   `json:"is_premium"`
	UserAgent string `json:"user_agent"`
}

// TappAdsAPIResponse represents the response from TapAds API
type TappAdsAPIResponse []TappAdsTask

// TappAdsTask represents a single task from TapAds API
type TappAdsTask struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Icon            string  `json:"icon"`
	Description     string  `json:"description"`
	NameWithDesc    string  `json:"name_with_description"`
	URL             string  `json:"url"`
	Payout          float64 `json:"payout"`
	Currency        string  `json:"currency"`
	IsDone          bool    `json:"is_done"`
	ClickPostback   string  `json:"click_postback"`
	ButtonLabel     string  `json:"btn_label"`
	LongDescription string  `json:"long_description"`
}

var tappAdsSplitters = []string{"-", "–", "—", "&", ":"}

func (ads *TappAdsTask) WithUserTask(t UserTask) UserTask {
	t.SubID = ads.ID
	t.Icon = ads.Icon
	t.ActionLink = ads.URL

	t.Name = ads.Name
	t.Desc = ads.Description

	for _, splitter := range tappAdsSplitters {
		if i := strings.Index(t.Name, splitter); i != -1 {
			desc := strings.TrimSpace(strings.TrimPrefix(t.Name[i:], splitter))
			if desc != "" {
				if t.Desc != "" && !strings.HasSuffix(desc, ".") && !strings.HasSuffix(desc, "!") {
					desc += "."
				}
				t.Desc = strings.TrimSpace(upperFirst(desc) + " " + t.Desc)
			}
			t.Name = strings.TrimSpace(t.Name[:i])
			break
		}
	}

	return t
}

func upperFirst(s string) string {
	for _, r := range s {
		first := string(r)
		return strings.ToUpper(first) + s[len(first):]
	}
	return s
}
