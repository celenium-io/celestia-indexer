// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	"sort"
	"strings"
	"testing"

	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	_ "github.com/celestiaorg/celestia-app/v10/app"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// knownUnhandledEvents are typed events the app can emit that EventType does not
// cover yet; parseEvent stores them as EventTypeUnknown. Listed explicitly so a
// version bump surfaces new events instead of quietly losing their type.
var knownUnhandledEvents = map[string]struct{}{
	"celestia.blob.v1.EventUpdateBlobParams":                                 {},
	"celestia.minfee.v1.EventUpdateMinfeeParams":                             {},
	"hyperlane.core.interchain_security.v1.EventAnnounceStorageLocation":     {},
	"hyperlane.core.interchain_security.v1.EventCreateMerkleRootMultisigIsm": {},
	"hyperlane.core.interchain_security.v1.EventCreateMessageIdMultisigIsm":  {},
	"hyperlane.core.interchain_security.v1.EventRemoveRoutingIsmDomain":      {},
	// Registered but never emitted: the valaddr handler emits the legacy string
	// event set_fibre_provider_info instead, which EventType does cover.
	"celestia.valaddr.v1.EventSetFibreProviderInfo": {},
}

// TestEventCoverage checks every typed event the celestia and hyperlane modules
// register against the EventType enum. Cosmos and IBC modules are out of scope:
// they emit legacy string events, which the enum lists under their own names.
func TestEventCoverage(t *testing.T) {
	files, err := proto.MergedRegistry()
	require.NoError(t, err)

	var registered []string
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		messages := fd.Messages()
		for i := range messages.Len() {
			name := string(messages.Get(i).FullName())
			if !strings.Contains(name, ".Event") {
				continue
			}
			if strings.HasPrefix(name, "celestia.") || strings.HasPrefix(name, "hyperlane.") {
				registered = append(registered, name)
			}
		}
		return true
	})
	require.NotEmpty(t, registered)

	var missing, listedButKnown []string
	for _, name := range registered {
		_, parseErr := storageTypes.ParseEventType(name)
		_, listed := knownUnhandledEvents[name]
		switch {
		case parseErr != nil && !listed:
			missing = append(missing, name)
		case parseErr == nil && listed:
			listedButKnown = append(listedButKnown, name)
		}
	}

	sort.Strings(missing)
	sort.Strings(listedButKnown)
	require.Empty(t, missing, "typed events are not in EventType; add them or list them in knownUnhandledEvents")
	require.Empty(t, listedButKnown, "typed events are in EventType now; drop them from knownUnhandledEvents")
}
