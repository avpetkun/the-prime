package api

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/avpetkun/the-prime/internal/cache"
	"github.com/avpetkun/the-prime/internal/common"
	"github.com/avpetkun/the-prime/internal/dbx"
	"github.com/avpetkun/the-prime/internal/worker"
	"github.com/avpetkun/the-prime/pkg/natsu"
	"github.com/avpetkun/the-prime/pkg/tgu"
	"github.com/avpetkun/the-prime/pkg/tonu"
)

type Service struct {
	cfg Config
	db  *dbx.DB
	ch  cache.Cache
	ns  *natsu.Stream
	bot *tgu.Bot
	log zerolog.Logger

	allTasks    []*common.FullTask
	allProducts []*common.Product
}

func NewService(cfg Config, log zerolog.Logger, db *dbx.DB, ch cache.Cache, ns *natsu.Stream, bot *tgu.Bot) *Service {
	return &Service{
		cfg: cfg,
		db:  db,
		ch:  ch,
		ns:  ns,
		bot: bot,
		log: log,
	}
}

func (s *Service) Start(ctx context.Context) error {
	if err := s.startLoadTasksLoop(ctx); err != nil {
		return err
	}
	if err := s.startLoadProductsLoop(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Service) ParseUserAuth(authData string) (*tgu.Auth, error) {
	return s.bot.AuthParse(authData)
}
func (s *Service) ParseUserAuth64(authData string) (*tgu.Auth, error) {
	return s.bot.AuthParseBase64(authData)
}

//

func (s *Service) GetAllBotChats(ctx context.Context) ([]tgu.Chat, error) {
	chats, err := s.db.GetAllBotChats(ctx)
	if err != nil {
		return nil, err
	}
	if chats == nil {
		return []tgu.Chat{}, nil
	}
	return chats, nil
}

//

func (s *Service) GetAllProducts(ctx context.Context) ([]*common.Product, error) {
	return s.db.GetAllProducts(ctx)
}

func (s *Service) AllProducts() []*common.Product {
	return s.allProducts
}

func (s *Service) SaveProduct(ctx context.Context, p *common.Product) error {
	if p == nil {
		return fmt.Errorf("empty product")
	}
	if err := p.Valid(); err != nil {
		return err
	}
	return s.db.SaveProduct(ctx, p)
}

func (s *Service) DeleteProduct(ctx context.Context, productID int64) error {
	return s.db.DeleteProduct(ctx, productID)
}

//

func sortTasks(tasks []*common.FullTask) {
	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Premium != tasks[j].Premium {
			return tasks[i].Premium
		}
		return tasks[i].Weight > tasks[j].Weight
	})
}
func sortUserTasks(tasks []common.UserTask) {
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Weight > tasks[j].Weight
	})
}

//

func (s *Service) GetAllTasksWithClicks(ctx context.Context) ([]*common.FullTask, error) {
	tasks, err := s.db.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}
	sortTasks(tasks)

	for _, t := range tasks {
		if t.MaxClicks > 0 && t.NowClicks < t.MaxClicks {
			clicks, err := s.ch.GetTaskClicks(ctx, t.TaskID)
			if err != nil {
				s.log.Error().Err(err).Msg("[clicks] failed to get task clicks")
			} else {
				t.NowClicks = clicks
			}
		}
	}

	return tasks, nil
}

func (s *Service) SaveTask(ctx context.Context, t *common.FullTask) error {
	if t == nil {
		return fmt.Errorf("empty task")
	}
	if err := t.Valid(); err != nil {
		return err
	}
	if t.TaskID != 0 {
		return s.db.UpdateTask(ctx, t)
	}
	return s.db.CreateTask(ctx, t)
}

func (s *Service) DeleteTask(ctx context.Context, taskID int64) error {
	if taskID <= 0 {
		return fmt.Errorf("invalid task id")
	}
	return s.db.DeleteTask(ctx, taskID)
}

//

