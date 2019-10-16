package handler

import (
	"blockChainWallet/database"

	"blockChainWallet/purseInterface"

	"fmt"
	"strings"

	"github.com/henrylee2cn/faygo"
)

type TradeParameterModel struct {
	SendAddress     string  `param:"<desc:> <in:formData> <required> <name:sendAddress> <err:>"`
	ReceiveAddress  string  `param:"<desc:> <in:formData> <required> <name:receiveAddress> <err:>"`
	ContractAddress string  `param:"<desc:> <in:formData> <name:contractAddress> <err:>"`
	MinerCosts      float64 `param:"<desc:> <in:formData> <required> <name:minerCosts> <nonzero> <err:>"`
	TradeNumber     float64 `param:"<desc:> <in:formData> <required> <name:tradeNumber> <nonzero> <err:>"`
	TradePassword   string  `param:"<desc:> <in:formData> <required> <name:tradePassword> <err:>"`
	ChainType       string  `param:"<desc:> <in:formData> <required> <name:chainType> <err:>"`
	Sign            string  `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

// Serve impletes Handler.
func (tpm *TradeParameterModel) Serve(ctx *faygo.Context) error {

	if tpm.SendAddress == tpm.ReceiveAddress {
		return ctx.JSON(200, HandlerCode(3003), true)
	}

	aesStr, err1 := purseInterface.AesOperate(tpm.TradePassword, false)
	if err1 == nil {
		var cm database.CurrencyManagement
		flag, err := database.Engine.Where("chain_type = ? and currency_contract_address = ?", tpm.ChainType, tpm.ContractAddress).Get(&cm)
		if err == nil {
			if flag {
				tradeHash, tradeErr := purseInterface.SendTransaction(&purseInterface.Mc, tpm.SendAddress, tpm.ContractAddress, tpm.ReceiveAddress, aesStr, tpm.MinerCosts, tpm.TradeNumber, cm.ContractPrecision)
				fmt.Println(tradeHash, tradeErr)
				if tradeErr == nil {
					
					return ctx.JSON(200, HandlerSucceed(tradeHash), true)
				} else if tradeErr.Error() == "intrinsic gasRemaining too low" {
					return ctx.JSON(200, HandlerCode(3002), true)
				} else if tradeErr.Error() == "could not decrypt key with given passphrase" {
					return ctx.JSON(200, HandlerCode(3001), true)
				} else if tradeErr.Error() == "insufficient funds for gasRemaining * price + value" {
					return ctx.JSON(200, HandlerCode(3002), true)
				} else if strings.Contains(tradeErr.Error(), "invalid JSON argument 0") {
					return ctx.JSON(200, HandlerCode(3004), true)
				} else if strings.Contains(tradeErr.Error(), "unknown account") {
					return ctx.JSON(200, HandlerCode(3005), true)
				} else if tradeErr.Error() == "intrinsic gasRemaining too low" {
					return ctx.JSON(200, HandlerCode(3006), true)
				}
			} else {
				return ctx.JSON(200, HandlerCode(3000), true)
			}
		}
	} else {
		return ctx.JSON(200, HandlerCode(3001), true)
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (tpm *TradeParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
