// database project database.go
package database

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

	"errors"
)

type UserInfo struct {
	Id int64 `xorm:"bigint(20) pk autoincr notnull 'id'"`

	IdentityId string `xorm:"varchar(18) 'identity_id'"`
	Name       string `xorm:"varchar(10) 'name'"`
	Sex        int    `xorm:"int(1) 'sex' default(1)"`
	Age        int    `xorm:"int(3) 'age' default(20)"`

	MachineType   string `xorm:"varchar(10) 'machine_type'"`
	PushId        string `xorm:"varchar(64) 'push_id'"`
	EffectiveFlag int    `xorm:"int(1) 'effective_flag' default(1)"`

	Nickname      string `xorm:"varchar(20) 'nickname' default('nickname')"`
	PhoneNumber   string `xorm:"varchar(11) 'phone_number'"`
	Password      string `xorm:"varchar(60) 'password'"`
	TradePassword string `xorm:"varchar(60) 'trade_password'"`

	RedPacketTotal float64 `xorm:"decimal(20,6) 'red_packet_total' default(0)"`

	CreateTimestamp int64 `xorm:"bigint(20) created 'create_timestamp'"`
}

type UserAddress struct {
	UserId            int64  `xorm:"bigint(20) notnull 'user_id'"`
	AddressName       string `xorm:"varchar(25) 'address_name' default('wallet')"`
	CurrencyAddress   string `xorm:"varchar(60)" 'currency_address'`
	ChainType         string `xorm:"varchar(10) 'chain_type'"`
	AddressPrivateKey string `xorm:"varchar(255) 'address_private_key'"`
}

type CurrencyManagement struct {
	CurrencyId               int64  `xorm:"bigint(20) pk autoincr notnull 'id'"`
	ChainType                string `xorm:"varchar(10) 'chain_type'"`
	CurrencyName             string `xorm:"varchar(255) 'currency_name'"`
	CurrencyTotal            int64  `xorm:"bigint(20) 'currency_total'"`
	CurrencyNameAbbreviation string `xorm:"varchar(255) 'currency_name_abbreviation'"`
	CurrencyImageUrl         string `xorm:"varchar(255) 'currency_image_url'"`
	CurrencyContractAddress  string `xorm:"varchar(60)" 'currency_contract_address'`
	ContractPrecision        int    `xorm:"int(2) 'contract_precision'"`
	RedPacketAddr            string `xorm:varchar(60) 'red_packet_addr'`
}

type UserAssets struct {
	UserId         int64   `xorm:"bigint(20) notnull 'user_id'"`
	CurrencyId     int64   `xorm:"bigint(20) 'currency_id'"`
	CurrencyNumber float64 `xorm:"decimal(24,8) 'currency_number' default(0)"`
}

type TradeRecording struct {
	Id int64 `xorm:"bigint(20) pk notnull autoincr 'id'"`

	SendUserId   int64 `xorm:"bigint(20) 'send_user_id'"`
	AcceptUserId int64 `xorm:"bigint(20) 'accept_user_id'"`

	ContractAddress string  `xorm:"varchar(64) 'contract_address'"`
	SendAddress     string  `xorm:"varchar(64) 'send_address'"`
	ReceiveAddress  string  `xorm:"varchar(64) 'receive_address'"`
	TransferNumber  float64 `xorm:"decimal(24,8) 'transfer_number'"`
	MinerCosts      float64 `xorm:"decimal(15,8) 'miner_costs'"`
	TransferHash    string  `xorm:"varchar(70) 'transfer_hash'"`
	TransferData    string  `xorm:"text 'transfer_data'"`
	BlockHeight     int64   `xorm:"int(10) 'block_height'"`
	TransferTime    int64   `xorm:"bigint(20) 'transfer_time'"`

	ChainType         string `xorm:"varchar(10) 'chain_type'"`
	CurrencyName      string `xorm:"varchar(255) 'currency_name'"`
	ContractPrecision int    `xorm:"int(2) 'contract_precision'"`
}

type UserNotify struct {
	Id     int64 `xorm:"bigint(20) pk autoincr notnull 'id'"`
	UserId int64 `xorm:"bigint(20) notnull 'user_id'"`

	NotifyReadFlag int    `bson:"int(1) 'notify_read_flag'"`
	NotifyTitle    string `bson:"varchar(128) 'notify_title'"`
	NotifyContent  string `bson:"varchar(255) 'notify_content'"`
	NotifyUrl      string `bson:"varchar(255) 'notify_url'"`

	NotifyTimestamp int64 `xorm:"bigint(20) created 'notify_timestamp'"`
}

