// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import "github.com/celenium-io/celestia-indexer/internal/storage"

func setNamespacesFromMessage(msg storage.Message, namespaces map[string]*storage.Namespace) {
	for i := range msg.Namespace {
		key := msg.Namespace[i].String()
		if ns, ok := namespaces[key]; !ok {
			msg.Namespace[i].PfbCount = 1
			namespaces[key] = &msg.Namespace[i]
		} else {
			ns.PfbCount += 1
			ns.Size += msg.Namespace[i].Size
		}
	}
}
