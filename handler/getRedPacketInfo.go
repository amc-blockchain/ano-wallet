package handler

import (
	"blockChainWallet/database"

	"github.com/henrylee2cn/faygo"
)

type RedPacketInfoParameterModel struct {
	RedPacketId int64 `param:"<desc:> <in:formData> <required> <name:redPacketId> <err:>"`
}

type redPacketInfoModel struct {
	NickName        string  `xorm:"varchar(20) 'nickname'" json:"nickName"`
	RedPacketDesc   string  `xorm:"varchar(255) 'red_packet_desc'" json:"redPacketDesc"`
	RedPacketAmount float64 `xorm:"decimal(20,3) 'red_packet_amount'" json:"redPacketAmount"`
	RedPacketNumber int     `xorm:"int(8) 'red_packet_number'" json:"redPacketNumber"`
	RemainingNumber int     `xorm:"int(8) 'remaining_number'" json:"remainingNumber"`
	CurrencyNumber  string  `xorm:"varchar(255) 'currency_number'" json:"currencyNumber"`

	ReceiveRecording []redPacketReceiveRecordingModel `json:"receiveRecording"`
}

type redPacketReceiveRecordingModel struct {
	ReceiveRecordId string  `xorm:"int(20) 'id'" json:"receiveRecordId"`
	NickName        string  `xorm:"varchar(20) 'nickname'" json:"nickName"`
	RedPacketAmount float64 `xorm:"decimal(20,3) 'red_packet_amount'" json:"redPacketAmount"`
	ReceiveTime     int64   `xorm:"bigint(20) 'receive_time'" json:"receiveTime"`
	CurrencyNumber  string  `xorm:"varchar(255) 'currency_number'" json:"currencyNumber"`
}

func (rpipm *RedPacketInfoParameterModel) Serve(ctx *faygo.Context) error {

	rpiModel := new(redPacketInfoModel)
	flag, err := database.Engine.Table("bcw_red_packet").Alias("rp").Where("rp.id = ?", rpipm.RedPacketId).Join("INNER", []string{"bcw_user_info", "ui"}, "rp.user_id = ui.id").Get(rpiModel)
	if err == nil {
		if flag {
			rprrModels := make([]redPacketReceiveRecordingModel, 0)
			err1 := database.Engine.Table("bcw_red_packet_recording").Alias("rpr").Where("red_packet_id = ?", rpipm.RedPacketId).Join("INNER", []string{"bcw_user_info", "ui"}, "rpr.phone_number = ui.phone_number").Find(&rprrModels)
			if err1 == nil {
				rpiModel.ReceiveRecording = rprrModels
				return ctx.JSON(200, HandlerSucceed(rpiModel), true)
			}
		} else {
			return ctx.JSON(200, HandlerCode(2005), true)
		}
	}

	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (rpipm *RedPacketInfoParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
