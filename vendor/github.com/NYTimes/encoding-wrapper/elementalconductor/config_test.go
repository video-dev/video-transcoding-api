package elementalconductor

import (
	"encoding/xml"
	"net/http"

	"gopkg.in/check.v1"
)

func (s *S) TestGetCloudConfig(c *check.C) {
	data := `<?xml version="1.0" encoding="UTF-8"?>
<cloud_config>
  <authorized_node_count>500</authorized_node_count>
  <max_cluster_size>30</max_cluster_size>
  <min_cluster_size>4</min_cluster_size>
  <worker_variant>production_server_cloud</worker_variant>
</cloud_config>`
	server, requests := s.startServer(http.StatusOK, data)
	defer server.Close()
	client := NewClient(server.URL, "myuser", "secret-key", 45, "aws-access-key", "aws-secret-key", "destination")
	config, err := client.GetCloudConfig()
	c.Assert(err, check.IsNil)
	c.Assert(*config, check.DeepEquals, CloudConfig{
		XMLName:             xml.Name{Local: "cloud_config"},
		AuthorizedNodeCount: 500,
		MaxNodes:            30,
		MinNodes:            4,
		WorkerVariant:       "production_server_cloud",
	})
	fakeReq := <-requests
	c.Assert(fakeReq.req.Method, check.Equals, "GET")
	c.Assert(fakeReq.req.URL.Path, check.Equals, "/api/config/cloud")
}
