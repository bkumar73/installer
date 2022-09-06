package manifests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/assisted-service/pkg/staticnetworkconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	k8syaml "sigs.k8s.io/yaml"
)

var (
	nmStateConfigFilename = filepath.Join(clusterManifestDir, "nmstateconfig.yaml")
)

// NMStateConfig generates the nmstateconfig.yaml file.
type NMStateConfig struct {
	File                *asset.File
	StaticNetworkConfig []*models.HostStaticNetworkConfig
	Config              []*aiv1beta1.NMStateConfig
}

type nmStateConfig struct {
	Interfaces []struct {
		IPV4 struct {
			Address []struct {
				IP string `yaml:"ip,omitempty"`
			} `yaml:"address,omitempty"`
		} `yaml:"ipv4,omitempty"`
		IPV6 struct {
			Address []struct {
				IP string `yaml:"ip,omitempty"`
			} `yaml:"address,omitempty"`
		} `yaml:"ipv6,omitempty"`
	} `yaml:"interfaces,omitempty"`
}

var _ asset.WritableAsset = (*NMStateConfig)(nil)

// Name returns a human friendly name for the asset.
func (*NMStateConfig) Name() string {
	return "NMState Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*NMStateConfig) Dependencies() []asset.Asset {
	return []asset.Asset{
		&agentconfig.AgentConfig{},
	}
}

// Generate generates the NMStateConfig manifest.
func (n *NMStateConfig) Generate(dependencies asset.Parents) error {

	agentConfig := &agentconfig.AgentConfig{}
	dependencies.Get(agentConfig)

	staticNetworkConfig := []*models.HostStaticNetworkConfig{}
	nmStateConfigs := []*aiv1beta1.NMStateConfig{}
	var data string
	var isNetworkConfigAvailable bool

	if agentConfig.Config != nil {
		if len(agentConfig.Config.Hosts) == 0 {
			return nil
		}
		for i, host := range agentConfig.Config.Hosts {
			if host.NetworkConfig.Raw != nil {
				isNetworkConfigAvailable = true

				nmStateConfig := aiv1beta1.NMStateConfig{
					TypeMeta: metav1.TypeMeta{
						Kind:       "NMStateConfig",
						APIVersion: "agent-install.openshift.io/v1beta1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprintf(getNMStateConfigName(agentConfig)+"-%d", i),
						Namespace: getNMStateConfigNamespace(agentConfig),
						Labels:    getNMStateConfigLabelsFromAgentConfig(agentConfig),
					},
					Spec: aiv1beta1.NMStateConfigSpec{
						NetConfig: aiv1beta1.NetConfig{
							Raw: []byte(host.NetworkConfig.Raw),
						},
					},
				}
				for _, hostInterface := range host.Interfaces {
					intrfc := aiv1beta1.Interface{
						Name:       hostInterface.Name,
						MacAddress: hostInterface.MacAddress,
					}
					nmStateConfig.Spec.Interfaces = append(nmStateConfig.Spec.Interfaces, &intrfc)

				}
				nmStateConfigs = append(nmStateConfigs, &nmStateConfig)

				staticNetworkConfig = append(staticNetworkConfig, &models.HostStaticNetworkConfig{
					MacInterfaceMap: buildMacInterfaceMap(nmStateConfig),
					NetworkYaml:     string(nmStateConfig.Spec.NetConfig.Raw),
				})

				// Marshal the nmStateConfig one at a time
				// and add a yaml seperator with new line
				// so as not to marshal the nmStateConfigs
				// as a yaml list in the generated nmstateconfig.yaml
				nmStateConfigData, err := k8syaml.Marshal(nmStateConfig)

				if err != nil {
					return errors.Wrap(err, "failed to marshal agent installer NMStateConfig")
				}
				data = fmt.Sprint(data, fmt.Sprint(string(nmStateConfigData), "---\n"))
			}
		}

		if isNetworkConfigAvailable {
			n.Config = nmStateConfigs
			n.StaticNetworkConfig = staticNetworkConfig

			n.File = &asset.File{
				Filename: nmStateConfigFilename,
				Data:     []byte(data),
			}
		}
	}
	return n.finish()
}

// Files returns the files generated by the asset.
func (n *NMStateConfig) Files() []*asset.File {
	if n.File != nil {
		return []*asset.File{n.File}
	}
	return []*asset.File{}
}

// Load returns the NMStateConfig asset from the disk.
func (n *NMStateConfig) Load(f asset.FileFetcher) (bool, error) {

	file, err := f.FetchByName(nmStateConfigFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrapf(err, "failed to load file %s", nmStateConfigFilename)
	}

	// Split up the file into multiple YAMLs if it contains NMStateConfig for more than one node
	var decoder nmStateConfigYamlDecoder
	yamlList, err := getMultipleYamls(file.Data, &decoder)
	if err != nil {
		return false, errors.Wrapf(err, "could not decode YAML for %s", nmStateConfigFilename)
	}

	var staticNetworkConfig []*models.HostStaticNetworkConfig
	var nmStateConfigList []*aiv1beta1.NMStateConfig

	for i := range yamlList {
		nmStateConfig := yamlList[i].(*aiv1beta1.NMStateConfig)
		staticNetworkConfig = append(staticNetworkConfig, &models.HostStaticNetworkConfig{
			MacInterfaceMap: buildMacInterfaceMap(*nmStateConfig),
			NetworkYaml:     string(nmStateConfig.Spec.NetConfig.Raw),
		})
		nmStateConfigList = append(nmStateConfigList, nmStateConfig)
	}

	log := logrus.New()
	log.Level = logrus.WarnLevel
	staticNetworkConfigGenerator := staticnetworkconfig.New(log.WithField("pkg", "manifests"), staticnetworkconfig.Config{MaxConcurrentGenerations: 2})

	// Validate the network config using nmstatectl
	if err = staticNetworkConfigGenerator.ValidateStaticConfigParams(context.Background(), staticNetworkConfig); err != nil {
		return false, errors.Wrapf(err, "staticNetwork configuration is not valid")
	}

	n.File, n.StaticNetworkConfig, n.Config = file, staticNetworkConfig, nmStateConfigList
	if err = n.finish(); err != nil {
		return false, err
	}
	return true, nil
}

