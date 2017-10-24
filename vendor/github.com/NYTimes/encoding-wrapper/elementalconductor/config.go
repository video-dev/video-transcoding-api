package elementalconductor

import "encoding/xml"

// CloudConfig contains configuration for Elemental Cloud, including Autoscaler
// Settings.
type CloudConfig struct {
	XMLName             xml.Name `xml:"cloud_config"`
	AuthorizedNodeCount int      `xml:"authorized_node_count"`
	MaxNodes            int      `xml:"max_cluster_size"`
	MinNodes            int      `xml:"min_cluster_size"`
	WorkerVariant       string   `xml:"worker_variant"`
}

// GetCloudConfig returns the current Elemental Cloud configuration. It
// includes Autoscaler Settings.
func (c *Client) GetCloudConfig() (*CloudConfig, error) {
	var config CloudConfig
	err := c.do("GET", "/config/cloud", nil, &config)
	return &config, err
}
