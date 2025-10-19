package village

import "strings"

type Village struct {
	Alive []VillageElement
}

func (V *Village) UpdateAll() {
	for i := range V.Alive {
		V.Alive[i].Update()
	}
}

func (V *Village) ShowAllInfo() string {
	var infos []string
	for _, person := range V.Alive {
		infos = append(infos, person.FlushInfo())
	}

	V.DeleteDied()

	return strings.Join(infos, "\n\n")
}

func (V *Village) AddElement(Element VillageElement) {
	V.Alive = append(V.Alive, Element)
}

func (V *Village) DeleteDied() {
	alive := make([]VillageElement, 0, len(V.Alive))

	for _, element := range V.Alive {
		if element.CheckAlive() {
			alive = append(alive, element)
		}
	}

	V.Alive = alive
}
