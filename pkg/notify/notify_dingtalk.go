package notify

import (
	"dms/config"
	"io"

	"github.com/CodyGuo/dingtalk"
	"github.com/CodyGuo/glog"
)

// var _ NotifyInterface = new(PagerDuty)

type NotifyDingTalk struct{}

// Notify send notify message to DingTalk
func (NotifyDingTalk) Notify(cd config.DingTalk) error {
	glog.SetFlags(glog.LglogFlags)
	webHook := cd.Url
	secret := cd.Secret
	dt := dingtalk.New(webHook, dingtalk.WithSecret(secret))

	// text类型
	textContent := cd.Message
	// atMobiles := robot.SendWithAtMobiles([]string{"176xxxxxx07", "178xxxxxx28"})
	if err := dt.RobotSendText(textContent); err != nil {
		return err
	}
	printResult(dt)
	return nil
}

func printResult(dt *dingtalk.DingTalk) {
	response, err := dt.GetResponse()
	if err != nil {
		glog.Fatal(err)
	}
	reqBody, err := response.Request.GetBody()
	if err != nil {
		glog.Fatal(err)
	}
	reqData, err := io.ReadAll(reqBody)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Infof("发送消息成功, message: %s", reqData)
}
