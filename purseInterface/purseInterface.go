// purseInterface project purseInterface.go
package purseInterface

//import (
//	"fmt"
//	"time"
//	"strconv"
//)

type PurseHandlerInterface interface {
	initNet()
	getBlockNumber() (blockNum int64, err error)
	createChainAddress(tradePassword string) (addr string, err error)
	getGasPrice() (gasPrice string, err error)
	getBalance(addr, contractAddr, place string, precision int) (balance float64, err error)
	getBlockByNumber(number int64) (err error)
	sendTransaction(fromAddr, contractAddr, toAddr, tradePassword string, gasTotal, value float64, precision int) (tradeHash string, err error)
	sendRawTransaction(fromAddr, contractAddr, toAddr, tradePassword string, gasTotal, value float64, precision int) (tradeHash string, err error)
	getContractInfo(contractAddr string) (err error)
	exportWalletPrivateKey(addr, addrPasswd string) (privateKey string, err error)

	getJsonStr(addr string) (jsonStr string, err error)
}

func InitNet(phi PurseHandlerInterface) {
	phi.initNet()
}

/*
 *
 * parameter - phi
 * parameter -
 * return - blockNum
 * return - err
 */
func GetBlockNumber(phi PurseHandlerInterface) (blockNum int64, err error) {
	return phi.getBlockNumber()
}

/*
 *
 * parameter - phi
 * parameter - tradePassword
 * return - addr
 * return - err
 */
func CreateChainAddress(phi PurseHandlerInterface, tradePassword string) (addr string, err error) {
	return phi.createChainAddress(tradePassword)
}

/*
 *
 * parameter - phi
 * parameter -
 * return - gasPrice
 * return - err
 */
func GetGasPrice(phi PurseHandlerInterface) (gasPrice string, err error) {
	return phi.getGasPrice()
}

/*
 *
 * parameter - phi
 * parameter - addr
 * parameter - place
 * parameter - precision
 * return - balance
 * return - err
 */
func GetBalance(phi PurseHandlerInterface, addr, contractAddr, place string, precision int) (balance float64, err error) {
	return phi.getBalance(addr, contractAddr, place, precision)
}

/*
 *
 * parameter - phi
 * parameter - number
 * return - err
 */
func GetBlockByNumber(phi PurseHandlerInterface, number int64) (err error) {
	return phi.getBlockByNumber(number)
}

/*
 *
 * parameter - phi
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
func SendTransaction(phi PurseHandlerInterface, fromAddr, contractAddr, toAddr, tradePassword string, gasTotal, value float64, precision int) (tradeHash string, err error) {
	return phi.sendTransaction(fromAddr, contractAddr, toAddr, tradePassword, gasTotal, value, precision)
}
func SendRawTransaction(phi PurseHandlerInterface, fromAddr, contractAddr, toAddr, tradePassword string, gasTotal, value float64, precision int) (tradeHash string, err error) {
	return phi.sendRawTransaction(fromAddr, contractAddr, toAddr, tradePassword, gasTotal, value, precision)
}

/*
 *
 * parameter - phi
 * parameter - contractAddr
 * return - err
 */
func GetContractInfo(phi PurseHandlerInterface, contractAddr string) (err error) {
	return phi.getContractInfo(contractAddr)
}

/*
 *
 * parameter - addr
 * parameter - addrPasswd
 * return - privateKey
 * return - err
 */
func ExportWalletPrivateKey(phi PurseHandlerInterface, addr, addrPasswd string) (privateKey string, err error) {
	return phi.exportWalletPrivateKey(addr, addrPasswd)
}

func GetJsonStr(phi PurseHandlerInterface, addr string) (jsonStr string, err error) {
	return phi.getJsonStr(addr)
}

var (
	Mc MoacChain
)

func init() {
	InitNet(&Mc)

}
