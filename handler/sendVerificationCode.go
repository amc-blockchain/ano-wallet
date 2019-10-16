
package handler

import (
	"blockChainWallet/database"
	"fmt"
	"math/rand"
	"time"

	"github.com/henrylee2cn/faygo"
)

type SendVerCodeParameterModel struct {
	PhoneNumber string `param:"<desc:> <in:formData> <required> <name:phoneNumber> <len:11> <regexp:^((13[0-9])|(14[5,7,9])|(15[^4])|(18[0-9])|(17[0,1,3,5,6,7,8]))\\d{8}$> <err:>"`
	CodeType    string `param:"<desc:> <in:formData> <required> <name:codeType> <len:1> <regexp:^([1,2,3])> <err:>"`
	Sign        string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

// Serve impletes Handler.
func (svcpm *SendVerCodeParameterModel) Serve(ctx *faygo.Context) error {


	vcode := fmt.Sprintf("%06v", (rand.New(rand.NewSource(time.Now().UnixNano()))).Int31n(1000000))
	if /*sendTextMessage(svcpm.PhoneNumber, vcode)*/ true { //
		vcode = "123456"
		err := new(database.RedisConn).InsertData(svcpm.PhoneNumber+svcpm.CodeType, vcode, 5*60)
		fmt.Println(err)
		if err == nil {
			return ctx.JSON(200, HandlerCode(200), true)
		}
	} else {
		return ctx.JSON(200, HandlerCode(2002), true)
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (svcpm *SendVerCodeParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
