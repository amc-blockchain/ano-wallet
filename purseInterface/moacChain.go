package purseInterface

import (
	"blockChainWallet/database"

	"blockChainWallet/jpushclient"

	"github.com/go-xorm/xorm"

	//	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mikemintang/go-curl"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethMath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"

	"io/ioutil"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	Chain3 "github.com/caivega/chain3go/chain3"
	Chain3common "github.com/caivega/chain3go/common"
	Chain3provider "github.com/caivega/chain3go/provider"
	Chain3rpc "github.com/caivega/chain3go/rpc"
)

var gourpS sync.Mutex
var addFlag bool = true

const (
	mocaInterfaceIp = ""
)

var moacNetRequset *curl.Request

type MocaContractInfo struct {
	ContractAddr      string
	ContractName      string
	ContractSymbol    string
	Contractprecision int64
	ContractTotal     int64
}

var (
	accountsManager *accounts.Manager = nil
	moacClient      *ethclient.Client = nil
	rpcClient       *rpc.Client       = nil
	chain3Client    *Chain3.Chain3    = nil
)

type MoacChain struct {
	MocaContractInfo
}

func (mc *MoacChain) makeMoacClient() (*ethclient.Client, *rpc.Client, error) {
	client, err := rpc.Dial(mocaInterfaceIp)
	if err != nil {
		return nil, nil, err
	}
	conn := ethclient.NewClient(client)

	return conn, client, nil
}

func (mc *MoacChain) initNet() {

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	moacNetRequset = curl.NewRequest().SetUrl(mocaInterfaceIp).SetHeaders(headers)

	var err error
	accountsManager, err = mc.makeAccountsManager()
	if err != nil {
		panic(err)
	}

	moacClient, rpcClient, err = mc.makeMoacClient()
	if err != nil {
		panic(err)
	}

	chain3Client = Chain3.NewChain3(Chain3provider.NewHTTPProvider(mocaInterfaceIp, Chain3rpc.GetDefaultMethod()))

}

/*
 *
 * parameter - contractAddr
 * return - err
 */
func (mc *MoacChain) addContractAddress(contractAddr string) (err error) {

	err = mc.getContractInfo(contractAddr)
	if err == nil {
		var cm database.CurrencyManagement
		flag2, err2 := database.Engine.Where("currency_contract_address = ?", contractAddr).Get(&cm)
		if err2 == nil {
			if flag2 {
				var countId int64 = 1
				var counter int = 0
				for {
					flag, err1 := database.Engine.Where("id = ?", countId).Exist(&database.UserInfo{})
					if err1 == nil {
						if flag {
							counter = 0

							var ua database.UserAssets
							ua.UserId = countId
							ua.CurrencyNumber = 0
							ua.CurrencyId = cm.CurrencyId
							err = database.SessionSubmit(func(session *xorm.Session) (err1 error) {
								_, err1 = session.Insert(ua)
								return err1
							})
							if err != nil {
								continue
							}
						} else {
							counter++
						}
						if counter == 11 {
							break
						}
						countId++
					}
				}
			}
		}
	}

	return err
}

/*
 *
 */
func (mc *MoacChain) timerBlock() {

	timer := time.NewTimer(time.Second * 5)
	for {
		select {

		case <-timer.C:
			bhr := new(database.BlockHeightRecording)
			flag, err := database.Engine.Where("chain_type = ?", "MOAC").Get(bhr)
			if err == nil && flag && addFlag {
				addFlag = false
				go mc.getBlockByNumber(bhr.BlockHeight)
			} else {
				fmt.Println(flag, err)
			}
			timer.Reset(time.Second * 5)
		}
	}
}

/*
 *
 */
func (mc *MoacChain) makeAccountsManager() (*accounts.Manager, error) {

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP

	var (
		keydir string = "./keystore/moac"
	)

	if err := os.MkdirAll(keydir, 0700); err != nil {
		return nil, err
	}
	// Assemble the accounts manager and supported backends
	backends := []accounts.Backend{
		keystore.NewKeyStore(keydir, scryptN, scryptP),
	}

	return accounts.NewManager(backends...), nil
}

