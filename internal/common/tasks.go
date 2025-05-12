package common

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

type TaskType string

const (
	TaskFree            TaskType = "free"
	TaskInvite          TaskType = "invite"
	TaskJoin            TaskType = "join"
	TaskFreeLink        TaskType = "free_link"
	TaskPartnerEvent    TaskType = "partner_event"
	TaskPartnerCheck    TaskType = "partner_check"
	TaskTonConnect      TaskType = "ton_connect"
	TaskTonDisconnect   TaskType = "ton_disconnect"
	TaskTonDeposit      TaskType = "ton_deposit"
	TaskStarsDeposit    TaskType = "stars_deposit"
	TaskAdsGramTask     TaskType = "ads_gram_task"
	TaskAdsGramRewarded TaskType = "ads_gram_rewarded"
	TaskTappAds         TaskType = "tapp_ads"
	TaskMonetagLink     TaskType = "monetag-link"
	TaskMonetagBanner   TaskType = "monetag-banner"
)

type TaskStatus string

const (
	TaskActive  TaskStatus = "active"
	TaskPending TaskStatus = "pending"
	TaskClaim   TaskStatus = "claim"
	TaskDone    TaskStatus = "done"
)

type TaskKey struct {
	TaskID int64 `json:"id"`
	SubID  int64 `json:"subID"`
}

type TaskEvent struct {
	TaskKey
	Time   int64      `json:"time"`
	Status TaskStatus `json:"status"`
}

type TaskFlow struct {
	TaskKey
	Start  int64      `json:"start"`
	Status TaskStatus `json:"status"`
}

func NewTaskFlow(taskID, subID int64) *TaskFlow {
	return &TaskFlow{
		TaskKey: TaskKey{
			TaskID: taskID,
			SubID:  subID,
		},
		Status: TaskActive,
	}
}

type UserTask struct {
	TaskKey
	Type    TaskType `json:"type"`
	Name    string   `json:"name"`
	Desc    string   `json:"desc"`
	Icon    string   `json:"icon"`
	Premium bool     `json:"premium"`
	Points  int64    `json:"points"`
	Weight  int64    `json:"weight"`
	Hidden  bool     `json:"hidden"`

	Interval int64 `json:"interval"`
	Pending  int64 `json:"pending"`

	ActionLink string `json:"actionLink,omitempty"`

	ActionAdsGramBlockID string `json:"actionAdsGramBlockID,omitempty"`

	// filled
	Status TaskStatus `json:"status"`
	Start  int64      `json:"start"`
}

func (t *UserTask) CanStart(flow *TaskFlow, nowTime int64) bool {
	if flow.Status == TaskClaim {
		return false
	}
	if flow.Status == TaskPending {
		if t.Pending <= 0 || flow.Start+t.Pending > nowTime {
			return false
		}
	}
	if t.Interval == 0 {
		if flow.Status == TaskDone {
			return false
		}
	} else if flow.Start+t.Interval > nowTime {
		return false
	}
	return true
}

