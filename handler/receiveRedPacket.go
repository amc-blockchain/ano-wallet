package handler

import (
	"blockChainWallet/database"
	"time"

	"github.com/go-xorm/xorm"

	"github.com/henrylee2cn/faygo"
)

type ReceiveRedPacketParameterModel struct {
	RedPacketId int64  `param:"<desc:> <in:formData> <required> <name:redPacketId> <err:>"`
	PhoneNumber string `param:"<desc:> <in:formData> <required> <name:phoneNumber> <len:11> <regexp:^((13[0-9])|(14[5,7,9])|(15[^4])|(18[0-9])|(17[0,1,3,5,6,7,8]))\\d{8}$> <err:>"`
}

// Serve impletes Handler.
func (rrppm *ReceiveRedPacketParameterModel) Serve(ctx *faygo.Context) error {

	flag2, err2 := database.Engine.Where("red_packet_id = ? and status = 1 and phone_number = ?", rrppm.RedPacketId, rrppm.PhoneNumber).Exist(new(database.RedPacketRecording))
	if err2 == nil {
		if flag2 {
			return ctx.JSON(200, HandlerCode(2012), true)
		} else {
			newTime := time.Now().Unix()

			var tmpRedPacket database.RedPacket
			tmpFlag, tmpErr := database.Engine.Where("id = ?", rrppm.RedPacketId).Get(&tmpRedPacket)

			var tmpRedPacketRecord database.RedPacketRecording
			flag, err := database.Engine.Where("red_packet_id = ? and status = 0", rrppm.RedPacketId).Get(&tmpRedPacketRecord)
			if err == nil && tmpErr == nil {
				if tmpFlag {
					if flag {
						if newTime <= tmpRedPacket.ExpireTime {
							err = database.SessionSubmit(func(session *xorm.Session) error {
								_, err1 := session.Table("bcw_red_packet_recording").Where("id = ?", tmpRedPacketRecord.Id).Update(map[string]interface{}{"phone_number": rrppm.PhoneNumber, "receive_time": newTime, "status": 1})
								if err1 == nil {
									_, err1 = session.Table("bcw_red_packet").Where("id = ?", tmpRedPacket.Id).Update(map[string]interface{}{"remaining_amount": tmpRedPacket.RedPacketAmount - tmpRedPacketRecord.RedPacketAmount, "remaining_number": tmpRedPacket.RedPacketNumber - 1, "finish_time": newTime})
									if err1 == nil {
										var uInfo database.UserInfo
										fFlag, fErr := session.Where("phone_number = ?", rrppm.PhoneNumber).Get(&uInfo)
										if fErr == nil {
											if !fFlag {
												var uiData database.UserInfo
												uiData.PhoneNumber = rrppm.PhoneNumber
												uiData.Nickname = rrppm.PhoneNumber
												uiData.RedPacketTotal = tmpRedPacketRecord.RedPacketAmount
												_, err1 = session.Insert(uiData)
											} else {
												_, err1 = session.Table("bcw_user_info").Where("phone_number = ?", rrppm.PhoneNumber).Update(map[string]interface{}{"red_packet_total": uInfo.RedPacketTotal + tmpRedPacketRecord.RedPacketAmount})
											}
										} else {
											return fErr
										}
									}
								}
								return err1
							})
							if err == nil {
								return ctx.JSON(200, HandlerCode(200), true)
							}
						} else {
							
							return ctx.JSON(200, HandlerCode(2009), true)
						}
					} else {
						
						return ctx.JSON(200, HandlerCode(2010), true)
					}
				} else {
					
					return ctx.JSON(200, HandlerCode(2011), true)
				}
			}
		}
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (rrppm *ReceiveRedPacketParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
