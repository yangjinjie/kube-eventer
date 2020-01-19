// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"kube-eventer/core"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

const (
	WEBHOOK_SINK          = "WebHookSink"
	WARNING           int = 2
	NORMAL            int = 1
	DEFAULT_MSG_TYPE      = "text"
	CONTENT_TYPE_JSON     = "application/json"
	LABE_TEMPLATE         = "%s\n"
	//发送消息使用的url
	SEND_MSG_URL = `http://172.17.202.22:9999/kube`
)

var (
	MSG_TEMPLATE = "ClusterName: %s\n\nLevel: %s\nKind: %s\nNamespace: %s\nName: %s\nReason: %s\nTimestamp: %s\nMessage: %s"

	MSG_TEMPLATE_ARR = [][]string{
		{"Level"},
		{"Kind"},
		{"ClusterName"},
		{"Namespace"},
		{"Name"},
		{"Reason"},
		{"Timestamp"},
		{"Message"},
	}
)

/**
wechat msg struct
*/
type WebHookMsg struct {
	MsgType string     `json:"msgtype"`
	Text    WechatText `json:"text"`

	ClusterName string `json:"clustername"`
	EventType   string `json:"eventtype"`
	Namespace   string `json:"namespace"`
	Pod         string `json:"pod"`
	Reason      string `json:"reason"`
	Message     string `json:"message"`
}

type WechatText struct {
	Content string `json:"content"`
}

type WebHookSink struct {
	Namespaces []string
	Kinds      []string
	Endpoint   string
	Level      int
	Labels     []string
}

func (d *WebHookSink) Name() string {
	return WEBHOOK_SINK
}

func (d *WebHookSink) Stop() {
	//do nothing
}

func (d *WebHookSink) ExportEvents(batch *core.EventBatch) {
	for _, event := range batch.Events {
		if d.isEventLevelDangerous(event.Type) {
			d.Send(event)
			// add threshold
			time.Sleep(time.Millisecond * 50)
		}
	}
}

func (d *WebHookSink) isEventLevelDangerous(level string) bool {
	score := getLevel(level)
	if score >= d.Level {
		return true
	}
	return false
}

func (d *WebHookSink) Send(event *v1.Event) {
	if d.Namespaces != nil {
		skip := true
		for _, namespace := range d.Namespaces {
			if namespace == event.Namespace {
				skip = false
				break
			}
		}
		if skip {
			return
		}
	}

	if d.Kinds != nil {
		skip := true
		for _, kind := range d.Kinds {
			if kind == event.InvolvedObject.Kind {
				skip = false
				break
			}
		}
		if skip {
			return
		}
	}

	msg := createMsgFromEvent(d, event)
	if msg == nil {
		klog.Warningf("failed to create msg from event,because of %v", event)
		return
	}

	msg_bytes, err := json.Marshal(msg)
	if err != nil {
		klog.Warningf("failed to marshal msg %v", msg)
		return
	}

	b := bytes.NewBuffer(msg_bytes)
	var sendMsgUrl string
	if d.Endpoint != "" {
		sendMsgUrl = d.Endpoint
	} else {
		sendMsgUrl = SEND_MSG_URL
	}
	resp, err := http.Post(sendMsgUrl, CONTENT_TYPE_JSON, b)
	if err != nil {
		klog.Errorf("failed to send msg to webhook. error: %s", err.Error())
		return
	}
	if resp != nil && resp.StatusCode != http.StatusOK {
		klog.Errorf("failed to send msg to webhook, because the response code is %d", resp.StatusCode)
		return
	}

}

func getClusterName() string {
	name := ""
	clusterName := os.Getenv("CLUSTER_NAME")
	if clusterName == "saas-stage" {
		name = clusterName + "🐦"
	} else {
		name = clusterName
	}

	return name
}

func getLevel(level string) int {
	score := 0
	switch level {
	case v1.EventTypeWarning:
		score += 2
	case v1.EventTypeNormal:
		score += 1
	default:
		//score will remain 0
	}
	return score
}

func createMsgFromEvent(d *WebHookSink, event *v1.Event) *WebHookMsg {

	msg := &WebHookMsg{}
	msg.MsgType = DEFAULT_MSG_TYPE

	//默认按文本模式推送
	template := MSG_TEMPLATE
	if len(d.Labels) > 0 {
		for _, label := range d.Labels {
			template = fmt.Sprintf(LABE_TEMPLATE, label) + template
		}
	}
	msg.Namespace = event.Namespace
	msg.Pod = event.Name
	msg.EventType = event.Type
	msg.Reason = event.Reason
	msg.Message = event.Message
	msg.ClusterName = getClusterName()

	msg.Text = WechatText{
		Content: fmt.Sprintf(template, getClusterName(), event.Type, event.InvolvedObject.Kind, event.Namespace, event.Name, event.Reason, event.LastTimestamp.String(), event.Message),
	}

	return msg
}

func NewWebHookSink(uri *url.URL) (*WebHookSink, error) {
	d := &WebHookSink{
		Level: WARNING,
	}
	if len(uri.Host) > 0 {
		d.Endpoint = "http://" + uri.Host + uri.Path
		klog.Info(fmt.Sprintf("endpoint: %s", d.Endpoint))
	} else {
		klog.Info(fmt.Sprintf("endpoint: %s", SEND_MSG_URL))
	}
	opts := uri.Query()

	if len(opts["level"]) >= 1 {
		d.Level = getLevel(opts["level"][0])
	}

	//add extra labels
	if len(opts["label"]) >= 1 {
		d.Labels = opts["label"]
	}

	d.Namespaces = getValues(opts["namespaces"])
	d.Kinds = getValues(opts["kinds"])

	return d, nil
}

func getValues(o []string) []string {
	if len(o) >= 1 {
		if len(o[0]) == 0 {
			return nil
		}
		return strings.Split(o[0], ",")
	}
	return nil
}
