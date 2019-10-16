package handler

import (
	"blockChainWallet/database"

	"github.com/go-xorm/xorm"

	"fmt"

	"github.com/henrylee2cn/faygo"
)

type AlterParameterModel struct {
	IdentityId string `param:"<desc:> <in:formData> <name:identityId> <err:>"`
	Name       string `param:"<desc:> <in:formData> <name:name>  <err:>"`
	Sex        string `param:"<desc:> <in:formData> <name:sex> <err:>"`
	Age        string `param:"<desc:> <in:formData> <name:age> <err:>"`
	Password   string `param:"<desc:> <in:formData> <name:password> <err:>"`
	NickName   string `param:"<desc:> <in:formData> <name:nickName> <err:>"`

	Code        string `param:"<desc:> <in:formData> <name:code> <err:>"`
	CodeType    string `param:"<desc:> <in:formData> <name:codeType> <err:>"`
	PhoneNumber string `param:"<desc:> <in:formData> <name:phoneNumber> <err:>"`

	Sign string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

// Serve impletes Handler.
func (apm *AlterParameterModel) Serve(ctx *faygo.Context) error {

	userId := ctx.HeaderParam("userId")
	uniqueMark := ctx.HeaderParam("uniqueMark")
	var aUI database.UserInfo
	var flag bool
	var err error
	if apm.CodeType == "2" {
		flag, err = database.Engine.Where("phone_number = ?", apm.PhoneNumber).Get(&aUI)
	} else {
		flag, err = database.Engine.ID(userId).Get(&aUI)
	}
	if err == nil {
		if flag {
			if apm.IdentityId != "" && apm.Name != "" && apm.Code != "" && apm.CodeType != "" {
				if len(apm.IdentityId) == 18 {
					nameLen := len(apm.Name) / 3
					if nameLen > 1 && nameLen < 6 {
						if apm.CodeType == "3" {
							if phoneIsEffective(apm.PhoneNumber) {
								vCode, err := new(database.RedisConn).GetData(apm.PhoneNumber + apm.CodeType)
								if err == nil {
									if len(apm.Code) == 6 && string(vCode) == apm.Code {
										if juheAuthenticationIsMatch(apm.Name, apm.IdentityId) {
											aUI.IdentityId = apm.IdentityId
											aUI.Name = apm.Name
											err := database.SessionSubmit(func(session *xorm.Session) error {
												_, err1 := session.Table("bcw_user_info").Where("id = ?", userId).Update(map[string]interface{}{"identity_id": aUI.IdentityId, "name": aUI.Name})
												return err1
											})
											if err == nil {
												alterRedisUserInfo(userId, uniqueMark, aUI)
												return ctx.JSON(200, HandlerCode(200), true)
											}
										} else {
											return ctx.JSON(200, HandlerCode(6000), true)
										}
									} else {
										return ctx.JSON(200, HandlerCode(2003), true)
									}
								} else {
									return ctx.JSON(200, HandlerCode(2003), true)
								}
							} else {
								return ctx.JSON(400, "", true)
							}
						} else {
							return ctx.JSON(200, HandlerCode(2004), true)
						}
					} else {
						return ctx.JSON(400, "", true)
					}
				} else {
					return ctx.JSON(400, "", true)
				}
			} else if apm.Sex != "" {
				if apm.Sex == "0" {
					aUI.Sex = 0
				} else if apm.Sex == "1" {
					aUI.Sex = 1
				} else {
					return ctx.JSON(400, "", true)
				}
				err := database.SessionSubmit(func(session *xorm.Session) error {
					_, err1 := session.Table("bcw_user_info").Where("id = ?", userId).Update(map[string]interface{}{"sex": aUI.Sex})
					return err1
				})
				if err == nil {
					alterRedisUserInfo(userId, uniqueMark, aUI)
					return ctx.JSON(200, HandlerCode(200), true)
				}
			} else if apm.PhoneNumber != "" && apm.Password != "" && apm.Code != "" && apm.CodeType != "" {
				if apm.CodeType == "2" {
					if phoneIsEffective(apm.PhoneNumber) {
						vCode, err := new(database.RedisConn).GetData(apm.PhoneNumber + apm.CodeType)
						if err == nil {
							if len(apm.Code) == 6 && string(vCode) == apm.Code {
								aUI.Password = apm.Password
								err := database.SessionSubmit(func(session *xorm.Session) error {
									_, err1 := session.Table("bcw_user_info").Where("phone_number = ?", apm.PhoneNumber).Update(map[string]interface{}{"password": aUI.Password})
									return err1
								})
								fmt.Println(err)
								if err == nil {
									alterRedisUserInfo(userId, uniqueMark, aUI)
									return ctx.JSON(200, HandlerCode(200), true)
								}
							} else {
								return ctx.JSON(200, HandlerCode(2003), true)
							}
						} else {
							return ctx.JSON(200, HandlerCode(2003), true)
						}
					} else {
						return ctx.JSON(400, "", true)
					}
				} else {
					return ctx.JSON(200, HandlerCode(2004), true)
				}
			} else if apm.NickName != "" {
				nickNameLen := len(apm.NickName)
				if nickNameLen > 0 && nickNameLen < 16 {
					aUI.Nickname = apm.NickName
					err := database.SessionSubmit(func(session *xorm.Session) error {
						_, err1 := session.Table("bcw_user_info").Where("id = ?", userId).Update(map[string]interface{}{"nickname": aUI.Nickname})
						return err1
					})
					if err == nil {
						alterRedisUserInfo(userId, uniqueMark, aUI)
						return ctx.JSON(200, HandlerCode(200), true)
					}
				} else {
					return ctx.JSON(400, "", true)
				}
			} else {
				return ctx.JSON(200, HandlerCode(4000), true)
			}
		} else {
			return ctx.JSON(200, HandlerCode(2001), true)
		}
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (apm *AlterParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
