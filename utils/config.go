package utils

import (
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/yaml.v3"
)

var (
	DefaultGasLimit uint64 = 1600000

	StakeHubAddress = common.HexToAddress("0x0000000000000000000000000000000000002002")
)

type Config struct {
	BscRpcUrl     string        `yaml:"BscRpcUrl"`
	BlsDataDir    string        `yaml:"BlsDataDir"`
	ValidatorInfo ValidatorInfo `yaml:"ValidatorInfo"`
}

type ValidatorInfo struct {
	Delegation       string      `yaml:"Delegation"`
	ConsensusAddress string      `yaml:"ConsensusAddress"`
	Description      Description `yaml:"Description"`
	Commission       Commission  `yaml:"Commission"`
}

type Description struct {
	Moniker  string `yaml:"moniker"`
	Identity string `yaml:"identity"`
	Website  string `yaml:"website"`
	Details  string `yaml:"details"`
}

type Commission struct {
	Rate          uint64 `yaml:"rate"`
	MaxRate       uint64 `yaml:"maxRate"`
	MaxChangeRate uint64 `yaml:"maxChangeRate"`
}

func NewConfig() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fPath := filepath.Join(wd, "config/config.yml")
	yamlFile, err := os.ReadFile(fPath)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