type FullTask struct {
	UserTask

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`

	MaxClicks int64 `json:"maxClicks"`
	NowClicks int64 `json:"nowClicks"`

	ActionChatID    int64   `json:"actionChatID"`
	ActionTonAmount float64 `json:"actionTonAmount"`

	ActionStarsAmount int    `json:"actionStarsAmount"`
	ActionStarsTitle  string `json:"actionStarsTitle"`
	ActionStarsDesc   string `json:"actionStarsDesc"`
	ActionStarsItem   string `json:"actionStarsItem"`

	ActionPartnerHook  string `json:"actionPartnerHook"`
	ActionPartnerMatch string `json:"actionPartnerMatch"`
	ActionTappAdsToken string `json:"actionTappAdsToken"`
}

func (t *FullTask) ActionTonAmountUnits() int64 {
	units := int64(math.Ceil(t.ActionTonAmount * 1e9))
	if units < 0 {
		panic("negative ton amount for task")
	}
	return units
}

func (t *FullTask) Valid() error {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" {
		return fmt.Errorf("empty task name")
	}
	t.Desc = strings.TrimSpace(t.Desc)
	if t.Desc == "" {
		return fmt.Errorf("empty task description")
	}
	if t.Points < 0 {
		return fmt.Errorf("negative task points")
	}
	if t.Interval < 0 {
		return fmt.Errorf("invalid task interval")
	}
	if t.Pending < 0 || t.Pending > 0 && t.Pending < t.Interval {
		return fmt.Errorf("invalid max pending time")
	}

	t.ActionPartnerHook = strings.TrimSpace(t.ActionPartnerHook)
	t.ActionPartnerMatch = strings.TrimSpace(t.ActionPartnerMatch)
	t.ActionTappAdsToken = strings.TrimSpace(t.ActionTappAdsToken)
	t.ActionStarsTitle = strings.TrimSpace(t.ActionStarsTitle)
	t.ActionStarsDesc = strings.TrimSpace(t.ActionStarsDesc)
	t.ActionStarsItem = strings.TrimSpace(t.ActionStarsItem)

	switch t.Type {
	case TaskFree:
	case TaskInvite:
	case TaskJoin:
		t.Interval = 0
		if t.ActionLink == "" || t.ActionChatID == 0 {
			return fmt.Errorf("task action chat is no set")
		}
	case TaskFreeLink:
		t.Interval = 0
		if t.ActionLink == "" {
			return fmt.Errorf("task link is empty")
		}
	case TaskTonConnect, TaskTonDisconnect:
		t.Interval = 0
	case TaskTonDeposit:
		if t.ActionTonAmount <= 0 {
			return fmt.Errorf("task ton amount is no set")
		}
	case TaskStarsDeposit:
		if t.ActionStarsAmount <= 0 {
			return fmt.Errorf("task stars amount is no set")
		}
		if t.ActionStarsTitle == "" {
			return fmt.Errorf("task stars title is no set")
		}
		if t.ActionStarsDesc == "" {
			return fmt.Errorf("task stars desc is no set")
		}
		if t.ActionStarsItem == "" {
			return fmt.Errorf("task stars label is no set")
		}
	case TaskPartnerEvent, TaskPartnerCheck:
		if t.ActionLink == "" {
			return fmt.Errorf("task partner link is empty")
		}
		if t.ActionPartnerHook == "" {
			return fmt.Errorf("task partner token is empty")
		}
		if t.Type == TaskPartnerCheck {
			if !strings.Contains(t.ActionPartnerHook, "123456789") {
				return fmt.Errorf("user_id (123456789) not found in link")
			}
			parts := strings.Fields(t.ActionPartnerHook)
			switch len(parts) {
			case 1:
				t.ActionPartnerHook = http.MethodGet + " " + makeHttps(parts[0])
			case 2:
				method := strings.ToUpper(parts[0])
				switch method {
				case http.MethodGet, http.MethodPost, http.MethodPut:
				default:
					return fmt.Errorf("wrong partner check method: %s", method)
				}
				t.ActionPartnerHook = method + " " + makeHttps(parts[1])
			default:
				return fmt.Errorf("partner hook link contains too many parts")
			}
		}
	case TaskAdsGramTask, TaskAdsGramRewarded:
		if t.ActionPartnerHook == "" {
			return fmt.Errorf("ads_gram partner token is empty")
		}
		if t.ActionAdsGramBlockID == "" {
			return fmt.Errorf("ads_gram block-id is empty")
		}
	case TaskTappAds:
		if t.ActionPartnerHook == "" {
			return fmt.Errorf("tapp_ads partner token is empty")
		}
		if t.ActionTappAdsToken == "" {
			return fmt.Errorf("tapp_ads token is empty")
		}
	case TaskMonetagLink, TaskMonetagBanner:
		if t.ActionPartnerHook == "" {
			return fmt.Errorf("monetag partner token is empty")
		}
	default:
		return fmt.Errorf("invalid task type")
	}

	return nil
}

func makeHttps(s string) string {
	if strings.HasPrefix(s, "https://") {
		return s
	}
	if strings.HasPrefix(s, "http://") {
		return s[:4] + "s" + s[4:]
	}
	return "https://" + s
}
