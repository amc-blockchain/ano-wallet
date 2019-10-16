package handler

import (
	"encoding/hex"

	"blockChainWallet/purseInterface"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethMath "github.com/ethereum/go-ethereum/common/math"
	"github.com/henrylee2cn/faygo"
)

type ExPrivateKeyParameterModel struct {
	Address         string `param:"<desc:> <in:formData> <required> <name:address> <err:>"`
	JsonStr         string `param:"<desc:> <in:formData> <required> <name:jsonStr> <err:>"`
	AddressPassword string `param:"<desc:> <in:formData> <required> <name:addressPassword> <err:>"`
	Sign            string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

type privateKeyModel struct {
	Key string `json:"privateKey"`
}

// Serve impletes Handler.
func (epkpm *ExPrivateKeyParameterModel) Serve(ctx *faygo.Context) error {

	//	commandStr := "find /root/.moac/keystore -name *" + epkpm.Address[2:]
	//	pathStr, err := execCommand(commandStr)
	//	if err == nil && len(pathStr) > 2 {

	//		pathStr = "cat " + pathStr[:len(pathStr)-1]
	//		jsonStr, err1 := execCommand(pathStr)
	jsonStr := epkpm.JsonStr
	//	if err1 == nil {
	aesStr, err2 := purseInterface.AesOperate(epkpm.AddressPassword, false)
	if err2 == nil {
		var pkModel privateKeyModel
		storeKey, err := keystore.DecryptKey([]byte(jsonStr), aesStr)
		if err == nil {
			pkModel.Key = hex.EncodeToString(ethMath.PaddedBigBytes(storeKey.PrivateKey.D, storeKey.PrivateKey.Params().BitSize/8))
			return ctx.JSON(200, HandlerSucceed(pkModel), true)
		} else {
			return ctx.JSON(200, HandlerCode(3001), true)
		}
	} else {
		return ctx.JSON(200, HandlerCode(3001), true)
	}
	//	}
	//	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (epkpm *ExPrivateKeyParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
