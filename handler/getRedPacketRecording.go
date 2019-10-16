package handler

import (
	"blockChainWallet/database"

	"github.com/henrylee2cn/faygo"
)

type RedPacketRecordingParameterModel struct {
	RedPacketDir int    `param:"<desc:0:> <in:formData> <required> <name:redPacketDir> <range:0:1> <err:>"`
	LastId       int64  `param:"<desc:ID> <in:formData> <required> <name:lastId> <err:>"`
	Sign         string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

type redPacketModel struct {
	RedPacketId     int64   `xorm:"bigint(20) 'id'" json:"redPacketId"`
	RedPacketAmount float64 `xorm:"decimal(20,3) 'red_packet_amount'" json:"redPacketType"`
	RedPacketNumber int     `xorm:"int(8) 'red_packet_number'" json:"redPacketNumber"`
	RemainingNumber int     `xorm:"int(8) 'remaining_number'" json:"remainingNumber"`
	CreateTime      int64   `xorm:"bigint(20) 'create_time'" json:"createTime"`
	CurrencyNumber  string  `xorm:"varchar(255) 'currency_number'" json:"currencyNumber"`
}

type redPacketRecordingModel struct {
	NickName        string  `xorm:"varchar(20) 'nickname'" json:"nickName"`
	RedPacketId     int64   `xorm:"bigint(20) 'red_packet_id'" json:"redPacketId"`
	RedPacketAmount float64 `xorm:"decimal(20,3) 'red_packet_amount'" json:"redPacketAmount"`
	ReceiveTime     int64   `xorm:"bigint(20) 'receive_time'" json:"receiveTime"`
	CurrencyNumber  string  `xorm:"varchar(255) 'currency_number'" json:"currencyNumber"`
}

func (rprpm *RedPacketRecordingParameterModel) Serve(ctx *faygo.Context) error {

	userId := ctx.HeaderParam("userId")

	if rprpm.RedPacketDir == 0 {

		redPackets := make([]redPacketModel, 0)
		err := database.Engine.Table("bcw_red_packet").Where("id > ? and user_id = ?", rprpm.LastId, userId).Limit(15).Find(&redPackets)
		if err == nil {
			return ctx.JSON(200, HandlerSucceed(redPackets), true)
		}
	} else if rprpm.RedPacketDir == 1 {

		redPacketRecordings := make([]redPacketRecordingModel, 0)
		err := database.Engine.Table("bcw_user_info").Alias("ui").Where("ui.id = ?", userId).Join("INNER", []string{"", "bcw_red_packet_recording", "rpr"}, "ui.phone_number = rpr.phone_number").Where("rpr.id > ?", rprpm.LastId).Limit(15).Find(&redPacketRecordings)
		if err == nil {
			return ctx.JSON(200, HandlerSucceed(redPacketRecordings), true)
		}
	} else {
		return ctx.JSON(200, HandlerCode(1001), true)
	}
	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (rprpm *RedPacketRecordingParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
