package handler

import (
	"blockChainWallet/database"
	"blockChainWallet/purseInterface"
	"fmt"
	"strings"

	"github.com/henrylee2cn/faygo"
)

type CurrencyInfoParameterModel struct {
	Address string `param:"<desc:> <in:formData> <required> <name:address> <err:>"`
	Sign    string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

type assetsModel struct {
	CurrencyId              int64  `json:"currencyId"`
	CurrencyNumber          string `json:"currencyNumber"`
	CurrencyName            string `json:"currencyName"`
	CurrencyImgUrl          string `json:"currencyImgUrl"`
	CurrencyContractAddress string `json:"currencyContractAddress"`
	ContractPrecision       int    `json:"contractPrecision"`
}

func (cipm *CurrencyInfoParameterModel) Serve(ctx *faygo.Context) error {

	user_id := ctx.HeaderParam("userId")

	type addressManagement struct {
		database.UserAddress        `xorm:"extends"`
		database.CurrencyManagement `xorm:"extends"`
	}
	am := make([]addressManagement, 0)
	err1 := database.Engine.Table("bcw_user_address").Where("user_id = ? and currency_address = ?", user_id, cipm.Address).Join("LEFT OUTER", "bcw_currency_management", "bcw_currency_management.chain_type = bcw_user_address.chain_type").Find(&am)
	if err1 == nil {
		assetsM := make([]assetsModel, 10)
		i := -1
		for index, value := range am {
			var ua database.UserAssets
			flag2, err2 := database.Engine.Where("user_id = ? and currency_id = ?", value.UserId, value.CurrencyId).Get(&ua)
			if err2 == nil {
				if flag2 {
					assetsM[index].ContractPrecision = value.ContractPrecision
					assetsM[index].CurrencyContractAddress = value.CurrencyContractAddress
					assetsM[index].CurrencyId = value.CurrencyId
					assetsM[index].CurrencyImgUrl = value.CurrencyImageUrl
					assetsM[index].CurrencyName = value.CurrencyName

					balance := ua.CurrencyNumber

					if strings.ToLower(value.CurrencyContractAddress) == "" {
						var err error
						balance, err = purseInterface.GetBalance(&purseInterface.Mc, cipm.Address, value.CurrencyContractAddress, "latest", 4)
						if err != nil {
							balance = ua.CurrencyNumber
						}
					}
					assetsM[index].CurrencyNumber = fmt.Sprintf("%.5f", balance)
				}
			}
			i = index
		}
		return ctx.JSON(200, HandlerSucceed(assetsM[:i+1]), true)
	}
	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (cipm *CurrencyInfoParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