func (s *Service) ClaimProduct(ctx context.Context, user tgu.User, productID int64) (spendPoints int64, err error) {
	product, err := s.getProduct(productID)
	if err != nil {
		return 0, err
	}
	if product.Amount <= 0 || product.Price <= 0 {
		return 0, fmt.Errorf("product claim not available")
	}

	if product.NeedUsername() && user.Username == "" {
		return 0, fmt.Errorf("username is empty")
	}

	balance, err := s.ch.GetUserPoints(ctx, user.ID)
	if err != nil {
		return 0, err
	}
	if balance < product.Price {
		return 0, fmt.Errorf("insufficient balance")
	}
	newBalance, err := s.ch.DecUserPoints(ctx, user.ID, product.Price)
	if err != nil {
		return 0, err
	}
	if newBalance < 0 {
		s.ch.IncUserPoints(ctx, user.ID, product.Price)
		return 0, fmt.Errorf("insufficient balance")
	}

	if product.NeedUsername() {
		if err = s.db.SetUserUsername(ctx, user.ID, user.Username); err != nil {
			return 0, err
		}
	}

	err = s.ns.Publish(ctx, worker.KeyProductClaim, worker.ProductClaimMessage{
		Product: product,
		User:    user,
		ClaimAt: time.Now(),
	})
	if err != nil {
		s.ch.IncUserPoints(ctx, user.ID, product.Price)
		return 0, err
	}
	return product.Price, nil
}

func (s *Service) ForceRewardUser(ctx context.Context, userID, points int64) error {
	if err := s.db.IncUserPoints(ctx, userID, points); err != nil {
		return err
	}
	if err := s.ch.IncUserPoints(ctx, userID, points); err != nil {
		return err
	}
	return nil
}

//

func (s *Service) ProcessUserStart(ctx context.Context, user tgu.User, startParam, ipAddress, userAgent string) error {
	exist, err := s.ch.CheckUser(ctx, user.ID)
	if exist || err != nil {
		return err
	}
	return s.ns.Publish(ctx, worker.KeyUserNew, worker.UserNewMessage{
		User:       user,
		JoinAt:     time.Now(),
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		StartParam: startParam,
	})
}

func (s *Service) ProcessUserInit(ctx context.Context, userID int64) error {
	inited, err := s.ch.CheckUserInit(ctx, userID)
	if inited || err != nil {
		return err
	}
	return s.ns.Publish(ctx, worker.KeyUserInit, worker.UserInitMessage{
		User: userID,
		Time: time.Now(),
	})
}

//

func (s *Service) GetTasksEvents(ctx context.Context, userID, fromTime int64) ([]common.TaskEvent, error) {
	events, err := s.ch.PopUserTasksEvents(ctx, userID)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(events); i++ {
		if events[i].Time < fromTime {
			copy(events[i:], events[i+1:])
			events = events[:len(events)-1]
			i--
		}
	}
	return events, err
}

func (s *Service) TaskStart(ctx context.Context, userID, taskID, subID int64, userIP, userAgent string) error {
	task, err := s.getTask(taskID)
	if err != nil {
		return err
	}
	if task.Hidden {
		return fmt.Errorf("user %d task %d hidden", userID, taskID)
	}
	flow, err := s.ch.GetUserTask(ctx, userID, taskID, subID)
	if err != nil {
		return err
	}
	now := time.Now()
	nowUnix := time.Now().Unix()

	if flow == nil {
		flow = common.NewTaskFlow(taskID, subID)
	} else if flow.Status == common.TaskClaim {
		return nil
	} else if !task.CanStart(flow, nowUnix) {
		return fmt.Errorf("can't start task: user %d task %d", userID, taskID)
	}

	if task.Type == common.TaskTappAds {
		if err = TappAdsSendClick(ctx, s.log, task.ActionTappAdsToken, userIP, userAgent, userID, subID); err != nil {
			return err
		}
	}

	err = s.ns.Publish(ctx, worker.KeyTaskStart, worker.TaskMessage{
		Time:   now,
		UserID: userID,
		TaskID: taskID,
		SubID:  subID,
	})
	if err != nil {
		return err
	}

	flow.Start = nowUnix
	flow.Status = common.TaskPending

	return s.ch.Tx(ctx, func(c cache.Cache) error {
		if err = c.SetUserTask(ctx, userID, flow); err != nil {
			return err
		}
		return c.IncTaskClicks(ctx, taskID)
	})
}