/*
 *
 */
func (mc *MoacChain) fetchKeystore() *keystore.KeyStore {
	if accountsManager == nil {
		accountsManager, _ = mc.makeAccountsManager()
	}
	return accountsManager.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
}

/*
 *
 * parameter -
 * return - blockNum
 * return - err
 */
func (mc *MoacChain) getBlockNumber() (blockNum int64, err error) {

	blockNum = -1
	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
	}()

	var blockNumStr string
	err = rpcClient.Call(&blockNumStr, "mc_blockNumber")
	if err == nil {
		blockNum, err = strconv.ParseInt(blockNumStr[2:], 16, 64)
	}

	return blockNum, err
}

/*
 *
 * parameter - tradePassword
 * return - addr
 * return - err
 */
func (mc *MoacChain) createChainAddress(tradePassword string) (addr string, err error) {

	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
	}()

	err = rpcClient.Call(&addr, "personal_newAccount", tradePassword)

	return addr, err
}

/*
 *
 * parameter -
 * return - gasPrice
 * return - err
 */
func (mc *MoacChain) getGasPrice() (gasPrice string, err error) {

	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
	}()

	err = rpcClient.Call(&gasPrice, "mc_gasPrice")

	return gasPrice, err
}

/*
 *
 * parameter - addr
 * parameter - precision
 * return - balance
 * return - err
 */
func (mc *MoacChain) getBalance(addr, contractAddr, place string, precision int) (balance float64, err error) {

	if place == "" {
		place = "latest"
	}
	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
	}()

	var tmpBalance string
	if contractAddr == "" {
		err = rpcClient.Call(&tmpBalance, "mc_getBalance", addr, place)
	} else {
		err = rpcClient.Call(&tmpBalance, "mc_call", map[string]interface{}{"to": contractAddr, "data": ("0x70a08231000000000000000000000000" + addr[2:len(addr)])}, place)
	}
	if err == nil {
		tmpBalanceBig, _, tErr := big.ParseFloat(tmpBalance[2:], 16, 'f', big.ToZero)
		if tErr == nil {
			//balance = float64(tmpBalance.Int64()) * math.Pow10(-precision)
			balance, _ = new(big.Float).Mul(big.NewFloat(math.Pow10(-precision)), tmpBalanceBig).Float64()
		} else {
			err = errors.New("")
		}
	}

	return balance, err
}

/*
 *
 * parameter - contractAddr
 * return - err
 */