func (n *NMStateConfig) finish() error {

	if n.Config == nil {
		return errors.New("missing configuration or manifest file")
	}

	if err := n.validateNMStateConfig().ToAggregate(); err != nil {
		return errors.Wrapf(err, "invalid NMStateConfig configuration")
	}
	return nil
}

func (n *NMStateConfig) validateNMStateConfig() field.ErrorList {
	allErrs := field.ErrorList{}

	if err := n.validateNMStateLabels(); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func (n *NMStateConfig) validateNMStateLabels() field.ErrorList {

	var allErrs field.ErrorList

	fieldPath := field.NewPath("ObjectMeta", "Labels")

	for _, nmStateConfig := range n.Config {
		if len(nmStateConfig.ObjectMeta.Labels) == 0 {
			allErrs = append(allErrs, field.Required(fieldPath, fmt.Sprintf("%s does not have any label set", nmStateConfig.Name)))
		}
	}

	return allErrs
}

func getFirstIP(nmStateConfig *nmStateConfig) string {
	for _, intf := range nmStateConfig.Interfaces {
		for _, addr4 := range intf.IPV4.Address {
			if addr4.IP != "" {
				return addr4.IP
			}
		}
		for _, addr6 := range intf.IPV6.Address {
			if addr6.IP != "" {
				return addr6.IP
			}
		}
	}
	return ""
}

// GetNodeZeroIP retrieves the first IP from the user provided NMStateConfigs to set as the node0 IP
func GetNodeZeroIP(nmStateConfigs []*aiv1beta1.NMStateConfig) (string, error) {
	for i := range nmStateConfigs {
		var nmStateConfig nmStateConfig
		err := yaml.Unmarshal(nmStateConfigs[i].Spec.NetConfig.Raw, &nmStateConfig)
		if err != nil {
			return "", fmt.Errorf("error unmarshalling NMStateConfig: %v", err)
		}
		if nodeZeroIP := getFirstIP(&nmStateConfig); nodeZeroIP != "" {
			if net.ParseIP(nodeZeroIP) == nil {
				return "", fmt.Errorf("could not parse static IP: %s", nodeZeroIP)
			}

			return nodeZeroIP, nil
		}

	}

	return "", fmt.Errorf("invalid NMStateConfig yaml, no interface IPs set")
}

// GetNMIgnitionFiles returns the list of NetworkManager configuration files
func GetNMIgnitionFiles(staticNetworkConfig []*models.HostStaticNetworkConfig) ([]staticnetworkconfig.StaticNetworkConfigData, error) {
	log := logrus.New()
	staticNetworkConfigGenerator := staticnetworkconfig.New(log.WithField("pkg", "manifests"), staticnetworkconfig.Config{MaxConcurrentGenerations: 2})

	networkConfigStr, err := staticNetworkConfigGenerator.FormatStaticNetworkConfigForDB(staticNetworkConfig)
	if err != nil {
		err = fmt.Errorf("error marshalling StaticNetwork configuration: %w", err)
		return nil, err
	}

	filesList, err := staticNetworkConfigGenerator.GenerateStaticNetworkConfigData(context.Background(), networkConfigStr)
	if err != nil {
		err = fmt.Errorf("failed to create StaticNetwork config data: %w", err)
		return nil, err
	}

	return filesList, err
}

type nmStateConfigYamlDecoder int

type decodeFormat interface {
	NewDecodedYaml(decoder *yaml.YAMLToJSONDecoder) (interface{}, error)
}

func (d *nmStateConfigYamlDecoder) NewDecodedYaml(yamlDecoder *yaml.YAMLToJSONDecoder) (interface{}, error) {
	decodedData := new(aiv1beta1.NMStateConfig)
	err := yamlDecoder.Decode(&decodedData)

	return decodedData, err
}

// Read a YAML file containing multiple YAML definitions of the same format
// Each specific format must be of type DecodeFormat
func getMultipleYamls(contents []byte, decoder decodeFormat) ([]interface{}, error) {

	r := bytes.NewReader(contents)
	dec := yaml.NewYAMLToJSONDecoder(r)

	var outputList []interface{}
	for {
		decodedData, err := decoder.NewDecodedYaml(dec)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, errors.Wrapf(err, "Error reading multiple YAMLs")
		}

		outputList = append(outputList, decodedData)
	}

	return outputList, nil
}

func buildMacInterfaceMap(nmStateConfig aiv1beta1.NMStateConfig) models.MacInterfaceMap {

	// TODO - this eventually will move to another asset so the interface definition can be shared with Butane
	macInterfaceMap := make(models.MacInterfaceMap, 0, len(nmStateConfig.Spec.Interfaces))
	for _, cfg := range nmStateConfig.Spec.Interfaces {
		logrus.Debug("adding MAC interface map to host static network config - Name: ", cfg.Name, " MacAddress:", cfg.MacAddress)
		macInterfaceMap = append(macInterfaceMap, &models.MacInterfaceMapItems0{
			MacAddress:     cfg.MacAddress,
			LogicalNicName: cfg.Name,
		})
	}
	return macInterfaceMap
}
