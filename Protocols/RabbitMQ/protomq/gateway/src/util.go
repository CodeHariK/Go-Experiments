package main

import (
	spec "gitlab.com/rahasak-labs/rabbit/spec"
	"strconv"
	"time"
)

func docMsg(name string) *spec.CreateDocumentMessage {
	uid := uid()
	doc := &spec.Document{
		Id:        "1945",
		Name:      name,
		Timestamp: timestamp(),
	}
	msg := &spec.CreateDocumentMessage{
		Uid:      uid,
		Document: doc,
	}

	return msg
}

func uid() string {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	return "ops" + strconv.FormatInt(t, 10)
}

func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
