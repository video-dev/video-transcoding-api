package elementalconductor

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"

	"github.com/video-dev/go-elementalconductor"
)

type nodeList struct {
	XMLName xml.Name                  `xml:"node_list"`
	Nodes   []elementalconductor.Node `xml:"node"`
}

type ElementalServer struct {
	*httptest.Server
	nodes  *nodeList
	config *elementalconductor.CloudConfig
}

func NewElementalServer(config *elementalconductor.CloudConfig, nodes []elementalconductor.Node) *ElementalServer {
	s := ElementalServer{
		nodes:  &nodeList{XMLName: xml.Name{Local: "node_list"}, Nodes: nodes},
		config: config,
	}
	s.Server = httptest.NewServer(&s)
	return &s
}

func (s *ElementalServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/nodes":
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(s.nodes)
	case "/api/config/cloud":
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(s.config)
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (s *ElementalServer) SetCloudConfig(config *elementalconductor.CloudConfig) {
	s.config = config
}

func (s *ElementalServer) SetNodes(nodes []elementalconductor.Node) {
	s.nodes.Nodes = nodes
}
