package handler

import (
	"blockChainWallet/database"

	"github.com/henrylee2cn/faygo"
)

type NotifyParameterModel struct {
	LastId int64  `param:"<desc:> <in:formData> <required> <name:lastId> <err:>"`
	Sign   string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

type NotifyModel struct {
	NotifyId  int64  `json:"notifyId"`
	ReadFlag  int    `json:"readFlag"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	NotifyUrl string `json:"notifyUrl"`
	Timestamp int64  `json:"timestamp"`
}

func (npm *NotifyParameterModel) Serve(ctx *faygo.Context) error {

	un := make([]database.UserNotify, 0)
	err1 := database.Engine.Where("id > ? and user_id = ?", npm.LastId, ctx.HeaderParam("userId")).Desc("id").Limit(15).Find(&un)
	if err1 == nil {
		nModel := make([]NotifyModel, 15)
		i := -1
		for index, value := range un {
			nModel[index].NotifyId = value.Id
			nModel[index].ReadFlag = value.NotifyReadFlag
			nModel[index].Title = value.NotifyTitle
			nModel[index].Content = value.NotifyContent
			nModel[index].NotifyUrl = value.NotifyUrl
			nModel[index].Timestamp = value.NotifyTimestamp
			i = index
		}
		return ctx.JSON(200, HandlerSucceed(nModel[:i+1]), true)
	}
	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (npm *NotifyParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
