package handler

import (
	"blockChainWallet/database"

	"github.com/go-xorm/xorm"

	"github.com/henrylee2cn/faygo"
)

type DeleteAddressParameterModel struct {
	AddressId int64  `param:"<desc:> <in:formData> <required> <name:addressId> <err:>"`
	Sign      string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

// Serve impletes Handler.
func (dapm *DeleteAddressParameterModel) Serve(ctx *faygo.Context) error {

	err := database.SessionSubmit(func(session *xorm.Session) error {
		_, err1 := session.ID(dapm.AddressId).Delete(&database.UserAddressManagement{})
		return err1
	})
	if err == nil {
		return ctx.JSON(200, HandlerCode(200), true)
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (dapm *DeleteAddressParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
