package elementalconductor

import (
	"encoding/xml"
	"time"
)

const (
	dateTimeLayout      = "2006-01-02 15:04:05 -0700"
	errorDateTimeLayout = "2006-01-02T15:04:05-07:00"
)

// DateTime is a custom struct for representing time within ElementalConductor.
// It customizes marshalling, and always store the underlying time in UTC.
type DateTime struct {
	time.Time
}

// MarshalXML implementation on DateTimeg to skip "zero" time values
func (jdt DateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !jdt.IsZero() {
		e.EncodeElement(jdt.Time.Format(dateTimeLayout), start)
	}
	return nil
}

// UnmarshalXML implementation on DateTimeg to use dateTimeLayout
func (jdt *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var err error
	jdt.Time, err = unmarshalTime(d, start, dateTimeLayout)
	return err
}

// JobErrorDateTime is a custom time struct to be used on Media items
type JobErrorDateTime struct {
	time.Time
}

// MarshalXML implementation on JobErrorDateTime to skip "zero" time values
func (jdt JobErrorDateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !jdt.IsZero() {
		e.EncodeElement(jdt.Time.Format(errorDateTimeLayout), start)
	}
	return nil
}

// UnmarshalXML implementation on JobErrorDateTime to use errorDateTimeLayout
func (jdt *JobErrorDateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var err error
	jdt.Time, err = unmarshalTime(d, start, errorDateTimeLayout)
	return err
}

func unmarshalTime(d *xml.Decoder, start xml.StartElement, format string) (time.Time, error) {
	var t time.Time
	var content string
	err := d.DecodeElement(&content, &start)
	if err != nil {
		return t, err
	}
	if content == "" {
		return t, nil
	}
	if content == "0001-01-01T00:00:00Z" {
		return t, nil
	}
	t, err = time.Parse(format, content)
	if err != nil {
		return t, err
	}
	t = t.UTC()
	return t, nil
}
