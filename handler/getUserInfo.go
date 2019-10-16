package handler

import (
	"blockChainWallet/database"
	"encoding/json"
	"fmt"

	"github.com/henrylee2cn/faygo"
)

type UserInfoParameterModel struct {
	Sign string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

// Serve impletes Handler.
func (uipm *UserInfoParameterModel) Serve(ctx *faygo.Context) error {

	value, err := new(database.RedisConn).GetData(ctx.HeaderParam("userId") + ctx.HeaderParam("uniqueMark"))
	if err == nil {
		var lm LoginModel
		err = json.Unmarshal(value, &lm)
		if err == nil {
			insertValue, _ := json.Marshal(lm)
			new(database.RedisConn).InsertData(fmt.Sprintf("%d%s", lm.UserId, lm.UniqueMark), insertValue, 7*24*3600)
			return ctx.JSON(200, HandlerSucceed(lm), true)
		}
	}
	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (uipm *UserInfoParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
