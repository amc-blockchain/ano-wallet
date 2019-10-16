package handler

import (
	"blockChainWallet/database"

	"github.com/henrylee2cn/faygo"
)

type AddressManagentParameterModel struct {
	LastId int64  `param:"<desc:> <in:formData> <required> <name:lastId> <err:>"`
	Sign   string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

type AddressManagentModel struct {
	AddressId    int64  `json:"'addressId'"`
	AddressDesc  string `json:"'addressDesc'"`
	MatchAddress string `json:"'matchAddress'"`
}

func (ampm *AddressManagentParameterModel) Serve(ctx *faygo.Context) error {

	uam := make([]database.UserAddressManagement, 0)
	err1 := database.Engine.Where("id > ? and user_id = ?", ampm.LastId, ctx.HeaderParam("userId")).Desc("id").Limit(15).Find(&uam)
	if err1 == nil {
		amModel := make([]AddressManagentModel, 15)
		i := -1
		for index, value := range uam {
			amModel[index].AddressId = value.Id
			amModel[index].AddressDesc = value.AddressDesc
			amModel[index].MatchAddress = value.MatchAddress
			i = index
		}
		return ctx.JSON(200, HandlerSucceed(amModel[:i+1]), true)
	}
	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (ampm *AddressManagentParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
