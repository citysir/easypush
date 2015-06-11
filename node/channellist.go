package main

import (
	"errors"
	"github.com/citysir/easypush/hash"
	log "github.com/log4go"
	"sync"
)

var (
	ErrChannelNotExist = errors.New("Channel: not exist")
	ErrChannelKey      = errors.New("Channel: Key not belong to this node")
	nodeHash           = hash.NewNodeHash(Conf.NodeSlotMap)
	nodeWeightMap      = map[string]int{}
)

// Connection bucket.
type ChannelBucket struct {
	data  map[string]*Channel
	mutex *sync.Mutex
}

// Connection list.
type ChannelList struct {
	channels []*ChannelBucket
}

// Lock lock the bucket mutex.
func (c *ChannelBucket) Lock() {
	c.mutex.Lock()
}

// Unlock unlock the bucket mutex.
func (c *ChannelBucket) Unlock() {
	c.mutex.Unlock()
}

// NewChannelList create a new conn bucket set.
func NewChannelList() *ChannelList {
	l := &ChannelList{channels: []*ChannelBucket{}}
	// split hashmap to many bucket
	log.Debug("create %d ChannelBucket", Conf.ChannelBucketCount)
	for i := 0; i < Conf.ChannelBucketCount; i++ {
		b := &ChannelBucket{
			data:  map[string]*Channel{},
			mutex: &sync.Mutex{},
		}
		l.channels = append(l.channels, b)
	}
	return l
}

// Count get the bucket total conn count.
func (l *ChannelList) Count() int {
	c := 0
	for i := 0; i < Conf.ChannelBucketCount; i++ {
		c += len(l.channels[i].data)
	}
	return c
}

// bucket return a ChannelBucket use murmurhash3.
func (l *ChannelList) Bucket(key string) *ChannelBucket {
	idx := int(hash.HashCrc32String(key) % uint32(Conf.ChannelBucketCount))
	return l.channels[idx]
}

// validate check the key is belong to this node.
func (l *ChannelList) validate(key string) error {
	if len(nodeWeightMap) == 0 {
		log.Debug("no node found")
		return ErrChannelKey
	}
	nodeId := nodeHash.Hash(key)
	log.Debug("match node:%s hash node:%s", Conf.ZookeeperNodeId, nodeId)
	if Conf.ZookeeperNodeId != nodeId {
		log.Warn("User: %s node:%s not match this node:%s", key, nodeId, Conf.ZookeeperNodeId)
		return ErrChannelKey
	}
	return nil
}

// New create a user conn.
func (l *ChannelList) Add(key string, c *Channel) error {
	// add a new conn
	b := l.Bucket(key)
	b.Lock()
	if _, ok := b.data[key]; ok {
		b.Unlock()
		log.Info("User: %s refresh conn bucket expire time", key)
		return nil
	} else {
		b.data[key] = c
		b.Unlock()
		log.Info("User: %s add a new conn", key)
		ChStat.IncrCreate()
		return nil
	}
}

// Get a user Channel from ChannelList.
func (l *ChannelList) Get(key string) *Channel {
	// get a conn bucket
	b := l.Bucket(key)
	b.Lock()
	if c, ok := b.data[key]; ok {
		b.Unlock()
		ChStat.IncrAccess()
		return c
	}
	b.Unlock()
	return nil
}

// Delete a user conn from ChannelList.
func (l *ChannelList) Remove(key string) *Channel {
	// get a conn bucket
	b := l.Bucket(key)
	b.Lock()
	if c, ok := b.data[key]; ok {
		delete(b.data, key)
		b.Unlock()
		c.Close()
		log.Info("User: %s removed conn", key)
		ChStat.IncrRemove()
		return c
	}
	b.Unlock()
	return nil
}

// Close close all conn.
func (l *ChannelList) Close() {
	channels := make([]*Channel, 0, l.Count())
	for _, c := range l.channels {
		c.Lock()
		for _, c := range c.data {
			channels = append(channels, c)
		}
		c.Unlock()
	}
	// close all channels
	for _, c := range channels {
		if err := c.Close(); err != nil {
			log.Error("c.Close() error(%v)", err)
		}
	}
}
