package base

import (
	"context"
	"github.com/dcalsky/gogo/common/logid"

	"github.com/cloudwego/hertz/pkg/app"
)

const (
	AccountIdKey    = "ACCOUNT_ID"
	AccountEmailKey = "ACCOUNT_EMAIL"
)

type AccountId string

type Args struct {
	TraceInfo
	AccountId    AccountId
	AccountEmail string
}

type TraceInfo struct {
	RequestId string
}

func GetArgs(ctx context.Context, c *app.RequestContext) Args {
	out := Args{
		TraceInfo: TraceInfo{
			RequestId: logid.GetLogId(ctx),
		},
		AccountId:    AccountId(c.GetString(AccountIdKey)),
		AccountEmail: c.GetString(AccountEmailKey),
	}
	return out
}

func SetArgsAccountId(c *app.RequestContext, accountId AccountId) {
	c.Set(AccountIdKey, accountId)
}

func SetArgsAccountEmail(c *app.RequestContext, email string) {
	c.Set(AccountEmailKey, email)
}
