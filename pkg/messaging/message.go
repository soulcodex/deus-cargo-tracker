package messaging

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	DefaultMessageSpecVersion = "1.0"
)

type Message interface {
	MessageIdentity

	Time() time.Time
	DataRaw() []byte
	Metadata() map[string]any
}

type MessageSubject interface {
	Subject() string
	SubjectName() string
	SubjectID() string
}

type MessageIdentity interface {
	MessageSubject

	Type() string
	SpecVersion() string
	Source() string
	Identifier() string
}

// AckAwaiter is an interface that exposes the methods to wait for an ack or nack.
type AckAwaiter interface {
	Acked() <-chan struct{}
	Nacked() <-chan struct{}
}

// Acknowledgeable exposes the methods to acknowledge or reject a message.
type Acknowledgeable interface {
	Ack()
	Nack()
}
type BaseMessage struct {
	msgID          string
	msgType        string
	msgSource      string
	msgSubjectID   string
	msgSubjectName string
	msgSpecVersion string
	msgTime        time.Time

	msgData     []byte
	msgMetadata map[string]any

	ackOnce  sync.Once
	nackOnce sync.Once
	ack      chan struct{}
	noAck    chan struct{}
}

func NewBaseMessage(
	msgType string,
	msgSpecVersion string,
	msgSource string,
	msgSubjectID string,
	msgSubjectName string,
	msgTime time.Time,
	msgData []byte,
) *BaseMessage {
	return &BaseMessage{
		msgID:          defaultIdentityProvider.Provide(),
		msgType:        msgType,
		msgSpecVersion: msgSpecVersion,
		msgSource:      msgSource,
		msgSubjectID:   msgSubjectID,
		msgSubjectName: msgSubjectName,
		msgTime:        msgTime,
		msgData:        msgData,
		msgMetadata:    make(map[string]any),
		ackOnce:        sync.Once{},
		nackOnce:       sync.Once{},
		ack:            make(chan struct{}),
		noAck:          make(chan struct{}),
	}
}

func (bm *BaseMessage) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	alias := new(baseMessageJSONAlias[json.RawMessage])
	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("unmarshal base message: %w", err)
	}

	err := bm.fromJSONAlias(alias)
	if err != nil {
		return err
	}

	return nil
}

func (bm *BaseMessage) MarshalJSON() ([]byte, error) {
	alias := new(baseMessageJSONAlias[json.RawMessage])
	alias.fromBaseMessage(bm)
	data, err := json.Marshal(alias)
	if err != nil {
		return nil, fmt.Errorf("marshal base message: %w", err)
	}
	return data, nil
}

func (bm *BaseMessage) fromJSONAlias(alias *baseMessageJSONAlias[json.RawMessage]) error {
	subjectID, subjectName, err := ParseSubject(alias.MsgSubject)
	if err != nil {
		return err
	}

	// Copy from alias
	bm.msgType = alias.MsgType
	bm.msgSpecVersion = alias.MsgSpecVersion
	bm.msgSource = alias.MsgSource
	bm.msgSubjectID = subjectID
	bm.msgSubjectName = subjectName
	bm.msgID = alias.MsgID
	bm.msgTime = alias.MsgTime
	bm.msgData = alias.MsgDataRaw

	if bm.ack == nil {
		bm.ack = make(chan struct{})
	}

	if bm.noAck == nil {
		bm.noAck = make(chan struct{})
	}

	return nil
}

func (bm *BaseMessage) Type() string {
	return bm.msgType
}

func (bm *BaseMessage) SpecVersion() string {
	return bm.msgSpecVersion
}

func (bm *BaseMessage) Source() string {
	return bm.msgSource
}

func (bm *BaseMessage) Subject() string {
	if bm.msgSubjectName == "" || bm.msgSubjectID == "" {
		return ""
	}

	return bm.msgSubjectName + "." + bm.msgSubjectID
}

func (bm *BaseMessage) SubjectName() string {
	return bm.msgSubjectName
}

func (bm *BaseMessage) SubjectID() string {
	return bm.msgSubjectID
}

func (bm *BaseMessage) Identifier() string {
	return bm.msgID
}

func (bm *BaseMessage) Time() time.Time {
	return bm.msgTime
}

func (bm *BaseMessage) DataRaw() []byte {
	return bm.msgData
}

func (bm *BaseMessage) Metadata() map[string]any {
	return bm.msgMetadata
}

// Ack marks the message as acknowledged
func (bm *BaseMessage) Ack() {
	bm.ackOnce.Do(func() {
		select {
		case <-bm.ack:
			// Channel is already closed
		default:
			close(bm.ack)
		}
	})
}

// Nack marks the message as not acknowledged
func (bm *BaseMessage) Nack() {
	bm.nackOnce.Do(func() {
		select {
		case <-bm.noAck:
			// Channel is already closed
		default:
			close(bm.noAck)
		}
	})
}

// Acked returns a channel that is closed when the message is acknowledged
func (bm *BaseMessage) Acked() <-chan struct{} {
	return bm.ack
}

// Nacked returns a channel that is closed when the message is not acknowledged
func (bm *BaseMessage) Nacked() <-chan struct{} {
	return bm.noAck
}

// ParseSubject splits the subject into name and id.
// The subject must be in the format {name}.{id}.
//
// Example: ParseSubject("subject.name") = "subject", "name"
func ParseSubject(subject string) (string, string, error) {
	const subjectPartsCount = 2
	subjectParts := strings.Split(subject, ".")
	if len(subjectParts) != subjectPartsCount {
		return "", "", fmt.Errorf("invalid subject format: %s", subject)
	}

	return subjectParts[1], subjectParts[0], nil
}
