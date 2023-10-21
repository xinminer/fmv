package consul

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/util/grand"
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

func Discovery(serviceName string, address string, tag string) (string, error) {
	config := consulapi.DefaultConfig()
	config.Address = address
	client, err := consulapi.NewClient(config)
	if err != nil {
		return "", err
	}
	services, _, err := client.Health().Service(serviceName, tag, false, nil)
	if err != nil {
		return "", err
	}
	ses := len(services)
	if ses == 0 {
		services, _, err := client.Health().Service(serviceName, "", false, nil)
		if err != nil {
			return "", err
		}
		ses := len(services)
		ses = len(services)
		if ses == 0 {
			return "", errors.New("not found service")
		}
	}
	rand := grand.N(0, ses-1)
	se := services[rand]

	result := fmt.Sprintf("%s:%d", se.Service.Address, se.Service.Port)

	return result, nil
}
