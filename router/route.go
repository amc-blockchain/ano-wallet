package router

import (
	"blockChainWallet/handler"

	"github.com/henrylee2cn/faygo"
)

// Route register router in a tree style.
func Route(frame *faygo.Framework) {

	frame.Route(
		frame.NewNamedAPI("", "POST", "/userRegister", &handler.RegisterParameterModel{}),
		frame.NewNamedAPI("", "POST", "/userLogin", &handler.LoginParameterModel{}),
		frame.NewNamedAPI("", "POST", "/userLogout", &handler.LogoutParameterModel{}),
		frame.NewNamedAPI("", "POST", "/getUserInfo", &handler.UserInfoParameterModel{}),
		frame.NewNamedAPI("", "POST", "/alterUserInfo", &handler.AlterParameterModel{}),
		frame.NewNamedAPI("", "POST", "/feedback", &handler.FeedbackParameterModel{}),
		frame.NewNamedAPI("", "POST", "/sendVerificationCode", &handler.SendVerCodeParameterModel{}),
		frame.NewNamedAPI("", "POST", "/getUserNotify", &handler.NotifyParameterModel{}),
		frame.NewNamedAPI("", "POST", "/readUserNotify", &handler.ReadNotifyParameterModel{}),
		frame.NewNamedAPI("", "POST", "/getAddressManagentList", &handler.AddressManagentParameterModel{}),
		frame.NewNamedAPI("", "POST", "/alterAddressManagentInfo", &handler.AlterAddressInfoParameterModel{}),
		frame.NewNamedAPI("", "POST", "/deleteAddressManagent", &handler.DeleteAddressParameterModel{}),
		frame.NewNamedAPI("", "POST", "/getRedPacketRecording", &handler.RedPacketRecordingParameterModel{}),
		frame.NewNamedAPI("", "POST", "/getRedPacketInfo", &handler.RedPacketInfoParameterModel{}),
		frame.NewNamedAPI("", "POST", "/sendRedPacket", &handler.SendRedPacketParameterModel{}),
		frame.NewNamedAPI("", "POST", "/receiveRedPacket", &handler.ReceiveRedPacketParameterModel{}),

		frame.NewNamedAPI("", "POST", "/searchBalance", &handler.SearchParameterModel{}),

		frame.NewNamedAPI("", "POST", "/getCurrencyInfo", &handler.CurrencyInfoParameterModel{}),
		frame.NewNamedAPI("", "POST", "/getTradeRecording", &handler.TradeRecordParameterModel{}),
		frame.NewNamedAPI("", "POST", "/trade", &handler.TradeParameterModel{}),
		frame.NewNamedAPI("", "POST", "/exportPrivateKey", &handler.ExPrivateKeyParameterModel{}),
	) /*.Use(middleware.Token)*/
}
