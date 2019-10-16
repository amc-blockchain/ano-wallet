package handler

import (
	"blockChainWallet/database"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"blockChainWallet/purseInterface"

	"github.com/henrylee2cn/faygo"
)

type LoginParameterModel struct {
	PhoneNumber   string `param:"<desc:> <in:formData> <required> <name:phoneNumber> <len:11> <regexp:^((13[0-9])|(14[5,7,9])|(15[^4])|(18[0-9])|(17[0,1,3,5,6,7,8]))\\d{8}$> <err:>"`
	LoginPassword string `param:"<desc:> <in:formData> <required> <name:loginPassword> <err:>"`
	Sign          string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

type LoginModel struct {
	UniqueMark     string  `json:"uniqueMark"`     
	UserId         int64   `json:userId`           
	PhoneNumber    string  `json:"phoneNumber"`    
	IdentityId     string  `json:"identityId"`     
	Name           string  `json:"name"`           
	Sex            int     `json:"sex"`           
	Age            int     `json:"age"`           
	Nickname       string  `json:"nickname"`       
	Address        string  `json:"address"`        
	JsonStr        string  `json:"jsonStr"`        
	AddressName    string  `json:"addressName"`    
	ChainType      string  `json:"chainType"`      
	RedPacketTotal float64 `json:"redPacketTotal"`
}

// Serve impletes Handler.
func (lpm *LoginParameterModel) Serve(ctx *faygo.Context) error {

	var reLM LoginModel
	var tmpUI database.UserInfo
	flag, err := database.Engine.Where("phone_number = ?", lpm.PhoneNumber).Get(&tmpUI)
	if err == nil {
		if flag {
			
			if tmpUI.Password == lpm.LoginPassword {
				reLM.UserId = tmpUI.Id
				reLM.PhoneNumber = tmpUI.PhoneNumber
				reLM.IdentityId = tmpUI.IdentityId
				reLM.Name = tmpUI.Name
				reLM.Sex = tmpUI.Sex
				reLM.Age = tmpUI.Age
				reLM.Nickname = tmpUI.Nickname
				reLM.RedPacketTotal = tmpUI.RedPacketTotal

				type addressManagement struct {
					database.UserAddress        `xorm:"extends"`
					database.CurrencyManagement `xorm:"extends"`
				}
				am := make([]addressManagement, 0)
				err1 := database.Engine.Table("bcw_user_address").Where("user_id = ?", reLM.UserId).Join("LEFT OUTER", "bcw_currency_management", "bcw_currency_management.chain_type = bcw_user_address.chain_type").Find(&am)
				if err1 == nil {
					for index, value := range am {
						if index == 0 {
							reLM.Address = value.CurrencyAddress
							reLM.JsonStr, _ = purseInterface.GetJsonStr(&purseInterface.Mc, value.CurrencyAddress)
							reLM.AddressName = value.AddressName
							reLM.ChainType = value.UserAddress.ChainType
						}
					}
					ts := strconv.FormatInt(time.Now().Unix(), 10)
					reLM.UniqueMark = ts
					insertValue, _ := json.Marshal(reLM)
					new(database.RedisConn).InsertData(fmt.Sprintf("%d%s", reLM.UserId, reLM.UniqueMark), insertValue, 7*24*3600)
					return ctx.JSON(200, HandlerSucceed(reLM), true)
				}
			} else {
				return ctx.JSON(200, HandlerCode(1003), true)
			}
		} else {
			
			return ctx.JSON(200, HandlerCode(2001), true)
		}
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (lpm *LoginParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
