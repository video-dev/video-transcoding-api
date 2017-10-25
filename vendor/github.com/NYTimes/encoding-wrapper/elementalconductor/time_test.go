package elementalconductor

import (
	"encoding/xml"
	"fmt"
	"time"

	"gopkg.in/check.v1"
)

func (s *S) TestDateTimeMarshalXML(c *check.C) {
	var tests = []struct {
		input    time.Time
		expected string
	}{
		{
			time.Time{},
			"<item></item>",
		},
		{
			time.Date(2016, 12, 7, 21, 28, 43, 0, time.UTC),
			"<item><date>2016-12-07 21:28:43 +0000</date></item>",
		},
	}
	for _, test := range tests {
		var data struct {
			XMLName xml.Name `xml:"item"`
			Date    DateTime `xml:"date,omitempty"`
		}
		data.Date.Time = test.input
		b, err := xml.Marshal(data)
		c.Check(err, check.IsNil)
		c.Check(string(b), check.Equals, test.expected)
	}
}

func (s *S) TestDateTimeUnmarshalXML(c *check.C) {
	var tests = []struct {
		input    string
		expected time.Time
	}{
		{
			"2016-02-01 11:59:20 -0800",
			time.Date(2016, time.February, 1, 19, 59, 20, 0, time.UTC),
		},
		{
			"2016-02-01 00:25:00 +0300",
			time.Date(2016, time.January, 31, 21, 25, 0, 0, time.UTC),
		},
		{
			"",
			time.Time{},
		},
		{
			"0001-01-01T00:00:00Z",
			time.Time{},
		},
	}
	for _, test := range tests {
		var output struct {
			XMLName xml.Name `xml:"item"`
			Date    DateTime `xml:"date"`
		}
		input := fmt.Sprintf("<item><date>%s</date></item>", test.input)
		err := xml.Unmarshal([]byte(input), &output)
		c.Check(err, check.IsNil)
		c.Check(output.Date.Time, check.DeepEquals, test.expected)
	}
}

func (s *S) TestDateTimeUnmarshalXMLInvalidFormat(c *check.C) {
	var output struct {
		XMLName xml.Name `xml:"item"`
		Date    DateTime `xml:"date"`
	}
	input := "<item><date>2016-13-01 15:03:02 -0300</date></item>"
	err := xml.Unmarshal([]byte(input), &output)
	c.Assert(err, check.NotNil)
}

func (s *S) TestJobErrorDateTimeMarshalXML(c *check.C) {
	var tests = []struct {
		input    time.Time
		expected string
	}{
		{
			time.Time{},
			"<item></item>",
		},
		{
			time.Date(2016, 12, 7, 21, 28, 43, 0, time.UTC),
			"<item><date>2016-12-07T21:28:43+00:00</date></item>",
		},
	}
	for _, test := range tests {
		var data struct {
			XMLName xml.Name         `xml:"item"`
			Date    JobErrorDateTime `xml:"date,omitempty"`
		}
		data.Date.Time = test.input
		b, err := xml.Marshal(data)
		c.Check(err, check.IsNil)
		c.Check(string(b), check.Equals, test.expected)
	}
}

func (s *S) TestJobErrorDateTimeUnmarshalXML(c *check.C) {
	var tests = []struct {
		input    string
		expected time.Time
	}{
		{
			"2016-02-01T11:59:20-08:00",
			time.Date(2016, time.February, 1, 19, 59, 20, 0, time.UTC),
		},
		{
			"2016-02-01T00:25:00+03:00",
			time.Date(2016, time.January, 31, 21, 25, 0, 0, time.UTC),
		},
		{
			"",
			time.Time{},
		},
		{
			"0001-01-01T00:00:00Z",
			time.Time{},
		},
	}
	for _, test := range tests {
		var output struct {
			XMLName xml.Name         `xml:"item"`
			Date    JobErrorDateTime `xml:"date"`
		}
		input := fmt.Sprintf("<item><date>%s</date></item>", test.input)
		err := xml.Unmarshal([]byte(input), &output)
		c.Check(err, check.IsNil)
		c.Check(output.Date.Time, check.DeepEquals, test.expected)
	}
}

func (s *S) TestJobErrorDateTimeUnmarshalXMLInvalidFormat(c *check.C) {
	var output struct {
		XMLName xml.Name         `xml:"item"`
		Date    JobErrorDateTime `xml:"date"`
	}
	input := "<item><date>2016-13-01 15:03:02 -0300</date></item>"
	err := xml.Unmarshal([]byte(input), &output)
	c.Assert(err, check.NotNil)
}
