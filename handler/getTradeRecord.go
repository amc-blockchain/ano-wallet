package handler

import (
	"blockChainWallet/database"

	//	"fmt"
	"math/big"
	"strconv"

	"github.com/henrylee2cn/faygo"
)

type TradeRecordParameterModel struct {
	TradeRecordId   int    `param:"<desc:ID> <in:formData> <required> <name:tradeRecordId> <err:>"`
	Address         string `param:"<desc:> <in:formData> <required> <name:address> <err:>"`
	ContractAddress string `param:"<desc:> <in:formData> <required> <name:contractAddress> <err:>"`
	Sign            string `param:"<desc:> <in:formData> <required> <name:sign> <err:>"`
}

type tradeRecordModel struct {
	RecordId          int64  `json:"recordId"`
	SendAddress       string `json:"sendAddress"`
	ReceiveAddress    string `json:"receiveAddress"`
	TransferNumber    string `json:"transferNumber"`
	MinerCosts        string `json:"minerCosts"`
	TransferHash      string `json:"transferHash"`
	BlockHeight       int64  `json:"blockHeight"`
	TransferTime      int64  `json:"transferTime"`
	CurrencyName      string `json:"currencyName"`
	ContractPercision int    `json:"contractPercision"`
}

// Serve impletes Handler.
func (trpm *TradeRecordParameterModel) Serve(ctx *faygo.Context) error {

	tr := make([]database.TradeRecording, 0)
	err1 := database.Engine.Where("id > ? and contract_address = ? and (send_address = ? or receive_address = ?)", trpm.TradeRecordId, trpm.ContractAddress, trpm.Address, trpm.Address).Desc("id").Limit(15).Find(&tr)
	if err1 == nil {
		trModel := make([]tradeRecordModel, 15)
		i := -1
		for index, value := range tr {
			trModel[index].RecordId = value.Id
			trModel[index].SendAddress = value.SendAddress
			trModel[index].ReceiveAddress = value.ReceiveAddress
			trModel[index].TransferNumber = big.NewFloat(value.TransferNumber).String()
			trModel[index].MinerCosts = strconv.FormatFloat(value.MinerCosts, 'f', 5, 64)
			trModel[index].TransferHash = value.TransferHash
			trModel[index].BlockHeight = value.BlockHeight
			trModel[index].TransferTime = value.TransferTime
			trModel[index].CurrencyName = value.CurrencyName
			trModel[index].ContractPercision = value.ContractPrecision
			i = index
		}
		return ctx.JSON(200, HandlerSucceed(trModel[:i+1]), true)
	}
	return ctx.JSON(200, HandlerCode(0), true)
}

// Doc returns the API's note, result or parameters information.
func (trpm *TradeRecordParameterModel) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "",
		Return: "// JSON\n{}",
	}
}
