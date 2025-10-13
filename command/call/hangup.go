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

type Hangup struct {
	UUID    string
	Cause   string
	Sync    bool
	SyncPri bool
}

func (h Hangup) BuildMessage() string {
	sendMsg := command.SendMessage{
		UUID:    h.UUID,
		Headers: make(textproto.MIMEHeader),
		Sync:    h.Sync,
		SyncPri: h.SyncPri,
	}
	sendMsg.Headers.Set("call-command", "hangup")
	sendMsg.Headers.Set("hangup-cause", h.Cause)

	return sendMsg.BuildMessage()
}
