package peer

type Vote struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Term    uint   `json:"term"`
}

type Map[K comparable, V any] struct {
	entry map[K]V
}

func (m *Map[K, V]) Get(key K) V {

	return m.entry[key]
}

func (m *Map[K, V]) GetEntries() map[K]V {
	return m.entry
}

func (m *Map[K, V]) Put(key K, value V) {
	m.entry[key] = value
}

func (m *Map[K, V]) Delete(key K) {
	delete(m.entry, key)
}

func NewMap[K comparable, V any]() Map[K, V] {
	return Map[K, V]{entry: map[K]V{}}
}

type PeerState string

const (
	LEADER    PeerState = "LEADER"
	CANDIDATE PeerState = "CANDIDATE"
	FOLLOWER  PeerState = "FOLLOWER"
	DEAD      PeerState = "DEAD"
)

var States = map[string]PeerState{
	"LEADER":    LEADER,
	"CANDIDATE": CANDIDATE,
	"FOLLOWER":  FOLLOWER,
	"DEAD":      DEAD,
}

type Event struct{}

var (
	PeerAdded         Event
	PeerRemoved       Event
	PeerStatusChanged Event
	PeerDead          Event
	PeerAlive         Event
	ElectionIntiated  Event
)
