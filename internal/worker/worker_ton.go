package worker

import (
	"context"
	"time"

	"github.com/avpetkun/the-prime/internal/common"
	"github.com/avpetkun/the-prime/pkg/tonu"
)

func (s *Service) listenTonTransactions(ctx context.Context) error {
	lis := &tonu.Listener{
		API: s.ton,
		Log: s.log,

		WatchAddress: s.cfg.Ton.DepositWallet,

		ProcessedTransaction: s.db.CheckTonTxByHash,
		ProcessTransaction:   s.processTonTransaction,
	}
	return lis.Start(ctx)
}

func (s *Service) processTonTransaction(ctx context.Context, transaction *tonu.Transaction) error {
	if !transaction.IsDeposit {
		return nil
	}
	log := s.log.With().Any("tx", transaction).Logger()

	userID, taskID, success := common.ParseTonTxComment(transaction.Comment)
	if !success {
		log.Warn().Msg("[ton] invalid tx comment")
		return nil
	}

	task := s.findTaskByID(taskID)
	if task == nil {
		log.Warn().Msg("[ton] task_id not found")
		return nil
	}
	if task.ActionTonAmountUnits() != transaction.TonAmount.Int64() {
		log.Warn().Msg("[ton] received invalid tx amount")
		return nil
	}

	if err := s.db.SaveTonTx(ctx, userID, taskID, transaction); err != nil {
		return err
	}
	return s.ns.Publish(ctx, KeyTaskDone, TaskMessage{
		Time:   time.Now(),
		UserID: userID,
		TaskID: taskID,
		SubID:  0,
	})
}
