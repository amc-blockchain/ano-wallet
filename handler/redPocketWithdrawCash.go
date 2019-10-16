package handler

import (
	"blockChainWallet/database"

	"blockChainWallet/purseInterface"

	"github.com/go-xorm/xorm"
	"github.com/henrylee2cn/faygo"
)

type WithdrawCashParameterModel struct {
	ToAddr         string  `param:"<desc:> <in:formData> <required> <name:toAddr> <err:>"`
	RedPacketTotal float64 `param:"<desc:> <in:formData> <required> <name:redPacketTotal> <err:>"`
	Sign           string  `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

// Serve impletes Handler.
func (wcpm *WithdrawCashParameterModel) Serve(ctx *faygo.Context) error {

	userId := ctx.HeaderParam("userId")
	uniqueMark := ctx.HeaderParam("uniqueMark")

	var uInfo database.UserInfo
	flag1, err1 := database.Engine.Where("id = ?", userId).Get(&uInfo)
	if err1 == nil {
		if flag1 && uInfo.RedPacketTotal >= wcpm.RedPacketTotal {
			var cm database.CurrencyManagement
			flag, err := database.Engine.Where("chain_type = MOAC and currency_name = PIG").Get(&cm)
			if err == nil {
				if flag {
					
					fromAddr := cm.RedPacketAddr
					contractAddr := cm.CurrencyContractAddress
					aesStr, aesErr := purseInterface.AesOperate("4hZ9pXENaa97wEL5nAZYWQ==", false)
					gasTotal := 0.001
					precision := cm.ContractPrecision
					if aesErr != nil {
						return ctx.JSON(200, HandlerCode(0), true)
					}
					_, hashErr := purseInterface.SendTransaction(&purseInterface.Mc, fromAddr, contractAddr, wcpm.ToAddr, aesStr, gasTotal, wcpm.RedPacketTotal, precision)
					if hashErr == nil {
						seErr := database.SessionSubmit(func(session *xorm.Session) (err2 error) {
							_, err2 = session.Table("bcw_user_info").Where("id = ?", userId).Update(map[string]interface{}{"red_packet_total": uInfo.RedPacketTotal - wcpm.RedPacketTotal})
							return err2
						})
						if seErr == nil {
							uInfo.RedPacketTotal = uInfo.RedPacketTotal - wcpm.RedPacketTotal
							alterRedisUserInfo(userId, uniqueMark, uInfo)
							return ctx.JSON(200, HandlerCode(200), true)
						}
					} else {
						return ctx.JSON(200, HandlerCode(7000), true)
					}
				} else {
					return ctx.JSON(200, HandlerCode(3007), true)
				}
			}
		} else {
			return ctx.JSON(200, HandlerCode(2001), true)
		}
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (wcpm *WithdrawCashParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