func (s *Service) TaskClaim(ctx context.Context, userID, taskID, subID int64) error {
	task, err := s.getTask(taskID)
	if err != nil {
		return err
	}
	flow, err := s.ch.GetUserTask(ctx, userID, taskID, subID)
	if err != nil {
		return err
	}
	if flow == nil {
		return fmt.Errorf("user %d task %d flow not found", userID, taskID)
	}
	if flow.Status != common.TaskClaim {
		return fmt.Errorf("user %d task %d not in state claim", userID, taskID)
	}

	err = s.ns.Publish(ctx, worker.KeyTaskClaim, worker.TaskMessage{
		Time:   time.Now(),
		UserID: userID,
		TaskID: taskID,
		SubID:  subID,
	})
	if err != nil {
		return err
	}

	if task.Interval == 0 || flow.Start+task.Interval > time.Now().Unix() {
		flow.Status = common.TaskDone
	} else {
		flow.Status = common.TaskActive
	}
	return s.ch.SetUserTask(ctx, userID, flow)
}

func (s *Service) GetUserRefs(ctx context.Context, userID int64) (refCount, refPoints int64, err error) {
	refPoints, err = s.ch.GetUserRefPoints(ctx, userID)
	if err != nil {
		return
	}
	refCount, err = s.ch.GetUserRefCount(ctx, userID)
	return
}

func (s *Service) GetUserProgress(
	ctx context.Context, user tgu.User,
	startParam, userIP, userAgent string,
) (userTasks []common.UserTask, points int64, err error) {
	tasksFlow, err := s.ch.GetUserTasks(ctx, user.ID)
	if err != nil {
		s.log.Error().Err(err).Int64("user_id", user.ID).Msg("[service] failed to get user tasks processes from cache")
		return nil, 0, err
	}

	err = s.ProcessUserStart(ctx, user, startParam, userIP, userAgent)
	if err != nil {
		s.log.Error().Err(err).
			Int64("user_id", user.ID).
			Str("start_param", startParam).
			Msg("[service] failed to process user start")
		return nil, 0, err
	}

	nowTime := time.Now().Unix()

	applyTaskFlow := func(t *common.UserTask) {
		t.Status = common.TaskActive
		if tasksFlow != nil {
			if flow, ok := tasksFlow[t.TaskKey]; ok {
				t.Start = flow.Start
				if !t.CanStart(&flow, nowTime) {
					t.Status = flow.Status
				}
			}
		}
		t.Name = parseAndGetLocale(t.Name, user.LanguageCode)
		t.Desc = parseAndGetLocale(t.Desc, user.LanguageCode)
	}

	allTasks := s.allTasks
	userTasks = make([]common.UserTask, 0, len(allTasks)+10)
	for _, ft := range allTasks {
		switch ft.Type {
		case common.TaskTappAds:
			if ft.Hidden {
				continue
			}
			adsList, err := TappAdsGetTasks(ctx, ft.ActionTappAdsToken, userIP, userAgent, user.LanguageCode, user.IsPremium, user.ID)
			if err != nil {
				s.log.Warn().Err(err).Int64("user_id", user.ID).Msg("[service] failed to get tapp_ads tasks")
			}
			for _, ads := range adsList {
				task := ads.WithUserTask(ft.UserTask)
				applyTaskFlow(&task)
				if ads.IsDone {
					if _, ok := tasksFlow[task.TaskKey]; !ok {
						continue
					}
					if task.Status != common.TaskClaim {
						task.Status = common.TaskDone
					}
				}
				userTasks = append(userTasks, task)
			}
		default:
			task := ft.UserTask
			applyTaskFlow(&task)
			if task.Status == common.TaskActive && (task.Hidden || ft.MaxClicks > 0 && ft.NowClicks > ft.MaxClicks) {
				continue
			}
			userTasks = append(userTasks, task)
		}
	}
	sortUserTasks(userTasks)

	points, err = s.ch.GetUserPoints(ctx, user.ID)
	if err != nil {
		s.log.Error().Err(err).Int64("user_id", user.ID).Msg("[service] failed to get user points from cache")
		return nil, 0, err
	}
	for i := range userTasks {
		if userTasks[i].Status == common.TaskClaim {
			points -= userTasks[i].Points
		}
	}
	if points < 0 {
		points = 0
	}
	return
}

