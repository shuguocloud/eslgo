/*
 * Copyright (c) 2020 Opensmarty
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Contributor(s):
 * Opensmarty  <opensmarty@163.com>
 */
package call

import (
	"net/textproto"

	"github.com/shuguocloud/eslgo/command"
)

type NoMedia struct {
	UUID        string
	NoMediaUUID string
	Sync        bool
	SyncPri     bool
}

func (n NoMedia) BuildMessage() string {
	sendMsg := command.SendMessage{
		UUID:    n.UUID,
		Headers: make(textproto.MIMEHeader),
		Sync:    n.Sync,
		SyncPri: n.SyncPri,
	}
	sendMsg.Headers.Set("call-command", "nomedia")
	sendMsg.Headers.Set("nomedia-uuid", n.NoMediaUUID)

	return sendMsg.BuildMessage()
}
