package jpushclient

const (
	appKey = "5a5b5a572707c5e030125e3e"
	secret = "e86e87c301cdddbee1fc9b33"
)

func SendJPush(sendStr, registrationId, platform string) (err error) {

	//Platform
	var pf Platform
	//Notice
	var notice Notice
	if platform == "iOS" {
		pf.Add(IOS)
		notice.SetIOSNotice(&IOSNotice{Alert: sendStr, Badge: "1", Sound: "defalut"})
	} else if platform == "Android" {
		pf.Add(ANDROID)
		notice.SetAndroidNotice(&AndroidNotice{Alert: sendStr})
	}

	//Audience
	var ad Audience
	s := []string{registrationId}
	ad.SetID(s)

	payload := NewPushPayLoad()
	payload.SetPlatform(&pf)
	payload.SetAudience(&ad)
	payload.SetNotice(&notice)

	bytes, _ := payload.ToBytes()

	//push
	client := NewPushClient(secret, appKey)
	_, err = client.Send(bytes)

	return err
}
