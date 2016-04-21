package usage

import (
	"fmt"
	"time"
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


func (c *UsageContainer) NewUsage(name string, port int, inactive time.Duration) *Usage {
	var usage *Usage

	usageID := fmt.Sprintf("%s-%s", name, port)
	if value, exists := c.Usages[usageID]; exists {
		usage = value
	} else {
		usage = &Usage{
			ID: usageID,
			Protocol: name,
			Port: port,
			InactiveTimeout: inactive,
		}

		c.Usages[usageID] = usage
		go usage.Launch()
	}

	return usage
}