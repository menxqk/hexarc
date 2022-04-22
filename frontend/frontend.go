package frontend

import "github.com/menxqk/hexarc/core"

type FrontEnd interface {
	Start(kv *core.KeyValueStore) error
}

type zeroFrontEnd struct{}

func (z zeroFrontEnd) Start(kv *core.KeyValueStore) error {
	return nil
}
