package village

type VillageElement interface {
	Update()
	FlushInfo() string
	CheckAlive() bool
}
