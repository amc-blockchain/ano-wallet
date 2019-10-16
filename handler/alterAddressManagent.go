package handler

import (
	"blockChainWallet/database"
	"strconv"

	"github.com/go-xorm/xorm"
	"github.com/henrylee2cn/faygo"
)

type AlterAddressInfoParameterModel struct {
	AddressId    int64  `param:"<desc:> <in:formData> <required> <name:addressId> <err:>"`
	AddressDesc  string `param:"<desc:> <in:formData> <required> <name:addressDesc> <err:>"`
	MatchAddress string `param:"<desc:> <in:formData> <required> <name:matchAddress> <err:>"`
	Sign         string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

// Serve impletes Handler.
func (aaipm *AlterAddressInfoParameterModel) Serve(ctx *faygo.Context) error {

	var tmpUAM database.UserAddressManagement

	if aaipm.AddressId == 0 {
		//add
		tmpUAM.AddressDesc = aaipm.AddressDesc
		tmpUAM.MatchAddress = aaipm.MatchAddress
		tmpUAM.UserId, _ = strconv.ParseInt(ctx.HeaderParam("userId"), 10, 64)
		err := database.SessionSubmit(func(session *xorm.Session) error {
			_, err1 := session.Insert(tmpUAM)
			return err1
		})
		if err == nil {
			return ctx.JSON(200, HandlerCode(200), true)
		}
	} else {
		//modify
		flag, err := database.Engine.ID(aaipm.AddressId).Get(&tmpUAM)
		if err == nil {
			if flag {
				tmpUAM.AddressDesc = aaipm.AddressDesc
				tmpUAM.MatchAddress = aaipm.MatchAddress
				err = database.SessionSubmit(func(session *xorm.Session) error {
					_, err1 := session.ID(aaipm.AddressId).Update(tmpUAM)
					return err1
				})
				if err == nil {
					return ctx.JSON(200, HandlerCode(200), true)
				}
			} else {
				return ctx.JSON(200, HandlerCode(4001), true)
			}
		}
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (aaipm *AlterAddressInfoParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
