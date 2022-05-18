package extend

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"strings"
)

type DubboGatewayOps struct {
	IsDirect         bool
	Protocol         string
	InterfaceName    string
	RegistryProtocol string
}

type DubboGatewayConfig struct {
	isDirect              bool
	protocol              string
	interfaceName         string
	registryProtocol      string
	rootConfigBuilder     *config.RootConfigBuilder
	consumerConfigBuilder *config.ConsumerConfigBuilder
}

// RegistryConfig is the configuration of the registry center
type RegistryConfig struct {
	Protocol  string `validate:"required" yaml:"protocol"  json:"protocol,omitempty" property:"protocol"`
	Timeout   string `default:"5s" validate:"required" yaml:"timeout" json:"timeout,omitempty" property:"timeout"` // unit: second
	Group     string `yaml:"group" json:"group,omitempty" property:"group"`
	Namespace string `yaml:"namespace" json:"namespace,omitempty" property:"namespace"`
	TTL       string `default:"10s" yaml:"ttl" json:"ttl,omitempty" property:"ttl"` // unit: minute
	// for registry
	Address    string `validate:"required" yaml:"address" json:"address,omitempty" property:"address"`
	Username   string `yaml:"username" json:"username,omitempty" property:"username"`
	Password   string `yaml:"password" json:"password,omitempty"  property:"password"`
	Simplified bool   `yaml:"simplified" json:"simplified,omitempty"  property:"simplified"`
	// Always use this registry first if set to true, useful when subscribe to multiple registriesConfig
	Preferred bool `yaml:"preferred" json:"preferred,omitempty" property:"preferred"`
	// The region where the registry belongs, usually used to isolate traffics
	Zone string `yaml:"zone" json:"zone,omitempty" property:"zone"`
	// Affects traffic distribution among registriesConfig,
	// useful when subscribe to multiple registriesConfig Take effect only when no preferred registry is specified.
	Weight       int64             `yaml:"weight" json:"weight,omitempty" property:"weight"`
	Params       map[string]string `yaml:"params" json:"params,omitempty" property:"params"`
	RegistryType string            `yaml:"registry-type"`
}

var appNameEndpointMap = make(map[string]string)

func NewDubboGatewayConfig(ops *DubboGatewayOps) *DubboGatewayConfig {
	return &DubboGatewayConfig{
		isDirect:              ops.IsDirect,
		protocol:              ops.Protocol,
		interfaceName:         ops.InterfaceName,
		registryProtocol:      ops.RegistryProtocol,
		rootConfigBuilder:     config.NewRootConfigBuilder(),
		consumerConfigBuilder: config.NewConsumerConfigBuilder(),
	}
}

func (gwConfig *DubboGatewayConfig) SetConsumerService(service common.RPCService) {
	config.SetConsumerService(service)
}

func (gwConfig *DubboGatewayConfig) AddRegistry(registryProtocol string, registryAddrs ...string) {
	if gwConfig.isDirect {
		return
	}

	registryAddress := convertRegistryAddress(registryAddrs)
	registryConfig := config.NewRegistryConfig([]config.RegistryConfigOpt{
		config.WithRegistryProtocol(registryProtocol),
		config.WithRegistryAddress(registryAddress),
	}...)

	/*opts := []config.RegistryConfigOpt{
		config.WithRegistryProtocol(registryProtocol),
		config.WithRegistryAddress(registryAddress),
	}
	registryConfig = config.NewRegistryConfig(opts...)*/

	gwConfig.consumerConfigBuilder.SetRegistryIDs(registryProtocol)
	gwConfig.rootConfigBuilder.AddRegistry(registryProtocol, registryConfig)
}

func convertRegistryAddress(registryAddrs []string) string {
	var registryAddress string
	if registryAddrs == nil {
		registryAddress = "127.0.0.1:2181"
	} else {
		registryAddress = strings.Join(registryAddrs, ",")
	}
	return registryAddress
}

func (gwConfig *DubboGatewayConfig) AddReferenceEndpoint(appName string, endpoint string) {
	appNameEndpointMap[appName] = endpoint
}

func (gwConfig *DubboGatewayConfig) AddReference(appName, referenceKey, interfaceName string) error {
	if gwConfig.isDirect {
		endpoint := appNameEndpointMap[appName]
		refConf := config.ReferenceConfig{
			Protocol:      gwConfig.protocol,
			URL:           endpoint,
			InterfaceName: interfaceName,
		}
		gwConfig.consumerConfigBuilder.AddReference(referenceKey, &refConf)
	} else {
		gwConfig.consumerConfigBuilder.
			AddReference(referenceKey, config.NewReferenceConfigBuilder().
				SetProtocol(gwConfig.protocol).
				SetInterface(interfaceName).
				Build())
	}
	return nil
}

func (gwConfig *DubboGatewayConfig) Load() error {
	// Load只会执行一次
	rootConfig := gwConfig.rootConfigBuilder.
		SetConsumer(gwConfig.consumerConfigBuilder.Build()).
		Build()
	if err := config.Load(config.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}
	return nil
}
