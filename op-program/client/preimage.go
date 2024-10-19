package client

import (
	"database/sql"
	"encoding/binary"
	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/mattn/go-sqlite3"
)

type MemoryPreimageOracle struct {
	data map[common.Hash][]byte
}

func NewMemoryPreimageOracle(path string) (*MemoryPreimageOracle, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT key, value FROM kv_store")
	if err != nil {
		return nil, err
	}
	data := map[common.Hash][]byte{}
	count := 0
	for rows.Next() {
		count++
		if count%10000 == 0 {
			println("current loaded", count)
		}
		var key []byte
		var value []byte
		if err = rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		data[common.BytesToHash(key)] = value
	}
	return &MemoryPreimageOracle{
		data: data,
	}, nil
}

func (p *MemoryPreimageOracle) Hint(v preimage.Hint) {

}

func (p *MemoryPreimageOracle) Get(k preimage.Key) []byte {
	switch k.PreimageKey() {
	case L1HeadLocalIndex.PreimageKey():
		return common.FromHex("0x93ba31bf89e54237af6e6564e69d328b2b5202adf643de4cb097431f74f4a6c1")
	case L2OutputRootLocalIndex.PreimageKey():
		return common.FromHex("0x056a42a72c62b0e80658cfc6ff0e87419cb38771d16a69c9399a58a28046e281")
	case L2ClaimLocalIndex.PreimageKey():
		return common.FromHex("0x0615473db962c6ccf828d01f5fe3f12b167588047435b4ff433660f5aa64875b")
	case L2ClaimBlockNumberLocalIndex.PreimageKey():
		return binary.BigEndian.AppendUint64(nil, 15378356)
	case L2ChainIDLocalIndex.PreimageKey():
		opSepoliaChainId := uint64(11155420)
		return binary.BigEndian.AppendUint64(nil, opSepoliaChainId)
	}
	key := k.PreimageKey()
	return p.data[common.BytesToHash(key[:])]
}