func (mc *MoacChain) getContractInfo(contractAddr string) (err error) {

	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
	}()

	var tmpCM database.CurrencyManagement
	flag2, err2 := database.Engine.Where("currency_contract_address = ?", contractAddr).Get(&tmpCM)
	if err2 == nil {
		if !flag2 {
			contractFunction := func(postData map[string]interface{}) (valueStr string, err error) {

				resp, netErr := moacNetRequset.SetPostData(postData).Post()
				if netErr == nil {
					var resultMap map[string]interface{}
					json.Unmarshal([]byte(resp.Body), &resultMap)
					if value, flag := resultMap["result"]; flag {
						valueStr = value.(string)
					} else {
						err = errors.New((((resultMap["error"]).(map[string]interface{}))["message"]).(string))
					}
				} else {
					err = netErr
				}
				return valueStr, err
			}
			strHandler := func(str string) string {

				var index int = 64
				for ; index < len(str); index++ {
					if str[index] == 0 {
						break
					}
				}
				return str[64:index]
			}

			mc.ContractAddr = contractAddr
			var tmpStr string

			tmpStr, err = contractFunction(map[string]interface{}{
				"id":      "101",
				"jsonrpc": "2.0",
				"method":  "mc_call",
				"params":  [2]interface{}{map[string]string{"to": contractAddr, "data": "0x06fdde03"}, "latest"},
			})
			data, _ := hex.DecodeString(tmpStr[2:len(tmpStr)])
			mc.ContractName = strHandler(string(data))
			if err != nil {
				return err
			}

			tmpStr, err = contractFunction(map[string]interface{}{
				"id":      "101",
				"jsonrpc": "2.0",
				"method":  "mc_call",
				"params":  [2]interface{}{map[string]string{"to": contractAddr, "data": "0x95d89b41"}, "latest"},
			})
			data, _ = hex.DecodeString(tmpStr[2:len(tmpStr)])
			mc.ContractSymbol = strHandler(string(data))
			if err != nil {
				return err
			}

			tmpStr, err = contractFunction(map[string]interface{}{
				"id":      "101",
				"jsonrpc": "2.0",
				"method":  "mc_call",
				"params":  [2]interface{}{map[string]string{"to": contractAddr, "data": "0x313ce567"}, "latest"},
			})
			mc.Contractprecision, _ = strconv.ParseInt(tmpStr[2:len(tmpStr)], 16, 64)
			if err != nil {
				return err
			}

			tmpStr, err = contractFunction(map[string]interface{}{
				"id":      "101",
				"jsonrpc": "2.0",
				"method":  "mc_call",
				"params":  [2]interface{}{map[string]string{"to": contractAddr, "data": "0x18160ddd"}, "latest"},
			})
			mc.ContractTotal, _ = strconv.ParseInt(tmpStr[2:len(tmpStr)], 16, 64)
			mc.ContractTotal = mc.ContractTotal / int64(math.Pow10(int(mc.Contractprecision)))
			if err != nil {
				return err
			}

			cm := new(database.CurrencyManagement)
			cm.ChainType = "MOAC"
			cm.ContractPrecision = int(mc.Contractprecision)
			cm.CurrencyContractAddress = contractAddr
			cm.CurrencyImageUrl = ""
			cm.CurrencyName = mc.ContractSymbol
			cm.CurrencyNameAbbreviation = mc.ContractName
			cm.CurrencyTotal = mc.ContractTotal
			cm.RedPacketAddr = ""
			err = database.SessionSubmit(func(session *xorm.Session) error {

				_, err1 := session.Insert(cm)
				return err1
			})
		} else {
			err = errors.New("")
		}
	} else {
		err = err2
	}

	return err
}

/*
 *
 * parameter - number
 * return - err
 */
func (mc *MoacChain) getBlockByNumber(number int64) (err error) {

	defer func() {
		addFlag = true
		if re := recover(); re != nil {
			err = re.(error)
			fmt.Println(err)
		}
	}()

	postData := map[string]interface{}{
		"id":      "101",
		"jsonrpc": "2.0",
		"method":  "mc_getBlockByNumber",
		"params":  [2]interface{}{"0x" + strconv.FormatInt(number, 16), true},
	}

	resp, netErr := moacNetRequset.SetPostData(postData).Post()
	fmt.Println(resp, netErr)
	if netErr == nil {
		var resultMap map[string]interface{}
		json.Unmarshal([]byte(resp.Body), &resultMap)
		if value, flag := resultMap["result"]; flag {
			timestamp, _ := strconv.ParseInt((value.(map[string]interface{}))["timestamp"].(string)[2:], 16, 64)
			ts := ((value.(map[string]interface{}))["transactions"]).([]interface{})
			for _, mapValue := range ts {
				err = mc.accountProcessing(mapValue.(map[string]interface{}), timestamp, number)
				if err != nil {
					return err
				}
			}

			var flag bool
			bhr := new(database.BlockHeightRecording)
			flag, err = database.Engine.Where("chain_type = ?", "MOAC").Get(bhr)
			if err != nil {
				return
			}
			err = database.SessionSubmit(func(session *xorm.Session) (err1 error) {
				if flag {
					bhr.BlockHeight = number + 1
					_, err1 = session.Table("bcw_block_height_recording").Where("chain_type = ?", "MOAC").Update(map[string]interface{}{"blockHeight": bhr.BlockHeight})
				} else {
					bhr.BlockHeight = number
					bhr.ChainType = "MOAC"
					_, err1 = session.Insert(bhr)
				}
				return err1
			})
			if err == nil {
				fmt.Println("：", number)
			}
		} else {
			err = errors.New((((resultMap["error"]).(map[string]interface{}))["message"]).(string))
		}
	} else {
		err = netErr
	}

	return err
}

