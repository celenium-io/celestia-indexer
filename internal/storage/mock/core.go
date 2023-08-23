package mock

import (
	models "github.com/dipdup-io/celestia-indexer/internal/storage"
	gomock "go.uber.org/mock/gomock"
)

type Storage struct {
	Blocks    models.IBlock
	Tx        models.ITx
	Message   models.IMessage
	Event     models.IEvent
	Address   models.IAddress
	Namespace models.INamespace
	State     models.IState

	ctrl *gomock.Controller
}

func Create(t gomock.TestReporter) Storage {
	ctrl := gomock.NewController(t)
	return Storage{
		Blocks:    NewMockIBlock(ctrl),
		Tx:        NewMockITx(ctrl),
		Message:   NewMockIMessage(ctrl),
		Event:     NewMockIEvent(ctrl),
		Address:   NewMockIAddress(ctrl),
		Namespace: NewMockINamespace(ctrl),
		State:     NewMockIState(ctrl),
		ctrl:      ctrl,
	}
}

func (s Storage) Close() error {
	s.ctrl.Finish()
	return nil
}

func (s Storage) Ctrl() *gomock.Controller {
	return s.ctrl
}
