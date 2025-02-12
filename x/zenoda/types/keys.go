package types

const (
	// ModuleName defines the module name
	ModuleName = "zenoda"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_zenoda"
)

var (
	ParamsKey = []byte("p_zenoda")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
