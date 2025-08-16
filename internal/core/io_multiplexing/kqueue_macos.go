//go:build darwin

package io_multiplexing

import (
	"log"
	"syscall"

	"github.com/lyxuansang91/redis-crash-course/internal/config"
)

type KQueue struct {
	fd            int
	kqEvents      []syscall.Kevent_t
	genericEvents []Event
}

func CreateIOMultiplexer(config *config.Config) (*KQueue, error) {
	epollFD, err := syscall.Kqueue()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &KQueue{
		fd:            epollFD,
		kqEvents:      make([]syscall.Kevent_t, config.MaxConnections),
		genericEvents: make([]Event, config.MaxConnections),
	}, nil
}

func (kq *KQueue) Monitor(event Event) error {
	kqEvent := event.toNative(syscall.EV_ADD)
	// Add event.Fd to the monitoring list of kq.fd
	_, err := syscall.Kevent(kq.fd, []syscall.Kevent_t{kqEvent}, nil, nil)
	return err
}

func (kq *KQueue) Wait() ([]Event, error) {
	n, err := syscall.Kevent(kq.fd, nil, kq.kqEvents, nil)
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		kq.genericEvents[i] = createEvent(kq.kqEvents[i])
	}

	return kq.genericEvents[:n], nil
}

func (kq *KQueue) Close() error {
	return syscall.Close(kq.fd)
}