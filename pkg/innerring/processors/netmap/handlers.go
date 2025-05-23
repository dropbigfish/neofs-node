package netmap

import (
	timerEvent "github.com/nspcc-dev/neofs-node/pkg/innerring/timers"
	"github.com/nspcc-dev/neofs-node/pkg/morph/event"
	netmapEvent "github.com/nspcc-dev/neofs-node/pkg/morph/event/netmap"
	"go.uber.org/zap"
)

func (np *Processor) HandleNewEpochTick(ev event.Event) {
	_ = ev.(timerEvent.NewEpochTick)
	np.log.Info("tick", zap.String("type", "epoch"))

	// send an event to the worker pool

	err := np.pool.Submit(func() { np.processNewEpochTick() })
	if err != nil {
		// there system can be moved into controlled degradation stage
		np.log.Warn("netmap worker pool drained",
			zap.Int("capacity", np.pool.Cap()))
	}
}

func (np *Processor) handleNewEpoch(ev event.Event) {
	epochEvent := ev.(netmapEvent.NewEpoch)
	np.log.Info("notification",
		zap.String("type", "new epoch"),
		zap.Uint64("value", epochEvent.EpochNumber()))

	// send an event to the worker pool

	err := np.pool.Submit(func() {
		np.processNewEpoch(epochEvent)
	})
	if err != nil {
		// there system can be moved into controlled degradation stage
		np.log.Warn("netmap worker pool drained",
			zap.Int("capacity", np.pool.Cap()))
	}
}

func (np *Processor) handleAddNode(ev event.Event) {
	newNode := ev.(netmapEvent.AddNode)

	np.log.Info("notification",
		zap.String("type", "add node"),
	)

	// send an event to the worker pool

	err := np.pool.Submit(func() {
		np.processAddNode(newNode)
	})
	if err != nil {
		// there system can be moved into controlled degradation stage
		np.log.Warn("netmap worker pool drained",
			zap.Int("capacity", np.pool.Cap()))
	}
}

func (np *Processor) handleUpdateState(ev event.Event) {
	updPeer := ev.(netmapEvent.UpdatePeer)
	np.log.Info("notification",
		zap.String("type", "update peer state"),
		zap.String("key", updPeer.PublicKey().StringCompressed()))

	// send event to the worker pool

	err := np.pool.Submit(func() {
		np.processUpdatePeer(updPeer)
	})
	if err != nil {
		// there system can be moved into controlled degradation stage
		np.log.Warn("netmap worker pool drained",
			zap.Int("capacity", np.pool.Cap()))
	}
}
