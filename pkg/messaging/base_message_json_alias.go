package messaging

import (
	"encoding/json"
	"time"
)

type baseMessageJSONAlias[T json.RawMessage | []byte] struct {
	MsgType        string    `json:"type"`
	MsgSpecVersion string    `json:"specversion"`
	MsgSource      string    `json:"source"`
	MsgSubject     string    `json:"subject"`
	MsgID          string    `json:"id"`
	MsgTime        time.Time `json:"time"`
	MsgDataRaw     T         `json:"data"`
}

func (bms *baseMessageJSONAlias[T]) fromBaseMessage(bm *BaseMessage) {
	bms.MsgType = bm.msgType
	bms.MsgSpecVersion = bm.msgSpecVersion
	bms.MsgSource = bm.msgSource
	bms.MsgSubject = bm.msgSubjectName + "." + bm.msgSubjectID
	bms.MsgID = bm.msgID
	bms.MsgTime = bm.msgTime
	bms.MsgDataRaw = T(bm.msgData)
}
