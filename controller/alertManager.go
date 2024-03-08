package controller

import (
	"dms/config"
	"dms/pkg/logger"
	"dms/pkg/notify"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/alertmanager/template"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//	type Alert struct {
//		*template.Alert
//	}
type Alert = template.Alert
type WatchdogAlertConfig = config.Config

// type WatchdogAlertConfig struct {
// 	*config.Config
// }

type AlertManager struct {
	alertConfigs *WatchdogAlertConfig
	alertTimers  map[string]*time.Timer // 存储不同消息对应的定时器
	alertMux     sync.Mutex
	dingTalkURL  string
	dingTalkMux  sync.Mutex
}

func NewAlertManager(configPath string) (*AlertManager, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}
	env := config.Env
	if env == "" || env == "production" {
		env = "production"
		logger.Log.SetLevel(log.WarnLevel)
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		logger.Log.SetLevel(log.InfoLevel)
	}

	for _, v := range config.Evaluate.Data.Alerts {
		logger.Log.Infof("添加客户%s监控", v)
	}

	return &AlertManager{
		alertConfigs: config,
		alertTimers:  make(map[string]*time.Timer),
		dingTalkURL:  config.Notify.DingTalk.Url,
	}, nil
}

func loadConfig(configPath string) (*WatchdogAlertConfig, error) {
	// 读取配置文件
	orgData, err := os.ReadFile(configPath)
	data := os.ExpandEnv(string(orgData))
	if err != nil {
		return nil, err
	}

	var config WatchdogAlertConfig
	// 解析配置文件
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func constructAlertKey(labels template.KV) string {
	// 将 labels 中的键名排序
	var keys []string
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建排序后的键值对字符串
	var keyValues []string
	for _, k := range keys {
		keyValues = append(keyValues, fmt.Sprintf("%s:%v", k, labels[k]))
	}

	// 使用排序后的键值对字符串作为 alertKey
	return strings.Join(keyValues, "-")
}

func (am *AlertManager) StartTimers() {
	for _, alert := range am.alertConfigs.Evaluate.Data.Alerts {
		alertKey := constructAlertKey(alert.Labels)
		interval := am.alertConfigs.Interval

		am.StartTimer(alertKey, interval)
	}
}

func (am *AlertManager) StartTimer(alertKey string, interval time.Duration) {
	// 创建定时器并启动
	timer := time.NewTimer(interval)
	go func() {
		for {
			select {
			case <-timer.C:
				// fmt.Printf("Alert: Haven't received message for %s\n", alertKey)
				logger.Log.Infof("Alert: Haven't received message for %s\n", alertKey)
				am.sendDingTalkAlert(alertKey)
				am.ResetTimer(alertKey)
			}
		}
	}()
	am.alertMux.Lock()
	defer am.alertMux.Unlock()
	am.alertTimers[alertKey] = timer
}

func (am *AlertManager) sendDingTalkAlert(alertKey string) {
	am.dingTalkMux.Lock()
	defer am.dingTalkMux.Unlock()

	// 发送告警到钉钉
	logger.Log.Infof("Sending DingTalk alert for %s\n", alertKey)
	// fmt.Printf("Sending DingTalk alert for %s\n", alertKey)
	msg := fmt.Sprintf("客户Prometheus已掉线：%s", alertKey)
	dt := *am.alertConfigs.Notify.DingTalk
	dt.Message = msg
	notify.NotifyDingTalk{}.Notify(dt)
	// 在这里添加发送钉钉告警的逻辑，使用am.dingTalkURL
}

func (am *AlertManager) ResetTimer(alertKey string) {
	am.alertMux.Lock()
	defer am.alertMux.Unlock()

	// 重置定时器
	if timer, ok := am.alertTimers[alertKey]; ok {
		timer.Reset(am.alertConfigs.Interval)
	}
}
func struct2json(i interface{}) string {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		logger.Log.Warn(err)
	}
	return string(jsonBytes)
}
func ProcessAlert(c *gin.Context, alertManager *AlertManager) {
	data := new(template.Data)
	if err := c.BindJSON(&data); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		ReturnError(c, http.StatusBadRequest, "error: Invalid JSON")
		return
	}
	// fmt.Fprintf(os.Stdout, "received webhook payload: %s\n", struct2json(data.Alerts))
	logger.Log.Infof("received webhook payload: %s\n", struct2json(data.Alerts))

	for _, v := range data.Alerts {
		alertKey := constructAlertKey(v.Labels)
		alertManager.ResetTimer(alertKey)
	}

	ReturnSuccess(c, http.StatusOK, "message: Received alerts successfully", "")
}
