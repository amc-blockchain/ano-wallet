package handler

import (
	"blockChainWallet/database"

	"github.com/henrylee2cn/faygo"
)

type LogoutParameterModel struct {
	Sign string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

// Serve impletes Handler.
func (lpm *LogoutParameterModel) Serve(ctx *faygo.Context) error {

	userId := ctx.HeaderParam("userId")
	uniqueMark := ctx.HeaderParam("uniqueMark")

	err := new(database.RedisConn).DeleteData(userId + uniqueMark)
	if err == nil {
		return ctx.JSON(200, HandlerCode(200), true)
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (lpm *LogoutParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
