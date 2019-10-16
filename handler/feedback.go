package handler

import (
	"blockChainWallet/database"

	"github.com/henrylee2cn/faygo"

	"strconv"

	"github.com/go-xorm/xorm"
)

type FeedbackParameterModel struct {
	Opinion string `param:"<desc:> <in:formData> <required> <name:opinion> <err:>"`
	Contact string `param:"<desc:> <in:formData> <required> <name:contact> <err:>"`
	Sign    string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

func (fpm *FeedbackParameterModel) Serve(ctx *faygo.Context) error {

	user_id, err := strconv.ParseInt(ctx.HeaderParam("userId"), 10, 64)
	if err == nil {
		uf := database.UserFeedback{
			UserId:         user_id,
			OpinionContact: fpm.Contact,
			OpinionContent: fpm.Opinion,
		}
		err := database.SessionSubmit(func(session *xorm.Session) error {
			_, err1 := session.Insert(uf)
			return err1
		})
		if err == nil {
			return ctx.JSON(200, HandlerCode(200), true)
		}
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (fpm *FeedbackParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