/*
 *
 * parameter - trade
 * parameter - timestamp
 * parameter - number
 * return -
 */
func (mc *MoacChain) accountProcessing(trade map[string]interface{}, timestamp int64, number int64) (err error) {

	gourpS.Lock()
	defer gourpS.Unlock()

	type tradeModel struct {
		BlockNumber      string `json:"blockNumber"`
		Gas              string `json:"gas"`
		GasPrice         string `json:"gasPrice"`
		From             string `json:"from"`
		To               string `json:"to"`
		Value            string `json:"value"`
		TransactionIndex string `json:"transactionIndex"`
		ShardingFlag     string `json:"shardingFlag"`
		Syscnt           string `json:"syscnt"`
		Nonce            string `json:"nonce"`
		V                string `json:"v"`
		R                string `json:"r"`
		S                string `json:"s"`
		Input            string `json:"input"`
		Hash             string `json:"hash"`
	}

	fmt.Println("trade", trade)

	var tm tradeModel
	bytes, err := json.Marshal(trade)
	if err == nil {
		json.Unmarshal(bytes, &tm)

		if err5 := mc.getTransactionReceipt(tm.Hash); err5 != nil && err5.Error() == "" {
			return err
		}

		tm.From = strings.ToLower(tm.From)
		tm.To = strings.ToLower(tm.To)
		tm.Hash = strings.ToLower(tm.Hash)
		if tm.From == "" || tm.To == "" {
			return
		}
		if tm.TransactionIndex != "0x0" {
			uAddrs := make([]database.UserAddress, 0)

			if len(tm.Input) > 10 && tm.Input[2:10] == "a9059cbb" && len(tm.Input) == 138 {
				err = database.Engine.Where("currency_address = ? or currency_address = ?", tm.From, "0x"+tm.Input[34:74]).Find(&uAddrs)
			} else {
				err = database.Engine.Where("currency_address = ? or currency_address = ?", tm.From, tm.To).Find(&uAddrs)
			}
			if err == nil {
				for _, value := range uAddrs {
					tr := new(database.TradeRecording)
					flag, tmpErr := database.Engine.Where("transfer_hash = ?", tm.Hash).Get(tr)

					var currencyId uint64
					var vn *big.Float
					if tmpErr == nil {
						if !flag {

							tr.ChainType = "MOAC"
							if len(tm.Input) > 10 && tm.Input[2:10] == "a9059cbb" && len(tm.Input) == 138 {

								tr.ContractAddress = tm.To
								tr.ReceiveAddress = "0x" + tm.Input[34:74]
								//								vn, _ = strconv.ParseUint(tm.Input[74:], 16, 64)
								vn, _, _ = big.ParseFloat(tm.Input[74:], 16, 'f', big.ToZero)
							} else {

								tr.ContractAddress = ""
								tr.ReceiveAddress = tm.To
								//								vn, _ = strconv.ParseUint(tm.Value[2:], 16, 64)

								vn, _, _ = big.ParseFloat(tm.Value[2:], 16, 'f', big.ToZero)
							}
							tr.SendAddress = tm.From
							gas, _, _ := big.ParseFloat(tm.Gas[2:], 16, 'f', big.ToZero)
							gasPrice, _, _ := big.ParseFloat(tm.GasPrice[2:], 16, 'f', big.ToZero)
							gasTotal, _ := gas.Mul(gas, gasPrice).Float64()
							tr.MinerCosts = gasTotal / math.Pow10(18)
							tr.TransferHash = tm.Hash
							tr.TransferData = tm.Input
							tr.BlockHeight, _ = strconv.ParseInt(tm.BlockNumber[2:], 16, 64)
							tr.TransferTime = timestamp
						}
						cm := new(database.CurrencyManagement)
						database.Engine.Where("currency_contract_address = ? and chain_type = ?", tr.ContractAddress, "MOAC").Get(cm)
						if cm.CurrencyName == "" {
							//
							//							err = mc.getContractInfo(tm.To)
							//							if err != nil {
							//								return
							//							} else {
							//								database.Engine.Where("currency_contract_address = ? and chain_type = ?", tr.ContractAddress, "MOAC").Get(cm)
							//								if cm.CurrencyName == "" {
							//									return
							//								}
							//							}
							return //
						}
						currencyId = uint64(cm.CurrencyId)
						tr.ChainType = cm.ChainType
						tr.CurrencyName = cm.CurrencyName
						tr.ContractPrecision = cm.ContractPrecision
						if vn != nil {
							if fv, _ := vn.Float64(); fv != 0 {
								//tr.TransferNumber = vn //float64(vn) / math.Pow10(cm.ContractPrecision)

								tr.TransferNumber, _ = new(big.Float).Mul(big.NewFloat(math.Pow10(-cm.ContractPrecision)), vn).Float64()
							}
						}

						tmpTr := new(database.TradeRecording)
						dirNumber := tr.TransferNumber
						if value.CurrencyAddress == tr.SendAddress {

							tr.SendUserId = value.UserId
							dirNumber = -tr.TransferNumber
							//							dirNumber = 0
							existFlag, existErr := database.Engine.Where("send_user_id = ? and transfer_hash = ?", value.UserId, tr.TransferHash).Get(tmpTr)
							if existErr == nil && existFlag {
								continue
							}

							moacCurrencyId := int64(2)
							var tmpCM database.CurrencyManagement
							database.Engine.Where("chain_type = ? and currency_name = ?", "MOAC", "MOAC").Get(&tmpCM)
							if tmpCM.CurrencyId != 0 {
								moacCurrencyId = tmpCM.CurrencyId
							}

							userA := new(database.UserAssets)
							eFlag, eErr := database.Engine.Where("user_id = ? and currency_id = ?", value.UserId, moacCurrencyId).Get(userA)
							if eErr == nil {
								moacnumber, _ := mc.getBalance(value.CurrencyAddress, "", "0x"+strconv.FormatInt(number, 16), 18)
								database.SessionSubmit(func(session *xorm.Session) (err1 error) {
									userA.CurrencyNumber = moacnumber
									if eFlag {
										_, err1 = session.Where("user_id = ? and currency_id = ?", value.UserId, moacCurrencyId).Update(userA)
									} else {
										userA.CurrencyId = int64(moacCurrencyId)
										userA.UserId = value.UserId
										_, err1 = session.Insert(tr)
									}
									return err1
								})
							}
						} else {

							tr.AcceptUserId = value.UserId
							existFlag, existErr := database.Engine.Where("accept_user_id = ? and transfer_hash = ?", value.UserId, tr.TransferHash).Get(tmpTr)
							if existErr == nil && existFlag {
								continue
							}
						}

						var flag1 bool
						ua := new(database.UserAssets)
						flag1, err = database.Engine.Where("user_id = ? and currency_id = ?", value.UserId, currencyId).Get(ua)
						if err != nil {
							return
						}
						ua.CurrencyNumber = ua.CurrencyNumber + dirNumber
						err = database.SessionSubmit(func(session *xorm.Session) (err1 error) {
							if flag {
								_, err1 = session.Where("transfer_hash = ?", tm.Hash).Update(tr)
							} else {
								_, err1 = session.Insert(tr)
							}
							if err1 == nil {
								if flag1 {
									_, err1 = session.Table("bcw_user_assets").Where("user_id = ? and currency_id = ?", value.UserId, currencyId).Update(map[string]interface{}{"currency_number": ua.CurrencyNumber})
								} else {
									ua.CurrencyId = int64(currencyId)
									ua.UserId = value.UserId
									ua.CurrencyNumber = dirNumber
									_, err1 = session.Insert(ua)
								}
							}
							return err1
						})

						if err == nil {
							uInfos := make([]database.UserInfo, 0)
							var err2 error
							if value.CurrencyAddress == tr.SendAddress {

								err2 = database.Engine.Where("id = ?", tr.SendUserId).Find(&uInfos)
							} else {

								err2 = database.Engine.Where("id = ?", tr.AcceptUserId).Find(&uInfos)
							}
							if err2 == nil {
								for _, tmpInfo := range uInfos {
									var uNotify database.UserNotify
									var tmpStr string
									if value.CurrencyAddress == tr.SendAddress {

										tmpStr = "" + big.NewFloat(tr.TransferNumber).String() + "" + tr.CurrencyName + ""
									} else {

										tmpStr = "" + big.NewFloat(tr.TransferNumber).String() + "" + tr.CurrencyName + ""
									}
									uNotify.NotifyContent = tmpStr
									uNotify.NotifyReadFlag = 0
									uNotify.NotifyTitle = ""
									uNotify.NotifyUrl = ""
									uNotify.UserId = tmpInfo.Id
									database.SessionSubmit(func(session *xorm.Session) (err1 error) {
										_, err1 = session.Insert(uNotify)
										return err1
									})
									if tmpInfo.PushId != "" {
										jpushclient.SendJPush(tmpStr, tmpInfo.PushId, tmpInfo.MachineType)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return err
}

/*
 *
 * parameter - hashStr
 * return - err
 */
func (mc *MoacChain) getTransactionReceipt(hashStr string) (err error) {

	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
	}()

	postData := map[string]interface{}{
		"id":      "101",
		"jsonrpc": "2.0",
		"method":  "mc_getTransactionReceipt",
		"params":  [1]interface{}{hashStr},
	}

	resp, netErr := moacNetRequset.SetPostData(postData).Post()
	if netErr == nil {
		var resultMap map[string]interface{}
		json.Unmarshal([]byte(resp.Body), &resultMap)
		if value, flag := resultMap["result"]; flag {
			if ((value).(map[string]interface{}))["status"] == "0x0" {

				err = errors.New("")
			}
		} else {
			err = errors.New((((resultMap["error"]).(map[string]interface{}))["message"]).(string))
		}
	} else {
		err = netErr
	}

	return err
}

/*
 *
 * parameter - addr
 * parameter - tradePassword
 * return - err
 */
func (mc *MoacChain) unlockAddr(addr, tradePassword string) (err error) {

	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
	}()

	postData := map[string]interface{}{
		"id":      "101",
		"jsonrpc": "2.0",
		"method":  "personal_unlockAccount",
		"params":  [2]string{addr, tradePassword},
	}

	resp, netErr := moacNetRequset.SetPostData(postData).Post()
	if netErr == nil {
		var resultMap map[string]interface{}
		json.Unmarshal([]byte(resp.Body), &resultMap)
		if _, flag := resultMap["result"]; flag {
			fmt.Println(addr, "")
		} else {
			err = errors.New((((resultMap["error"]).(map[string]interface{}))["message"]).(string))
		}
	} else {
		err = netErr
	}
	return err
}

/*
 *
 */
func (mc *MoacChain) lockAddr(addr string) (err error) {

	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
	}()

	postData := map[string]interface{}{
		"id":      "101",
		"jsonrpc": "2.0",
		"method":  "personal_lockAccount",
		"params":  [1]string{addr},
	}

	resp, netErr := moacNetRequset.SetPostData(postData).Post()
	if netErr == nil {
		var resultMap map[string]interface{}
		json.Unmarshal([]byte(resp.Body), &resultMap)
		if _, flag := resultMap["result"]; flag {
			fmt.Println(addr, "")
		} else {
			err = errors.New((((resultMap["error"]).(map[string]interface{}))["message"]).(string))
		}
	} else {
		err = netErr
	}
	return err
}

/*
 *
 * parameter - fromAddr
 * parameter - contractAddr
 * parameter - toAddr
 * parameter - tradePassword
 * parameter - gasTotal
 * parameter - value
 * parameter - precision
 * return - tradeHash
 * return - err
 */
func (mc *MoacChain) sendTransaction(fromAddr, contractAddr, toAddr, tradePassword string, gasTotal, value float64, precision int) (tradeHash string, err error) {

	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
	}()

	if precision > 4 {
		return "", errors.New("")
	}

	hexStr := hex.EncodeToString(new(big.Int).Mul(big.NewInt(int64(value*10000)), big.NewInt(int64(math.Pow10(precision-4)))).Bytes())

	var dict map[string]interface{}
	if contractAddr == "" {

		dict = map[string]interface{}{
			"from":     fromAddr,
			"to":       toAddr,
			"gas":      "0x" + strconv.FormatInt(int64(gasTotal*50000000), 16),
			"gasPrice": "0x4a817c800", //20000000000
			"value":    "0x" + hexStr,
			"data":     "0x",
		}
	} else {
		placeholderStr := "0000000000000000000000000000000000000000000000000000000000000000"
		numberStr := hexStr

		dict = map[string]interface{}{
			"from":     fromAddr,
			"to":       contractAddr,
			"gas":      "0x" + strconv.FormatInt(int64(gasTotal*50000000), 16),
			"gasPrice": "0x4a817c800", //20000000000
			"value":    "0x0",
			"data":     "0xa9059cbb" + placeholderStr[:(64-len(toAddr[2:len(toAddr)]))] + toAddr[2:len(toAddr)] + placeholderStr[:(64-len(numberStr))] + numberStr,
		}
	}

	err = rpcClient.Call(&tradeHash, "personal_sendTransaction", dict, tradePassword)

	return tradeHash, err
}

func (mc *MoacChain) sendRawTransaction(fromAddr, contractAddr, toAddr, tradePassword string, gasTotal, value float64, precision int) (tradeHash string, err error) {

	var prStr string
	prStr, err = mc.exportWalletPrivateKey(fromAddr, tradePassword)
	fmt.Println(prStr)
	if err != nil {
		return
	}

	txcount, _ := chain3Client.Mc.GetTransactionCount(Chain3common.StringToAddress(fromAddr), "pending")

	gasPrice := 25000000000
	gasLimit := 100000

	chainId, _ := chain3Client.Net.Version()
	fmt.Println(chainId)

	var txData *types.Transaction
	var v, r, s *big.Int
	var acc accounts.Account
	acc, err = mc.fetchKeystore().Find(accounts.Account{Address: common.HexToAddress(fromAddr)})
	if err == nil {
		var jsonBytes []byte
		jsonBytes, err = ioutil.ReadFile(acc.URL.Path)
		if err == nil {
			var storeKey *keystore.Key
			storeKey, err = keystore.DecryptKey(jsonBytes, tradePassword)
			if err == nil {
				txData = types.NewTransaction(txcount.Uint64(), common.HexToAddress(toAddr[:]), new(big.Int).Mul(big.NewInt(int64(0.01*1e8)), big.NewInt(1e10)), uint64(gasLimit), big.NewInt(int64(gasPrice)), nil)

				txData, _ = types.SignTx(txData, types.HomesteadSigner{}, storeKey.PrivateKey)

				//				jsonStr, _ := json.Marshal(txData)
				//				fmt.Println(string(jsonStr))

				v, r, s = txData.RawSignatureValues()

			}
		}
	}

	type txdata struct {
		AccountNonce   uint64          `json:"nonce"    gencodec:"required"`
		SystemContract uint64          `json:"syscnt" gencodec:"required"`
		Price          *big.Int        `json:"gasPrice" gencodec:"required"`
		GasLimit       *big.Int        `json:"gas"      gencodec:"required"`
		Recipient      *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
		Amount         *big.Int        `json:"value"    gencodec:"required"`
		Payload        []byte          `json:"input"    gencodec:"required"`
		ShardingFlag   uint64          `json:"shardingFlag" gencodec:"required"`
		Via            *common.Address `json:"via"       rlp:"nil"`

		// Signature values
		V *big.Int `json:"v" gencodec:"required"`
		R *big.Int `json:"r" gencodec:"required"`
		S *big.Int `json:"s" gencodec:"required"`

		// This is only used when marshaling to JSON.
		Hash *common.Hash `json:"hash" rlp:"-"`
	}

	var rawTx = map[string]interface{}{
		"hash":         txData.Hash().String(),
		"nonce":        "0xa",
		"to":           txData.To().String(),
		"gas":          "0x186a0",
		"gasPrice":     "0x5d21dba00",
		"value":        "0x2386f26fc10000",
		"input":        "0x",
		"chainId":      "0x63",
		"shardingFlag": "0x0",
		"syscnt":       "0x0",
		"via":          "0x",
		"v":            v.String(),
		"r":            r.String(),
		"s":            s.String(),
	}

	txByet, _ := json.Marshal(rawTx)

	fmt.Println(string(txByet))

	//	//	chain3Client.Sha3(string(txByet), prStr)

	txHash, err := chain3Client.Mc.SendRawTransaction(txByet)
	if err == nil {
		fmt.Println(txHash)
	} else {
		fmt.Println(rawTx)
		fmt.Println("err:", err)
	}

	return
}

/*
 *
 * privateKey -
 * passwd -
 * return - err
 */
func (mc *MoacChain) ImportWalletPrivateKey(privateKey, passwd string) (err error) {

	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
	}()

	return errors.New("")
}

/*
 *
 * parameter - addr
 * parameter - addrPasswd
 * return - privateKey
 * return - err
 */
func (mc *MoacChain) exportWalletPrivateKey(addr, addrPasswd string) (privateKey string, err error) {

	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
	}()

	var acc accounts.Account
	acc, err = mc.fetchKeystore().Find(accounts.Account{Address: common.HexToAddress(addr)})
	if err == nil {
		var jsonBytes []byte
		jsonBytes, err = ioutil.ReadFile(acc.URL.Path)
		if err == nil {
			var storeKey *keystore.Key
			storeKey, err = keystore.DecryptKey(jsonBytes, addrPasswd)
			if err == nil {
				privateKey = hex.EncodeToString(ethMath.PaddedBigBytes(storeKey.PrivateKey.D, storeKey.PrivateKey.Params().BitSize/8))
			}
		}
	}

	return privateKey, err
}

/*
 *
 * parameter - addr
 * return - count
 * return - err
 */
func (mc *MoacChain) getTransactionCount(addr string) (count string, err error) {

	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
	}()

	postData := map[string]interface{}{
		"id":      "101",
		"jsonrpc": "2.0",
		"method":  "mc_getTransactionCount",
		"params":  [2]string{addr, "pending"},
	}

	resp, netErr := moacNetRequset.SetPostData(postData).Post()
	if netErr == nil {
		var resultMap map[string]interface{}
		json.Unmarshal([]byte(resp.Body), &resultMap)
		if value, ok := resultMap["result"]; ok {
			count = value.(string)
		} else {
			err = errors.New((((resultMap["error"]).(map[string]interface{}))["message"]).(string))
		}
	} else {
		err = netErr
	}

	return count, err
}

func (mc *MoacChain) getJsonStr(addr string) (jsonStr string, err error) {

	//	fmt.Println(addr)
	var acc accounts.Account
	acc, err = mc.fetchKeystore().Find(accounts.Account{Address: common.HexToAddress(addr)})
	if err == nil {
		var jsonBytes []byte
		jsonBytes, err = ioutil.ReadFile(acc.URL.Path)
		if err == nil {
			jsonStr = string(jsonBytes)

		}
	}
	return jsonStr, err
}
