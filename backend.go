package etclib

import (
	"time"
)

type Backend interface {
	DialTimeout(servers []string, timeout time.Duration) error
	Close()
	IsConnected() bool
	NodeExistsError(err error) bool
	WatchChildren(path string, ch chan []string) (err error)
	Create(path string, value string, flags int32) error
	CreateService(path, value string) error
	CreateOrUpdate(path, value string, flags int32) error
	Exists(path string) (bool, error)
	Set(path, value string) error
	Children(parentKey string) ([]string, error)
	Delete(key string) error
	Get(key string) (string, error)
}
