package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config contains the full list of default parameters to initialize the app
type Config struct {
	Version string       `mapstructure:"version" json:"version,omitempty"`
	FabCfg  FabricConfig `mapstructure:"fabric" json:"fabric,omitempty"`
}

// FabricConfig contains all the attributes related to Hyperledger Fabric including docker image information
type FabricConfig struct {
	Images  []DockerImage `mapstructure:"images" json:"images,omitempty"`
	Network FabricNetwork `mapstructure:"network" json:"network,omitempty"`
}

// DockerImage is defined by image name and tag
type DockerImage struct {
	Image string `mapstructure:"image" json:"image,omitempty"`
	Tag   string `mapstructure:"tag" json:"tag,omitempty"`
}

// FabricNetwork contains all the necessary information to run an Hyperledger Fabric network
type FabricNetwork struct {
	Profile       NetworkProfile `mapstructure:"profile" json:"profile,omitempty"`
	Organizations []Org          `mapstructure:"organizations" json:"organizations,omitempty"`
	Orderers      []Orderer      `mapstructure:"orderers" json:"orderers,omitempty"`
	TLS           bool           `mapstructure:"tls" json:"tls,omitempty"`
}

// NetworkProfile maps network profile information contained in the default configtx.yaml file
type NetworkProfile struct {
	NetworkConfig string `mapstructure:"network" json:"network,omitempty"`
	ChannelConfig string `mapstructure:"channel" json:"channel,omitempty"`
}

// Org contains information related to a specific organization in the network
type Org struct {
	MSP        string      `mapstructure:"msp" json:"msp,omitempty"`
	Chaincodes []Chaincode `mapstructure:"chaincodes" json:"chaincodes,omitempty"`
	Channels   []string    `mapstructure:"channels" json:"channels,omitempty"`
}

// Chaincode is defined by name and version (together they must be unique)
type Chaincode struct {
	Name    string `mapstructure:"name" json:"name,omitempty"`
	Version string `mapstructure:"version" json:"version,omitempty"`
}

// Orderer defines an orderer organization in the network
type Orderer struct {
	MSP     string `mapstructure:"msp" json:"msp,omitempty"`
	Address string `mapstructure:"address" json:"address,omitempty"`
	Channel string `mapstructure:"channel" json:"channel,omitempty"`
}

// GetConfig returns a previously initialized configuration
func GetConfig() Config {
	config := Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalln(err)
	}

	return config
}
