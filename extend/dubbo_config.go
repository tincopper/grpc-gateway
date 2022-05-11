package extend

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
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

func (gwConfig *DubboGatewayConfig) AddRegistry(registryProtocol string) error {
	if !gwConfig.isDirect {
		return nil
	}
	gwConfig.consumerConfigBuilder.SetRegistryIDs(registryProtocol)
	gwConfig.rootConfigBuilder.AddRegistry(registryProtocol, config.NewRegistryConfigWithProtocolDefaultPort(registryProtocol))
	return nil
}

func (gwConfig *DubboGatewayConfig) AddReferenceEndpoint(appName string, endpoint string) {
	appNameEndpointMap[appName] = endpoint
}

func (gwConfig *DubboGatewayConfig) AddReference(appName, referenceKey string) error {
	if gwConfig.isDirect {
		endpoint := appNameEndpointMap[appName]
		refConf := config.ReferenceConfig{
			Protocol: gwConfig.protocol,
			URL:      endpoint,
		}
		gwConfig.consumerConfigBuilder.AddReference(referenceKey, &refConf)
	} else {
		gwConfig.consumerConfigBuilder.
			AddReference(referenceKey, config.NewReferenceConfigBuilder().
				SetProtocol(gwConfig.protocol).
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
