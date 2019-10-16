package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/crypto/ssh"

	"blockChainWallet/database"

	"os/exec"
)

var netRequestClient = &http.Client{}

/*
 *
 *realname -
 *idcard -
 */
func juheAuthenticationIsMatch(realname, idcard string) (flag bool) {

	urlStr := fmt.Sprintf("", idcard, realname)

	request, _ := http.NewRequest("GET", urlStr, nil)
	response, _ := netRequestClient.Do(request)

	flag = false
	if response.StatusCode == 200 {
		bytes, _ := ioutil.ReadAll(response.Body)
		var etpJson map[string]interface{}
		if err := json.Unmarshal(bytes, &etpJson); err == nil {
			//			fmt.Println(etpJson, "\n", etpJson["error_code"])
			if value, ok := etpJson["error_code"]; ok && value.(float64) == 0 {
				//				fmt.Println("===", ((etpJson["result"]).(map[string]interface{}))["res"])
				if (((etpJson["result"]).(map[string]interface{}))["res"]).(float64) == 1 {
					flag = true
				}
			}
		}
	}
	return flag
}

/*
 *
 *phone -
 *vCode -
 */
func sendTextMessage(phone, vCode string) (flag bool) {

	urlStr := fmt.Sprintf("/?ac=send&uid=leonardwang&pwd=09e24d5d30fadc6f724ed8e2d44adea4&mobile=%s&content=：%s，", phone, vCode)

	request, _ := http.NewRequest("GET", urlStr, nil)
	response, _ := netRequestClient.Do(request)

	flag = false
	if response.StatusCode == 200 {
		bytes, _ := ioutil.ReadAll(response.Body)
		var codeJson map[string]interface{}
		if err := json.Unmarshal(bytes, &codeJson); err == nil {
			fmt.Println(codeJson)
			if value, ok := codeJson["stat"]; ok && value == "100" {
				flag = true
			}
		}
	}
	return flag
}

/*
 *
 * phone -
 */
func phoneIsEffective(phone string) (flag bool) {

	//	reg := regexp.MustCompile(`^((13[0-9])|(14[5,7,9])|(15[^4])|(18[0-9])|(17[0,1,3,5,6,7,8]))\d{8}$`)
	reg := regexp.MustCompile(`^((13[0-9])|(15[^4])|(18[0-9])|(17[0-8])|(145)|(147)|(149)|(166)|(199))\d{8}$`)

	if reg.MatchString(phone) {
		flag = true
	} else {
		flag = false
	}

	return flag
}

/*
 *
 * userId -
 * uniqueMark -
 * aUI -
 */
func alterRedisUserInfo(userId, uniqueMark string, aUI database.UserInfo) {

	value1, err2 := new(database.RedisConn).GetData(userId + uniqueMark)
	if err2 == nil {
		var lm LoginModel
		err2 = json.Unmarshal(value1, &lm)
		lm.Sex = aUI.Sex
		lm.Nickname = aUI.Nickname
		lm.IdentityId = aUI.IdentityId
		lm.Name = aUI.Name
		if err2 == nil {
			insertValue, _ := json.Marshal(lm)
			new(database.RedisConn).InsertData(userId+uniqueMark, insertValue, 7*24*3600)
		}
	}
}

//codemodel
type CodeModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func HandlerCode(code int) (cm CodeModel) {

	cm.Code = code
	switch code {
	case 200:
		cm.Message = ""
	case 1000:
		cm.Message = ""
	case 1001:
		cm.Message = ""
	case 1002:
		cm.Message = ""
	case 1003:
		cm.Message = ""

	case 2000:
		cm.Message = ""
	case 2001:
		cm.Message = ""
	case 2002:
		cm.Message = ""
	case 2003:
		cm.Message = ""
	case 2004:
		cm.Message = ""
	case 2005:
		cm.Message = ""
	case 2006:
		cm.Message = ""
	case 2007:
		cm.Message = ""
	case 2008:
		cm.Message = ""
	case 2009:
		cm.Message = ""
	case 2010:
		cm.Message = ""
	case 2011:
		cm.Message = ""
	case 2012:
		cm.Message = ""

	case 3000:
		cm.Message = ""
	case 3001:
		cm.Message = ""
	case 3002:
		cm.Message = ""
	case 3003:
		cm.Message = ""
	case 3004:
		cm.Message = ""
	case 3005:
		cm.Message = ""
	case 3006:
		cm.Message = ""
	case 3007:
		cm.Message = ""

	case 4000:
		cm.Message = ""
	case 4001:
		cm.Message = ""

	case 6000:
		cm.Message = ""

	case 7000:
		cm.Message = ""

	case 9000:
		cm.Message = ""

	default:
		cm.Message = ""
	}
	return cm
}

type SucceedData struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

func HandlerSucceed(data interface{}) (sd SucceedData) {
	sd.Data = data
	sd.Code = 200
	return sd
}

func connect(user, password, hostPort string) (*ssh.Session, error) {

	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	if client, err = ssh.Dial("tcp", hostPort, clientConfig); err != nil {
		return nil, err
	}

	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil
}

func execCommand(cmdStr string) (string, error) {

	//	var (
	//		serverUserName     string = "root"
	//		serverUserPassWord string = "Hs1234567"
	//		hostPort           string = "39.108.15.102:22"
	//	)

	//
	//	seesion, err := connect(serverUserName, serverUserPassWord, hostPort)
	//	if err != nil {
	//		return "", err
	//	}
	//	defer seesion.Close()

	//
	//	bytes, err := seesion.Output(cmdStr)

	//	return string(bytes), err

	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
