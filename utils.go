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
	"fmt"
	"strings"
)

// BuildVars - A helper that builds channel variable strings to be included in various commands to FreeSWITCH
func BuildVars(format string, vars map[string]string) string {
	// No vars do not format
	if vars == nil || len(vars) == 0 {
		return ""
	}

	var builder strings.Builder
	for key, value := range vars {
		if builder.Len() > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(key)
		builder.WriteString("=")
		if strings.ContainsAny(value, " ") {
			builder.WriteString("'")
			builder.WriteString(value)
			builder.WriteString("'")
		} else {
			builder.WriteString(value)
		}
	}
	return fmt.Sprintf(format, builder.String())
}
