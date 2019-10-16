package handler

import (
	"blockChainWallet/database"

	"github.com/go-xorm/xorm"

	"blockChainWallet/purseInterface"

	"errors"
	"strings"

	"github.com/henrylee2cn/faygo"
)

type RegisterParameterModel struct {
	NickName      string `param:"<desc:> <in:formData> <required> <name:nickName> <len:1:15> <err:>"`
	PhoneNumber   string `param:"<desc:> <in:formData> <required> <name:phoneNumber> <len:11> <regexp:^((13[0-9])|(14[5,7,9])|(15[^4])|(18[0-9])|(17[0,1,3,5,6,7,8]))\\d{8}$> <err:>"`
	Code          string `param:"<desc:> <in:formData> <required> <name:code> <len:6> <err:>"`
	CodeType      string `param:"<desc:> <in:formData> <required> <name:codeType> <len:1> <regexp:^([1,2,3])> <err:>"`
	LoginPassword string `param:"<desc:> <in:formData> <required> <name:loginPassword> <err:>"`
	TradePassword string `param:"<desc:> <in:formData> <required> <name:tradePassword> <err:>"`
	Sign          string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

// Serve impletes Handler.
func (rpm *RegisterParameterModel) Serve(ctx *faygo.Context) error {

	vCode, err := new(database.RedisConn).GetData(rpm.PhoneNumber + rpm.CodeType)
	if err == nil {
		if string(vCode) == rpm.Code {
			uInfo := database.UserInfo{
				MachineType:   ctx.HeaderParam("machineType"),
				PushId:        ctx.HeaderParam("pushId"),
				PhoneNumber:   rpm.PhoneNumber,
				Password:      rpm.LoginPassword,
				TradePassword: rpm.TradePassword,
				Nickname:      rpm.NickName,
			}

			//			var tmpUI database.UserInfo
			//database.Engine.Where("phone_number = ?", rpm.PhoneNumber).Get(&tmpUI)
			//Table("bcw_user_info")
			existFlag, existErr := database.Engine.Where("phone_number = ?", rpm.PhoneNumber).Exist(&database.UserInfo{})
			if existErr == nil {
				if existFlag {
					
					existFlag, existErr = database.Engine.Where("phone_number = ? and password = ?", rpm.PhoneNumber, "").Exist(&database.UserInfo{})
					if existErr == nil {
						if existFlag {
							
							if rpm.addUserInfo(ctx, uInfo, false) == nil {
								return ctx.JSON(200, HandlerCode(200), true)
							}
						} else {
							return ctx.JSON(200, HandlerCode(2000), true)
						}
					}
				} else {
					
					if rpm.addUserInfo(ctx, uInfo, true) == nil {
						return ctx.JSON(200, HandlerCode(200), true)
					}
				}
			}
		} else {
			return ctx.JSON(200, HandlerCode(2003), true)
		}
	} else {
		return ctx.JSON(200, HandlerCode(2003), true)
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

func (rpm *RegisterParameterModel) addUserInfo(ctx *faygo.Context, uInfo database.UserInfo, flag bool) error {
	tp, tmpErr := purseInterface.AesOperate(rpm.TradePassword, false)
	if tmpErr == nil {
		newAddr, err := purseInterface.CreateChainAddress(&purseInterface.Mc, tp)
		if err == nil {
			err = database.SessionSubmit(func(session *xorm.Session) error {
				var err1 error
				if flag {
					_, err = session.Insert(uInfo)
				} else {
					_, err = session.Table("bcw_user_info").Where("phone_number = ?", rpm.PhoneNumber).Update(map[string]interface{}{"machine_type": ctx.HeaderParam("machineType"), "push_id": ctx.HeaderParam("pushId"), "password": rpm.LoginPassword, "trade_password": rpm.TradePassword, "nickname": rpm.NickName})
				}
				if err1 == nil {
					var tmpUI database.UserInfo
					var flag bool
					flag, err1 = session.Where("phone_number = ?", rpm.PhoneNumber).Get(&tmpUI)
					if err1 == nil {
						if flag {
							uAddr := database.UserAddress{
								UserId:          tmpUI.Id,
								AddressName:     "MOAC",
								CurrencyAddress: strings.ToLower(newAddr),
								ChainType:       "MOAC",
							}
							_, err1 = session.Insert(uAddr)
							if err1 == nil {
								var cm []database.CurrencyManagement
								session.Find(&cm)
								for _, value := range cm {
									uA := database.UserAssets{
										UserId:         tmpUI.Id,
										CurrencyId:     value.CurrencyId,
										CurrencyNumber: 0,
									}
									_, err1 = session.Insert(uA)
									if err1 != nil {
										break
									}
								}
							}
						} else {
							err1 = errors.New("")
						}
					}
				}
				return err1
			})
			return err
		}
	}
	return tmpErr
}

// Doc returns the API's note, result or parameters information.
func (rpm *RegisterParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
