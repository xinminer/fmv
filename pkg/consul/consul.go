package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
)

type DiscoveryConfig struct {
	ID      string
	Name    string
	Tags    []string
	Port    int
	Address string
}

func RegisterService(addr string, dis DiscoveryConfig) error {
	config := consulapi.DefaultConfig()
	config.Address = addr
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Printf("create consul client : %v\n", err.Error())
		return err
	}
	registration := &consulapi.AgentServiceRegistration{
		ID:      dis.ID,
		Name:    dis.Name,
		Port:    dis.Port,
		Tags:    dis.Tags,
		Address: dis.Address,
	}

	check := &consulapi.AgentServiceCheck{}
	check.TCP = fmt.Sprintf("%s:%d", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "60s"
	registration.Check = check

	if err := client.Agent().ServiceRegister(registration); err != nil {
		return err
	}

	return nil
}
