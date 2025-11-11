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
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/shuguocloud/eslgo/command"
	"github.com/shuguocloud/eslgo/command/call"
)

func (c *Conn) EnableEvents(ctx context.Context, format ...string) error {
	var err error
	eventFormat := "plain" // default to plain text
	if len(format) > 0 && format[0] != "" {
		eventFormat = format[0]
	}
	if c.outbound {
		_, err = c.SendCommand(ctx, command.MyEvents{
			Format: eventFormat,
		})
	} else {
		_, err = c.SendCommand(ctx, command.Event{
			Format: eventFormat,
			Listen: []string{"all"},
		})
	}
	return err
}

// DebugEvents - A helper that will output all events to a logger
func (c *Conn) DebugEvents(w io.Writer) string {
	logger := log.New(w, "EventLog: ", log.LstdFlags|log.Lmsgprefix)
	return c.RegisterEventListener(EventListenAll, func(event *Event) {
		logger.Println(event)
	})
}

func (c *Conn) DebugOff(id string) {
	c.RemoveEventListener(EventListenAll, id)
}

// Phrase - Executes the mod_dptools phrase app
func (c *Conn) Phrase(ctx context.Context, uuid, macro string, times int) (*RawResponse, error) {
	return c.audioCommand(ctx, "phrase", uuid, macro, times, true)
}

// PhraseAsync - Executes the mod_dptools phrase app with async mode
func (c *Conn) PhraseAsync(ctx context.Context, uuid, macro string, times int) (*RawResponse, error) {
	return c.audioCommand(ctx, "phrase", uuid, macro, times, false)
}

// PhraseWithArg - Executes the mod_dptools phrase app with arguments
func (c *Conn) PhraseWithArg(ctx context.Context, uuid, macro string, argument interface{}, times int) (*RawResponse, error) {
	return c.audioCommand(ctx, "phrase", uuid, fmt.Sprintf("%s,%v", macro, argument), times, true)
}

// PhraseWithArgAsync - Executes the mod_dptools phrase app with arguments by async mode
func (c *Conn) PhraseWithArgAsync(ctx context.Context, uuid, macro string, argument interface{}, times int) (*RawResponse, error) {
	return c.audioCommand(ctx, "phrase", uuid, fmt.Sprintf("%s,%v", macro, argument), times, false)
}

// Playback - Executes the mod_dptools playback app
func (c *Conn) Playback(ctx context.Context, uuid, audioArgs string, times int) (*RawResponse, error) {
	return c.audioCommand(ctx, "playback", uuid, audioArgs, times, true)
}

// PlaybackAsync - Executes the mod_dptools playback app with async mode
func (c *Conn) PlaybackAsync(ctx context.Context, uuid, audioArgs string, times int) (*RawResponse, error) {
	return c.audioCommand(ctx, "playback", uuid, audioArgs, times, false)
}

// Say - Executes the mod_dptools say app
func (c *Conn) Say(ctx context.Context, uuid, audioArgs string, times int) (*RawResponse, error) {
	return c.audioCommand(ctx, "say", uuid, audioArgs, times, true)
}

// SayAsync - Executes the mod_dptools say app with async mode
func (c *Conn) SayAsync(ctx context.Context, uuid, audioArgs string, times int) (*RawResponse, error) {
	return c.audioCommand(ctx, "say", uuid, audioArgs, times, false)
}

// Speak - Executes the mod_dptools speak app
func (c *Conn) Speak(ctx context.Context, uuid, audioArgs string, times int) (*RawResponse, error) {
	return c.audioCommand(ctx, "speak", uuid, audioArgs, times, true)
}

// SpeakAsync - Executes the mod_dptools speak app with async mode
func (c *Conn) SpeakAsync(ctx context.Context, uuid, audioArgs string, times int) (*RawResponse, error) {
	return c.audioCommand(ctx, "speak", uuid, audioArgs, times, false)
}

// Execute - Executes the mod_dptools app
func (c *Conn) Execute(ctx context.Context, uuid, command, appArgs string) (*RawResponse, error) {
	return c.executeCommand(ctx, command, uuid, appArgs, true)
}

// ExecuteAsync - Executes the mod_dptools app with async mode
func (c *Conn) ExecuteAsync(ctx context.Context, uuid, command, appArgs string) (*RawResponse, error) {
	return c.executeCommand(ctx, command, uuid, appArgs, false)
}

// Hangup - Executes the mod_dptools hangup app
func (c *Conn) Hangup(ctx context.Context, uuid, cause string) (*RawResponse, error) {
	return c.hangupCommand(ctx, uuid, cause, true)
}

