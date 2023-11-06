package store

import (
	"errors"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

func NewMemoryStore(opts ...Option) Store {
	s := &memoryStore{
		options: Options{
			Database: "micro",
			Table:    "micro",
		},
		store: cache.New(cache.NoExpiration, 5*time.Minute),
	}
	for _, o := range opts {
		o(&s.options)
	}
	return s
}

type memoryStore struct {
	options Options

	store *cache.Cache
}

type storeRecord struct {
	key       string
	value     []byte
	metadata  map[string]interface{}
	expiresAt time.Time
}

func (m *memoryStore) key(prefix, key string) string {
	return filepath.Join(prefix, key)
}

func (m *memoryStore) prefix(database, table string) string {
	if len(database) == 0 {
		database = m.options.Database
	}
	if len(table) == 0 {
		table = m.options.Table
	}
	return filepath.Join(database, table)
}

func (m *memoryStore) get(prefix, key string) (*Record, error) {
	key = m.key(prefix, key)

	var storedRecord *storeRecord
	r, found := m.store.Get(key)
	if !found {
		return nil, ErrNotFound
	}

	storedRecord, ok := r.(*storeRecord)
	if !ok {
		return nil, errors.New("Retrieved a non *storeRecord from the cache")
	}

	newRecord := &Record{}
	newRecord.Key = strings.TrimPrefix(storedRecord.key, prefix+"/")
	newRecord.Value = make([]byte, len(storedRecord.value))
	newRecord.Metadata = make(map[string]interface{})

	copy(newRecord.Value, storedRecord.value)

	if !storedRecord.expiresAt.IsZero() {
		newRecord.Expiry = time.Until(storedRecord.expiresAt)
	}

	for k, v := range storedRecord.metadata {
		newRecord.Metadata[k] = v
	}

	return newRecord, nil
}

func (m *memoryStore) set(prefix string, r *Record) {
	key := m.key(prefix, r.Key)

	i := &storeRecord{}
	i.key = r.Key
	i.value = make([]byte, len(r.Value))
	i.metadata = make(map[string]interface{})

	copy(i.value, r.Value)

	if r.Expiry != 0 {
		i.expiresAt = time.Now().Add(r.Expiry)
	}

	for k, v := range r.Metadata {
		i.metadata[k] = v
	}

	m.store.Set(key, i, r.Expiry)
}

func (m *memoryStore) delete(prefix, key string) {
	key = m.key(prefix, key)
	m.store.Delete(key)
}

func (m *memoryStore) list(prefix string, limit, offset uint) []string {
	allItems := m.store.Items()
	foundKeys := make([]string, 0, len(allItems))

	for k := range allItems {
		if !strings.HasPrefix(k, prefix+"/") {
			continue
		}
		foundKeys = append(foundKeys, strings.TrimPrefix(k, prefix+"/"))
	}

	if limit != 0 || offset != 0 {
		sort.Slice(foundKeys, func(i, j int) bool { return foundKeys[i] < foundKeys[j] })
		min := func(i, j uint) uint {
			if i < j {
				return i
			}
			return j
		}
		return foundKeys[offset:min(limit, uint(len(foundKeys)))]
	}

	return foundKeys
}

func (m *memoryStore) Close() error {
	m.store.Flush()
	return nil
}

func (m *memoryStore) Init(opts ...Option) error {
	for _, o := range opts {
		o(&m.options)
	}
	return nil
}

func (m *memoryStore) String() string {
	return "memory"
}

func (m *memoryStore) Read(key string, opts ...ReadOption) ([]*Record, error) {
	readOpts := ReadOptions{}
	for _, o := range opts {
		o(&readOpts)
	}

	prefix := m.prefix(readOpts.Database, readOpts.Table)

	var keys []string
	if readOpts.Prefix || readOpts.Suffix {
		// 获取所有的key
		k := m.list(prefix, 0, 0)
		limit := int(readOpts.Limit)
		offset := int(readOpts.Offset)

		if limit > len(k) {
			limit = len(k)
		}

		if offset > len(k) {
			offset = len(k)
		}

		for i := offset; i < limit; i++ {
			kk := k[i]

			// 过滤readOpts中设置的前缀
			if readOpts.Prefix && !strings.HasPrefix(kk, key) {
				continue
			}

			if readOpts.Suffix && !strings.HasSuffix(kk, key) {
				continue
			}

			keys = append(keys, kk)
		}
	} else {
		keys = []string{key}
	}

	var results []*Record

	for _, k := range keys {
		r, err := m.get(prefix, k)
		if err != nil {
			return results, err
		}
		results = append(results, r)
	}

	return results, nil
}

func (m *memoryStore) Write(r *Record, opts ...WriteOption) error {
	writeOpts := WriteOptions{}
	for _, o := range opts {
		o(&writeOpts)
	}

	prefix := m.prefix(writeOpts.Database, writeOpts.Table)

	if len(opts) > 0 {
		newRecord := Record{}
		newRecord.Key = r.Key
		newRecord.Value = make([]byte, len(r.Value))
		newRecord.Metadata = make(map[string]interface{})
		copy(newRecord.Value, r.Value)
		newRecord.Expiry = r.Expiry

		if !writeOpts.Expiry.IsZero() {
			newRecord.Expiry = time.Until(writeOpts.Expiry)
		}
		if writeOpts.TTL != 0 {
			newRecord.Expiry = writeOpts.TTL
		}

		for k, v := range r.Metadata {
			newRecord.Metadata[k] = v
		}

		m.set(prefix, &newRecord)
		return nil
	}

	m.set(prefix, r)

	return nil
}

func (m *memoryStore) Delete(key string, opts ...DeleteOption) error {
	deleteOptions := DeleteOptions{}
	for _, o := range opts {
		o(&deleteOptions)
	}

	prefix := m.prefix(deleteOptions.Database, deleteOptions.Table)
	m.delete(prefix, key)
	return nil
}

func (m *memoryStore) Options() Options {
	return m.options
}

func (m *memoryStore) List(opts ...ListOption) ([]string, error) {
	listOptions := ListOptions{}

	for _, o := range opts {
		o(&listOptions)
	}

	prefix := m.prefix(listOptions.Database, listOptions.Table)
	// 获取对应dabatase/table下的所有key
	keys := m.list(prefix, listOptions.Limit, listOptions.Offset)

	if len(listOptions.Prefix) > 0 {
		var prefixKeys []string
		for _, k := range keys {
			if strings.HasPrefix(k, listOptions.Prefix) {
				prefixKeys = append(prefixKeys, k)
			}
		}
		keys = prefixKeys
	}

	if len(listOptions.Suffix) > 0 {
		var suffixKeys []string
		for _, k := range keys {
			if strings.HasSuffix(k, listOptions.Suffix) {
				suffixKeys = append(suffixKeys, k)
			}
		}
		keys = suffixKeys
	}

	return keys, nil
}
