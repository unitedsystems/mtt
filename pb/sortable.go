package pb

func (mp *MessagePack) Len() int {
	return len(mp.Messages)
}

func (mp *MessagePack) Less(i, j int) bool {
	return mp.Messages[i].Timestamp < mp.Messages[j].Timestamp
}

func (mp *MessagePack) Swap(i, j int) {
	mp.Messages[i], mp.Messages[j] = mp.Messages[j], mp.Messages[i]
}
