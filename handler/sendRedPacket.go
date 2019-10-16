
package handler

import (
	"blockChainWallet/database"
	"strings"

	"math/rand"
	"strconv"
	"time"

	"blockChainWallet/purseInterface"

	"github.com/go-xorm/xorm"
	"github.com/henrylee2cn/faygo"
)

type SendRedPacketParameterModel struct {
	SendAddress     string  `param:"<desc:> <in:formData> <required> <name:sendAddress> <err:>"`
	RedPacketTotal  float64 `param:"<desc:> <in:formData> <required> <name:redPacketTotal> <err:>"`
	RedPacketNumber int64   `param:"<desc:> <in:formData> <required> <name:redPacketNumber> <range:1:10000> <err:>"`
	RedPacketDesc   string  `param:"<desc:> <in:formData> <required> <name:redPacketDesc> <len:1:256> <err:>"`
	CurrencyId      string  `param:"<desc:> <in:formData> <required> <name:currencyId> <err:>"`
	TradePassword   string  `param:"<desc:> <in:formData> <required> <name:tradePassword> <err:>"`
	Sign            string  `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

// Serve impletes Handler.
func (srppm *SendRedPacketParameterModel) Serve(ctx *faygo.Context) error {

	userId := ctx.HeaderParam("userId")

	var cm database.CurrencyManagement
	tmpFlag1, tmpErr1 := database.Engine.Where("id = ?", srppm.CurrencyId).Get(&cm)
	if tmpErr1 != nil {
		return ctx.JSON(200, HandlerCode(0), true)
	}
	if !tmpFlag1 {
		return ctx.JSON(200, HandlerCode(3000), true)
	}

	
	to := cm.RedPacketAddr
	contractAddr := cm.CurrencyContractAddress
	aesStr, aesErr := purseInterface.AesOperate(srppm.TradePassword, false)
	gasTotal := 0.001
	precision := cm.ContractPrecision
	if aesErr != nil {
		return ctx.JSON(200, HandlerCode(0), true)
	}
	_, hashErr := purseInterface.SendTransaction(&purseInterface.Mc, cm.RedPacketAddr, contractAddr, to, aesStr, gasTotal, srppm.RedPacketTotal, precision)
	if hashErr == nil {
		
		var ua database.UserAssets
		tmpFlag, tmpErr := database.Engine.Where("user_id = ? and currency_id = ?", userId, srppm.CurrencyId).Get(&ua)
		if tmpErr != nil {
			return ctx.JSON(200, HandlerCode(0), true)
		}
		if tmpFlag {
			if ua.CurrencyNumber < srppm.RedPacketTotal {
				return ctx.JSON(200, HandlerCode(2008), true)
			}
		} else {
			return ctx.JSON(200, HandlerCode(3000), true)
		}

		creatTime := time.Now().Unix()
		rand.Seed(creatTime)

		total := srppm.RedPacketTotal
		number := srppm.RedPacketNumber

		var min float64 = 0.01
		var max float64

		var i int64 = 1

		if total < 1 {
			
			return ctx.JSON(200, HandlerCode(2006), true)
		}
		if total/float64(number) < 0.01 {
			
			return ctx.JSON(200, HandlerCode(2007), true)
		}

		var redPacket database.RedPacket
		redPacket.UserId, _ = strconv.ParseInt(userId, 10, 64)
		redPacket.RedPacketType = 1
		redPacket.PayType = 1
		redPacket.RedPacketAmount = total
		redPacket.RemainingAmount = total
		redPacket.RedPacketNumber = number
		redPacket.RemainingNumber = number
		redPacket.RedPacketDesc = srppm.RedPacketDesc
		redPacket.CreateTime = creatTime
		redPacket.FinishTime = 0
		redPacket.ExpireTime = creatTime + 3600*24
		redPacket.HasRefund = 0
		redPacket.CurrencyNumber = cm.CurrencyName
		err := database.SessionSubmit(func(session *xorm.Session) error {

			_, err1 := session.Insert(redPacket)
			flag, err2 := session.Where("user_id = ? and create_time = ?", userId, creatTime).Get(&redPacket)
			if err1 != nil && err2 != nil && !flag {
				return err1
			}

			redPackets := make([]database.RedPacketRecording, number)

			for ; i < number; i++ {
				max = total - min*(float64(number)-1)
				k := int64(number-1) / 2
				if number-i <= 2 {
					k = number - i
				}
				max = max / float64(k)
				tmpMoney := int64(min*100 + rand.Float64()*(max*100-min*100+1))
				money := float64(tmpMoney) / 100

				total = total - money

				var redPacketOne database.RedPacketRecording
				redPacketOne.RedPacketId = redPacket.Id
				redPacketOne.RedPacketAmount = money
				redPacketOne.PhoneNumber = ""
				redPacketOne.CreateTime = creatTime
				redPacketOne.ReceiveTime = 0
				redPacketOne.Status = 0
				redPacketOne.RedPacketType = redPacket.RedPacketType
				redPacketOne.CurrencyNumber = redPacket.CurrencyNumber
				redPackets[i-1] = redPacketOne
			}

			if i == number {
				var redPacketOne database.RedPacketRecording
				redPacketOne.RedPacketId = redPacket.Id
				redPacketOne.RedPacketAmount = total
				redPacketOne.PhoneNumber = ""
				redPacketOne.CreateTime = creatTime
				redPacketOne.ReceiveTime = 0
				redPacketOne.Status = 0
				redPacketOne.RedPacketType = redPacket.RedPacketType
				redPacketOne.CurrencyNumber = redPacket.CurrencyNumber
				redPackets[i-1] = redPacketOne
				_, err1 = session.Insert(srppm.sliceOutOfOrder(redPackets))
				if err1 != nil {
					return err1
				}
			}

			_, err1 = session.Table("bcw_user_assets").Where("user_id = ? and currency_id = ?", userId, srppm.CurrencyId).Update(map[string]interface{}{"currency_number": ua.CurrencyNumber - srppm.RedPacketTotal})
			if err1 != nil {
				return err1
			}

			return nil
		})
		if err == nil {
			return ctx.JSON(200, HandlerCode(200), true)
		}
	} else if hashErr.Error() == "intrinsic gasRemaining too low" {
		return ctx.JSON(200, HandlerCode(3002), true)
	} else if hashErr.Error() == "could not decrypt key with given passphrase" {
		return ctx.JSON(200, HandlerCode(3001), true)
	} else if hashErr.Error() == "insufficient funds for gasRemaining * price + value" {
		return ctx.JSON(200, HandlerCode(3002), true)
	} else if strings.Contains(hashErr.Error(), "invalid JSON argument 0") {
		return ctx.JSON(200, HandlerCode(3004), true)
	} else if strings.Contains(hashErr.Error(), "unknown account") {
		return ctx.JSON(200, HandlerCode(3005), true)
	} else if hashErr.Error() == "intrinsic gasRemaining too low" {
		return ctx.JSON(200, HandlerCode(3006), true)
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

func (srppm *SendRedPacketParameterModel) sliceOutOfOrder(inData []database.RedPacketRecording) []database.RedPacketRecording {
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(inData)
	for i := l - 1; i > 0; i-- {
		r := rr.Intn(i)
		inData[r], inData[i] = inData[i], inData[r]
	}
	return inData
}

// Doc returns the API's note, result or parameters information.
func (srppm *SendRedPacketParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
