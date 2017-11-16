package elementalconductor

import "encoding/xml"

// NodeProduct is the product that is running inside a node.
type NodeProduct string

const (
	// ProductConductorFile is condutor file product.
	ProductConductorFile = NodeProduct("Conductor File")

	// ProductServer is the server product.
	ProductServer = NodeProduct("Server")
)

type nodeList struct {
	XMLName xml.Name `xml:"node_list"`
	Nodes   []Node   `xml:"node"`
}

// Node is a server running one of Elemental products in one of its platforms.
type Node struct {
	Href            string      `xml:"href,attr"`
	Name            string      `xml:"name"`
	HostName        string      `xml:"hostname"`
	IPAddress       string      `xml:"ip_addr"`
	PublicIPAddress string      `xml:"public_ip_addr,omitempty"`
	Eth0Mac         string      `xml:"eth0_mac"`
	Status          string      `xml:"status"`
	Product         NodeProduct `xml:"product"`
	Version         string      `xml:"version"`
	Platform        string      `xml:"platform"`
	Packages        []string    `xml:"packages>package"`
	Licenses        []string    `xml:"licenses>license"`
	CreatedAt       DateTime    `xml:"created_at"`
	RunningCount    int         `xml:"running_count,omitempty"`
}

// GetNodes returns the list of nodes currently available in the Elemental
// setup.
func (c *Client) GetNodes() ([]Node, error) {
	var result nodeList
	err := c.do("GET", "/nodes", nil, &result)
	if err != nil {
		return nil, err
	}
	return result.Nodes, nil
}
