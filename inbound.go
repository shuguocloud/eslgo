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
package eslgo

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/shuguocloud/eslgo/command"
)

// InboundOptions - Used to dial a new inbound ESL connection to FreeSWITCH
type InboundOptions struct {
	Options                    // Generic common options to both Inbound and Outbound Conn
	Network      string        // The network type to use, should always be tcp, tcp4, tcp6.
	Password     string        // The password used to authenticate with FreeSWITCH. Usually ClueCon
	OnDisconnect func()        // An optional function to be called with the inbound connection gets disconnected
	AuthTimeout  time.Duration // How long to wait for authentication to complete
}

// DefaultOutboundOptions - The default options used for creating the inbound connection
var DefaultInboundOptions = InboundOptions{
	Options:     DefaultOptions,
	Network:     "tcp",
	Password:    "ClueCon",
	AuthTimeout: 5 * time.Second,
}

// Dial - Connects to FreeSWITCH ESL at the provided address and authenticates with the provided password. onDisconnect is called when the connection is closed either by us, FreeSWITCH, or network error
func Dial(address, password string, onDisconnect func()) (*Conn, error) {
	opts := DefaultInboundOptions
	opts.Password = password
	opts.OnDisconnect = onDisconnect
	return opts.Dial(address)
}

// Dial - Connects to FreeSWITCH ESL on the address with the provided options. Returns the connection and any errors encountered
func (opts InboundOptions) Dial(address string) (*Conn, error) {
	// 建立 TCP 连接加超时，避免对端不可达时 Dial 长时间阻塞导致重连 goroutine 挂起
	dialer := net.Dialer{Timeout: opts.AuthTimeout}
	c, err := dialer.Dial(opts.Network, address)
	if err != nil {
		return nil, err
	}
	connection := newConnection(c, false, opts.Options)

	// First auth：等待 auth/request 加超时兜底，避免对端半开/卡死时永久阻塞（reconnectFlag 无法复位）
	// 快照 channel 到局部变量后再 select，避免与 close() 并发读写 responseChannels map
	// （concurrent map read and map write 是 fatal error，recover 无法捕获）。
	// 若 close() 已执行，快照为 nil，由 time.After 超时分支兜底（ExitAndClose 幂等）。
	authRequestChan := connection.snapshotResponseChan(TypeAuthRequest)
	select {
	case <-authRequestChan:
	case <-time.After(opts.AuthTimeout):
		connection.ExitAndClose()
		return nil, fmt.Errorf("timeout waiting for auth/request from %s", address)
	}
	authCtx, cancel := context.WithTimeout(connection.runningContext, opts.AuthTimeout)
	err = connection.doAuth(authCtx, command.Auth{Password: opts.Password})
	cancel()
	if err != nil {
		// Try to gracefully disconnect, we have the wrong password.
		connection.ExitAndClose()
		if opts.OnDisconnect != nil {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						connection.logger.Warn("panic in onDisconnect: %v", r)
					}
				}()
				opts.OnDisconnect()
			}()
		}
		return nil, err
	} else {
		connection.logger.Info("Successfully authenticated %s\n", connection.conn.RemoteAddr())
	}

	// Inbound only handlers
	go connection.authLoop(command.Auth{Password: opts.Password}, opts.AuthTimeout)
	go connection.disconnectLoop(opts.OnDisconnect)

	return connection, nil
}

func (c *Conn) disconnectLoop(onDisconnect func()) {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Warn("panic in disconnectLoop: %v", r)
		}
	}()
	// Snapshot the channels under the read-lock: raw map access here would race
	// with close() mutating responseChannels (concurrent map read and map write).
	// If close() already ran, the snapshot is nil and runningContext (cancelled
	// by close()'s stopFunc) terminates the select below.
	disconnectChan := c.snapshotResponseChan(TypeDisconnect)
	select {
	case <-disconnectChan:
		c.Close()
		if onDisconnect != nil {
			onDisconnect()
		}
		return
	case <-c.runningContext.Done():
		// Connection was closed by another path (panic-recovery Close, ExitAndClose
		// after authLoop failure, ...). close() cancels runningContext before closing
		// the disconnect channel, so this branch can win the select race; it must
		// still notify the caller or an inbound connection would never reconnect.
		// Callers must make onDisconnect idempotent (fsctl uses a CAS-guarded
		// triggerReconnect, so duplicate notification is harmless).
		if onDisconnect != nil {
			onDisconnect()
		}
		return
	}
}

func (c *Conn) authLoop(auth command.Auth, authTimeout time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Warn("panic in authLoop: %v", r)
		}
	}()
	for {
		// Snapshot under the read-lock, same race protection as disconnectLoop
		authRequestChan := c.snapshotResponseChan(TypeAuthRequest)
		select {
		case <-authRequestChan:
			authCtx, cancel := context.WithTimeout(c.runningContext, authTimeout)
			err := c.doAuth(authCtx, auth)
			cancel()
			if err != nil {
				c.logger.Warn("Failed to auth %s\n", err.Error())
				// Close the connection, we have the wrong password
				c.ExitAndClose()
				return
			} else {
				c.logger.Info("Successfully authenticated %s\n", c.conn.RemoteAddr())
			}
		case <-c.runningContext.Done():
			return
		}
	}
}

func (c *Conn) doAuth(ctx context.Context, auth command.Auth) error {
	response, err := c.SendCommand(ctx, auth)
	if err != nil {
		return err
	}
	if !response.IsOk() {
		return fmt.Errorf("failed to auth %#v", response)
	}
	return nil
}
