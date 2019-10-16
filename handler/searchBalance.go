package handler

import (
	"blockChainWallet/purseInterface"

	"github.com/henrylee2cn/faygo"
)

type SearchParameterModel struct {
	Address         string `param:"<desc:> <in:formData> <required> <name:Address> <err:>"`
	ContractAddress string `param:"<desc:> <in:formData> <name:contractAddress> <err:>"`
	Sign            string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

type SearchModel struct {
	Balance float64 `json:"balance"` 
}

// Serve impletes Handler.
func (spm *SearchParameterModel) Serve(ctx *faygo.Context) error {

	var reSM SearchModel
	var err error
	reSM.Balance, err = purseInterface.GetBalance(&purseInterface.Mc, spm.Address, spm.ContractAddress, "latest", 4)
	if err != nil {
		return ctx.JSON(200, HandlerSucceed(reSM), true)
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (spm *SearchParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
