// Package store define the Store interface with basic operations of a key value store.
// The store package also defines a MemoryStore with fields to track access and modify time
// with their respective mutexes so that access can be used to clean keys with low rate of
// use and modify time is used to stream keys in last modified order.
package store
