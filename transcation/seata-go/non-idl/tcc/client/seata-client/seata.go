package seata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"dubbo.apache.org/dubbo-go/v3/client"
	"github.com/dubbogo/gost/log/logger"
	"github.com/seata/seata-go/pkg/common"
	"github.com/seata/seata-go/pkg/common/net"
	"github.com/seata/seata-go/pkg/protocol/branch"
	"github.com/seata/seata-go/pkg/rm"
	"github.com/seata/seata-go/pkg/rm/tcc"
	"github.com/seata/seata-go/pkg/tm"
)

func registeBranch(t *tcc.TCCServiceProxy, ctx context.Context) error {
	// register transaction branch
	if !tm.IsTransactionOpened(ctx) {
		err := errors.New("BranchRegister error, transaction should be opened")
		logger.Errorf(err.Error())
		return err
	}
	tccContext := make(map[string]interface{}, 0)
	tccContext[common.StartTime] = time.Now().UnixNano() / 1e6
	tccContext[common.HostName] = net.GetLocalIp()
	tccContextStr, _ := json.Marshal(map[string]interface{}{
		common.ActionContext: tccContext,
	})

	branchId, err := rm.GetRMRemotingInstance().BranchRegister(branch.BranchTypeTCC, t.GetActionName(), "", tm.GetXID(ctx), string(tccContextStr), "")
	if err != nil {
		err = fmt.Errorf("BranchRegister error: %v", err.Error())
		logger.Errorf(err.Error())
		return err
	}

	actionContext := &tm.BusinessActionContext{
		Xid:        tm.GetXID(ctx),
		BranchId:   branchId,
		ActionName: t.GetActionName(),
		//ActionContext: param,
	}
	tm.SetBusinessActionContext(ctx, actionContext)
	return nil
}

func Prepare(t *tcc.TCCServiceProxy, ctx context.Context, conn *client.Connection, param ...interface{}) (resp interface{}, err error) {
	if tm.IsTransactionOpened(ctx) {
		err := registeBranch(t, ctx)
		if err != nil {
			return nil, err
		}
	}
	err = conn.CallUnary(ctx, []interface{}{1}, &resp, "Prepare")
	return
}

func CommitOrRollback(conn *client.Connection, ctx context.Context, isSuccess bool) error {
	role := *tm.GetTransactionRole(ctx)
	if role == tm.PARTICIPANT {
		// Participant has no responsibility of rollback
		logger.Debugf("Ignore Rollback(): just involved in global transaction [%s]", tm.GetXID(ctx))
		return nil
	}
	tx := &tm.GlobalTransaction{
		Xid:    tm.GetXID(ctx),
		Status: *tm.GetTxStatus(ctx),
		Role:   role,
	}
	var (
		err           error
		retry         = 10
		retryInterval = 200 * time.Millisecond
	)

	for ; retry > 0; retry-- {
		if isSuccess {
			err = tm.GetGlobalTransactionManager().Commit(ctx, tx)
			if err != nil {
				logger.Infof("transactionTemplate: commit transaction failed, error %v", err)
			}
		} else {
			err = tm.GetGlobalTransactionManager().Rollback(ctx, tx)
			if err != nil {
				logger.Infof("transactionTemplate: Rollback transaction failed, error %v", err)
			}
		}
		if err == nil {
			break
		} else {
			time.Sleep(retryInterval)
		}
	}
	return err
}

// 	businessActionCtx := tm.GetBusinessActionContext(ctx)
// 	if isSuccess {
// 		if err := conn.CallUnary(context.Background(), []interface{}{businessActionCtx}, &resp, "Commit"); err != nil {
// 			logger.Errorf("response commit", err)
// 			return
// 		}
// 	}
// 	if err := conn.CallUnary(context.Background(), []interface{}{businessActionCtx}, &resp, "Rollback"); err != nil {
// 		logger.Errorf("response Rollback", err)
// 		return
// 	}
// 	logger.Infof("get resp %#v", resp)
// 	return
// }
