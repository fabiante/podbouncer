package controller

import (
	"sync"
	"time"
)

// PodReconcilerConfig is the configuration object used by PodReconciler.
//
// You should only keep pointers to a PodReconcilerConfig value since it embeds sync.Mutex.
//
// Use only the getter / setter funcs for thread-safe access.
type PodReconcilerConfig struct {
	sync.Mutex

	maxPodAge time.Duration
}

func NewPodReconcilerConfig() *PodReconcilerConfig {
	return &PodReconcilerConfig{
		maxPodAge: time.Hour,
	}
}

func (c *PodReconcilerConfig) SetMaxPodAge(d time.Duration) {
	c.Lock()
	defer c.Unlock()
	c.maxPodAge = d
}

func (c *PodReconcilerConfig) MaxPodAge() time.Duration {
	c.Lock()
	defer c.Unlock()

	return c.maxPodAge
}
