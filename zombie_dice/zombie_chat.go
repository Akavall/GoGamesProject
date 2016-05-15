package zombie_dice

import (
	"sync"
)

type ZombieChat struct {
	sync.RWMutex
	Messages []string 
}

func (zc *ZombieChat) ThreadSafeAppend(message string) {
	(*zc).Lock()
	defer (*zc).Unlock()
	(*zc).Messages = append((*zc).Messages, message)
}

// I don't think I need to thead safe read
