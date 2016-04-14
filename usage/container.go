package usage

import (
	"fmt"
)


type UsageContainer struct {
	Usages map[string]*Usage
}



func NewUsageContainer() *UsageContainer {
	container := &UsageContainer{
		Usages: make(map[string]*Usage),
	}

	return container
}


func (c *UsageContainer) NewUsage(name string, port int) *Usage {
	var usage *Usage

	key := fmt.Sprintf("%s-%s", name, port)
	if value, exists := c.Usages[key]; exists {
		usage = value
	} else {
		usage = &Usage{
			Protocol: name,
			Port: port,
		}

		c.Usages[key] = usage
		go usage.Launch()
	}

	return usage
}