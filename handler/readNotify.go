
package handler

import (
	"blockChainWallet/database"

	"github.com/go-xorm/xorm"

	"github.com/henrylee2cn/faygo"
)

type ReadNotifyParameterModel struct {
	NotifyId int64  `param:"<desc:> <in:formData> <required> <name:notifyId> <err:>"`
	Sign     string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

// Serve impletes Handler.
func (rnpm *ReadNotifyParameterModel) Serve(ctx *faygo.Context) error {

	var tmpUN database.UserNotify
	flag, err := database.Engine.ID(rnpm.NotifyId).Get(&tmpUN)
	if err == nil {
		if flag {
			tmpUN.NotifyReadFlag = 1
			err = database.SessionSubmit(func(session *xorm.Session) error {
				_, err1 := session.ID(rnpm.NotifyId).Update(tmpUN)
				return err1
			})
			if err == nil {
				return ctx.JSON(200, HandlerCode(200), true)
			}
		} else {
			return ctx.JSON(200, HandlerCode(4001), true)
		}
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (rnpm *ReadNotifyParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
