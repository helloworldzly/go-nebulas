// Copyright (C) 2017 go-nebulas authors
//
// This file is part of the go-nebulas library.
//
// the go-nebulas library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-nebulas library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-nebulas library.  If not, see <http://www.gnu.org/licenses/>.
//

package net

// MessageTypeNewBlock
const (
	MessageTypeNewBlock = "NewBlockMessage"
	MessageTypeNewTx    = "NewTxMessage"
)

// MessageType a string for message type.
type MessageType string

// Message interface for message.
type Message interface {
	MessageType() MessageType
	Data() interface{}
}

// Manager manager interface
type Manager interface {
	Start()
	Stop()

	Register(subscribers ...*Subscriber)
	Deregister(subscribers ...*Subscriber)

	BroadcastBlock(block interface{})
}

// Subscriber subscriber.
type Subscriber struct {
	// id usually the owner/creator, used for troubleshooting .
	id interface{}

	// msgChan chan for subscribed message.
	msgChan chan Message

	// msgType message types to subscribe
	msgTypes []MessageType
}

// NewSubscriber return new Subscriber instance.
func NewSubscriber(id interface{}, msgChan chan Message, msgTypes ...MessageType) *Subscriber {
	return &Subscriber{id, msgChan, msgTypes}
}

// ID return id.
func (s *Subscriber) ID() interface{} {
	return s.id
}

// MessageType return msgTypes.
func (s *Subscriber) MessageType() []MessageType {
	return s.msgTypes
}

// MessageChan return msgChan.
func (s *Subscriber) MessageChan() chan Message {
	return s.msgChan
}