// HangupAsync - Executes the mod_dptools hangup app with async mode
func (c *Conn) HangupAsync(ctx context.Context, uuid, cause string) (*RawResponse, error) {
	return c.hangupCommand(ctx, uuid, cause, false)
}

// Answer - Executes the mod_dptools answer app
func (c *Conn) Answer(ctx context.Context, uuid string) (*RawResponse, error) {
	return c.executeCommand(ctx, "answer", uuid, "", true)
}

// AnswerAsync - Executes the mod_dptools answer app with async mode
func (c *Conn) AnswerAsync(ctx context.Context, uuid string) (*RawResponse, error) {
	return c.executeCommand(ctx, "answer", uuid, "", false)
}

// Set - Executes the mod_dptools set app
func (c *Conn) Set(ctx context.Context, uuid, key, value string) (*RawResponse, error) {
	return c.setCommand(ctx , uuid, key, value, true)
}

// SetAsync - Executes the mod_dptools set app with async mode
func (c *Conn) SetAsync(ctx context.Context, uuid, key, value string)(*RawResponse, error)   {
	return c.setCommand(ctx , uuid, key, value, false)
}

// Linger - Executes the mod_dptools linger app
func (c *Conn) Linger(ctx context.Context, enabled bool, time time.Duration) (*RawResponse, error) {
	response, err := c.SendCommand(ctx, command.Linger{
		Enabled: enabled,
		Seconds: time,
	})
	if err != nil {
		return response, err
	}
	if !response.IsOk() {
		return response, errors.New("linger response is not okay")
	}
	return response, nil
}

// Conference - Executes the mod_dptools conference app
func (c *Conn) Conference(ctx context.Context, uuid, audioArgs string) (*RawResponse, error) {
	return c.executeCommand(ctx, "conference", uuid, audioArgs, true)
}

// ConferenceAsync - Executes the mod_dptools conference app with async mode
func (c *Conn) ConferenceAsync(ctx context.Context, uuid, audioArgs string) (*RawResponse, error) {
	return c.executeCommand(ctx, "conference", uuid, audioArgs, false)
}

// WaitForDTMF, waits for a DTMF event. Requires events to be enabled!
func (c *Conn) WaitForDTMF(ctx context.Context, uuid string) (byte, error) {
	done := make(chan byte, 1)
	listenerID := c.RegisterEventListener(uuid, func(event *Event) {
		if event.GetName() == "DTMF" {
			dtmf := event.GetHeader("DTMF-Digit")
			if len(dtmf) > 0 {
				select {
				case done <- dtmf[0]:
				default:
				}
			} else {
				select {
				case done <- 0:
				default:
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	})
	defer func() {
		c.RemoveEventListener(uuid, listenerID)
		close(done)
	}()

	select {
	case digit := <-done:
		if digit != 0 {
			return digit, nil
		}
		return digit, errors.New("invalid DTMF digit received")
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

// Helper for mod_dptools apps since they are very similar in invocation
func (c *Conn) audioCommand(ctx context.Context, command, uuid, audioArgs string, times int, wait bool) (*RawResponse, error) {
	response, err := c.SendCommand(ctx, &call.Execute{
		UUID:    uuid,
		AppName: command,
		AppArgs: audioArgs,
		Loops:   times,
		Sync:    wait,
	})
	if err != nil {
		return response, err
	}
	if !response.IsOk() {
		return response, errors.New(command + " response is not okay")
	}
	return response, nil
}

// Helper fuck to execute commands with its args and sync/async mode
func (c *Conn) executeCommand(ctx context.Context, command, uuid, audioArgs string, wait bool) (*RawResponse, error) {
	return c.audioCommand(ctx, command, uuid, audioArgs, 0, wait)
}

// Helper fuck to set commands with its args and sync/async mode
func (c *Conn) setCommand(ctx context.Context, uuid, key, value string, wait bool) (*RawResponse, error) {
	response, err := c.SendCommand(ctx, &call.Set{
		UUID:  uuid,
		Key:   key,
		Value: value,
		Sync:  wait,
	})
	if err != nil {
		return response, err
	}
	if !response.IsOk() {
		return response, errors.New("set response is not okay")
	}
	return response, nil
}

func (c *Conn) hangupCommand(ctx context.Context, uuid, cause string, wait bool) (*RawResponse, error) {
	response, err := c.SendCommand(ctx, call.Hangup{
		UUID:  uuid,
		Cause: cause,
		Sync:  wait,
	})
	if err != nil {
		return response, err
	}
	if !response.IsOk() {
		return response, errors.New("set response is not okay")
	}
	return response, nil
}