//

func (s *Service) GetInviteMessage(ctx context.Context, user tgu.User) (msgID string, err error) {
	msgID, err = s.ch.GetInviteMessage(ctx, user.ID)
	if err != nil || msgID != "" {
		return
	}

	loc := s.cfg.Miniapp.En
	if user.LanguageCode == "ru" {
		loc = s.cfg.Miniapp.Ru
	}

	msg, err := s.bot.PrepareInviteMessage(
		ctx, user.ID,
		loc.InviteText, s.cfg.Miniapp.InviteImage, loc.InviteButton, s.cfg.Miniapp.MiniappURL,
	)
	if err != nil {
		return "", err
	}

	ttl := time.Until(time.Unix(msg.ExpirationDate-30, 0))
	if err = s.ch.SaveInviteMessage(ctx, user.ID, msg.ID, ttl); err != nil {
		return "", err
	}

	return msg.ID, nil
}

func (s *Service) GetTaskStarsInvoice(ctx context.Context, userID, taskID int64) (invoiceLink string, err error) {
	task, err := s.getTask(taskID)
	if err != nil {
		return "", err
	}
	if task.Type != common.TaskStarsDeposit {
		return "", fmt.Errorf("task %d type not %s", taskID, common.TaskStarsDeposit)
	}

	invoiceLink, err = s.ch.GetStarsInvoice(ctx, userID, taskID)
	if err != nil || invoiceLink != "" {
		return invoiceLink, err
	}

	payload := common.NewStarsInvoicePayload(userID, taskID)
	invoiceLink, err = s.bot.CreateStarsInvoice(
		ctx, task.ActionStarsAmount, payload, task.ActionStarsTitle, task.ActionStarsDesc, task.ActionStarsItem,
	)
	if err != nil {
		return "", err
	}

	const invoiceTTL = time.Minute * 5
	if err = s.ch.SaveStarsInvoice(ctx, userID, taskID, invoiceLink, invoiceTTL); err != nil {
		s.log.Error().Err(err).Int64("user_id", userID).Int64("task_id", taskID).Msg("[service] failed to save stars invoice")
	}

	return invoiceLink, nil
}

func (s *Service) GetTaskTonTransaction(ctx context.Context, userID, taskID int64) (*common.TonCTx, error) {
	task, err := s.getTask(taskID)
	if err != nil {
		return nil, err
	}
	if task.Type != common.TaskTonDeposit {
		return nil, fmt.Errorf("task %d type not %s", taskID, common.TaskTonDeposit)
	}

	txPayload, err := tonu.NewTxWithComment(
		common.NewTonTxComment(s.cfg.Ton.CommentPrefix, userID, taskID),
	)
	if err != nil {
		return nil, err
	}

	const tonInvoiceTTL = 60 * 10 // 10 min
	tonTx := &common.TonCTx{
		ValidUntil: time.Now().Unix() + tonInvoiceTTL,
		Messages: []common.TonCTxMessage{{
			Address: s.cfg.Ton.DepositWallet,
			Amount:  strconv.FormatInt(task.ActionTonAmountUnits(), 10),
			Payload: txPayload,
		}},
	}
	return tonTx, nil
}
