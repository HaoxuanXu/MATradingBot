package channel

type ChanMap struct {
	Map map[string](chan bool)
}

func CreateMap(assets []string) *ChanMap {
	var resultMap ChanMap
	for _, asset := range assets {
		resultMap.Map[asset] = make(chan bool)
	}
	return &resultMap
}

func (chanMap *ChanMap) TriggerWorkers() {
	for _, channel := range chanMap.Map {
		channel <- true
	}
}

func (chanMap *ChanMap) CloseWorkers() {
	for _, channel := range chanMap.Map {
		channel <- false
	}
}

func (chanMap *ChanMap) CloseChannels() {
	for _, channel := range chanMap.Map {
		close(channel)
	}
}