type UserFeedback struct {
	UserId                 int64  `xorm:"bigint(20) notnull 'user_id'"`
	OpinionContact         string `xorm:"varchar(30) 'opinion_contact'"`
	OpinionContent         string `xorm:"varchar(255) 'opinion_content'"`
	OpinionCreateTimestamp int64  `xorm:"bigint(20) created 'opinion_create_timestamp'"`
}

type UserAddressManagement struct {
	Id               int64  `xorm:"bigint(20) pk autoincr notnull 'id'"`
	UserId           int64  `xorm:"bigint(20) notnull 'user_id'"`
	AddressDesc      string `xorm:"varchar(255) 'address_desc'"`
	MatchAddress     string `xorm:"varchar(70) 'match_address'"`
	AddressAlterTime int64  `xorm:"bigint(20) updated 'address_alter_time'"`
}

type BlockHeightRecording struct {
	ChainType   string `xorm:"varchar(10) 'chain_type'"`
	BlockHeight int64  `xorm:"int(10) 'blockHeight'"`
	UpdateTime  int64  `xorm:"bitint(20) updated 'update_time'"`
}

type RedPacket struct {
	Id              int64   `xorm:"bigint(20) pk autoincr notnull 'id'"`
	UserId          int64   `xorm:"bigint(20) 'user_id'"`
	RedPacketType   int     `xorm:"int(2) 'red_packet_type'"`
	PayType         int     `xorm:"int(2) 'pay_type'"`
	RedPacketAmount float64 `xorm:"decimal(20,6) 'red_packet_amount'"`
	RemainingAmount float64 `xorm:"decimal(20,6) 'remaining_amount'"`
	RedPacketNumber int64   `xorm:"int(8) 'red_packet_number'"`
	RemainingNumber int64   `xorm:"int(8) 'remaining_number'"`
	RedPacketDesc   string  `xorm:"varchar(255) 'red_packet_desc'"`
	CreateTime      int64   `xorm:"bigint(20) 'create_time'"`
	FinishTime      int64   `xorm:"bigint(20) 'finish_time'"`
	ExpireTime      int64   `xorm:"bigint(20) 'expire_time'"`
	HasRefund       int     `xorm:"int(2) 'has_refund'"`
	CurrencyNumber  string  `xorm:"varchar(255) 'currency_number'"`
}

type RedPacketRecording struct {
	Id              int64   `xorm:"bigint(20) pk autoincr notnull 'id'"`
	RedPacketId     int64   `xorm:"bigint(20) 'red_packet_id'"`
	RedPacketAmount float64 `xorm:"decimal(20,6) 'red_packet_amount'"`
	PhoneNumber     string  `xorm:"varchar(11) 'phone_number'"`
	CreateTime      int64   `xorm:"bigint(20) 'create_time'"`
	ReceiveTime     int64   `xorm:"bigint(20) 'receive_time'"`
	Status          int     `xorm:"int(2) 'status'"`
	RedPacketType   int     `xorm:"int(2) 'red_packet_type'"`
	CurrencyNumber  string  `xorm:"varchar(255) 'currency_number'"`
}

var Engine *xorm.Engine

func init() {

	var err error
	Engine, err = xorm.NewEngine("", "")
	if err != nil {
		panic(err)
	}
	Engine.TZLocation, _ = time.LoadLocation("Asia/shanghai")

	f, err1 := os.OpenFile("sql.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err1 != nil {
		panic(err1)
	}
	Engine.SetLogger(xorm.NewSimpleLogger(f))
	Engine.ShowSQL(true)

	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "bcw_")
	Engine.SetTableMapper(tbMapper)

	createTable("bcw_user_info", UserInfo{})
	createTable("bcw_user_address", UserAddress{})
	createTable("bcw_currency_management", CurrencyManagement{})
	createTable("bcw_user_assets", UserAssets{})
	createTable("bcw_trade_recording", TradeRecording{})
	createTable("bcw_user_notify", UserNotify{})
	createTable("bcw_user_feedback", UserFeedback{})
	createTable("bcw_block_height_recording", BlockHeightRecording{})
	createTable("bcw_user_address_management", UserAddressManagement{})
	createTable("bcw_red_packet", RedPacket{})
	createTable("bcw_red_packet_recording", RedPacketRecording{})
}

func createTable(tableName string, tableType interface{}) {

	flag, err := Engine.IsTableExist(tableName)
	if err == nil {
		if !flag {
			err := Engine.CreateTables(tableType)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println(tableName)
		}
	}
}

type operate func(*xorm.Session) error

func SessionSubmit(op operate) (err error) {

	session := Engine.NewSession()
	defer session.Close()

	session.Begin()
	err = op(session)
	if err == nil {
		err = session.Commit()
		if err != nil {
			session.Rollback()
			err = errors.New("operating fail")
		}
	} else {
		session.Rollback()
	}
	return err
}
