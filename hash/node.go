package hash

type NodeHash struct {
	slotCount   uint32
	slotNodeMap map[int]int
}

func NewNodeHash(nodeSlotMap map[int][]int) *NodeHash {
	h := new(NodeHash)
	h.Init(nodeSlotMap)
	return h
}

func (h *NodeHash) Init(nodeSlotMap map[int][]int) {
	h.slotNodeMap = make(map[int]int)
	for nodeId, slotRange := range nodeSlotMap {
		for slotId := slotRange[0]; slotId <= slotRange[1]; slotId++ {
			h.slotNodeMap[slotId] = nodeId
		}
	}
	h.slotCount = uint32(len(h.slotNodeMap))
}

func (h *NodeHash) Hash(key string) int {
	slotId := int(HashCrc32String(key)%h.slotCount) + 1
	return h.slotNodeMap[slotId]
}
