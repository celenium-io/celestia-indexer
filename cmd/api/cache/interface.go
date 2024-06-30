// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

type ICache interface {
	Get(key string) ([]byte, bool)
	Set(key string, data []byte)
	Clear()
}
