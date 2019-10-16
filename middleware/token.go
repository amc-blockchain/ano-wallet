package middleware

import (
	"blockChainWallet/handler"

	"blockChainWallet/database"
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"sort"

	"github.com/henrylee2cn/faygo"
)

/*
	"", "POST", "/userRegister"
	"", "POST", "/userLogin"
	"", "POST", "/getUserInfo"
	"", "POST", "/alterUserInfo"
	"", "POST", "/feedback"
	"", "POST", "/sendVerificationCode"
	"", "POST", "/getUserNotify"
	"", "POST", "/readUserNotify"
	"", "POST", "/getAddressManagentList"
	"", "POST", "/alterAddressInfo"

	
	"", "POST", "/getCurrencyInfo"
	"", "POST", "/getTradeRecording"
	"", "POST", "/trade"
*/

/*
Token
*/
var Token = faygo.HandlerFunc(func(ctx *faygo.Context) error {

	flag := false
	path := ctx.URL().Path[1:]
	switch path {
	case "userRegister":
	case "userLogin":
	case "sendVerificationCode":
	case "userLogout":
		flag = true
	case "getUserInfo":
		flag = true
	case "alterUserInfo":
		if ctx.Param("codeType") != "2" {
			flag = true
		}
	case "feedback":
		flag = true
	case "getUserNotify":
		flag = true
	case "readUserNotify":
		flag = true
	case "getAddressManagentList":
		flag = true
	case "alterAddressManagentInfo":
		flag = true
	case "deleteAddressManagent":
		flag = true
	case "getCurrencyInfo":
		flag = true
	case "getTradeRecording":
		flag = true
	case "trade":
		flag = true
	case "exportPrivateKey":
		flag = true
	case "getRedPacketRecording":
		flag = true
	case "sendRedPacket":
		flag = true
	}
	//	fmt.Println(flag)
	if flag {
		userId := ctx.HeaderParam("userId")
		uniqueMark := ctx.HeaderParam("uniqueMark")
		flag1, err := new(database.RedisConn).ExitsData(userId + uniqueMark)
		if err == nil && !flag1 {
			ctx.JSON(200, handler.HandlerCode(9000), true)
			return errors.New("")
		}
	}

	if !argsHandler(ctx) && path != "getRedPacketInfo" && path != "receiveRedPacket" {
		ctx.JSON(200, handler.HandlerCode(1000), true)
		return errors.New("")
	}
	return nil
})

/*
 *
 */
func argsHandler(ctx *faygo.Context) (signFlag bool) {

	argsDict := make(map[string]string)
	keys := make([]string, 15)
	i := 0
	for tmpKey, tmpValue := range ctx.FormParamAll() {
		argsDict[tmpKey] = tmpValue[0]
		keys[i] = tmpKey
		i++
	}

	signFlag = false
	buf := bytes.Buffer{}
	buf.WriteString(ctx.HeaderParam("userId") + "blockChainWallet")
	sort.Strings(keys)
	for _, value := range keys {
		if value != "sign" && argsDict[value] != "" {
			buf.WriteString(argsDict[value])
			buf.WriteString("#$")
		}
	}
	hashValue := sha1.New()
	hashValue.Write(buf.Bytes())
	hashStr := fmt.Sprintf("%x", hashValue.Sum(nil))

	captial, ok := argsDict["sign"]
	if ok {
		if captial == hashStr {
			signFlag = true
		}
	}
	//	fmt.Println(captial, "-", hashStr)
	return signFlag
	//	return true
}
