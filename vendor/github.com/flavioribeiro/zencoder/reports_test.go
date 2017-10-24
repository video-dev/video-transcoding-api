package zencoder

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetReportQuery(t *testing.T) {
	resp := GetReportQuery("/path", nil)
	if resp != "/path" {
		t.Fatal("Expected /path, got", resp)
	}

	var settings *ReportSettings
	resp = GetReportQuery("/path", settings)
	if resp != "/path" {
		t.Fatal("Expected /path, got", resp)
	}

	now := time.Date(2013, time.November, 22, 0, 0, 0, 0, time.UTC)

	settings = Report()

	settings = settings.ReportFrom(now)
	resp = GetReportQuery("/path", settings)
	if resp != "/path?from=2013-11-22" {
		t.Fatal("Expected /path?from=2013-11-22, got", resp)
	}

	settings = settings.ReportTo(now)
	resp = GetReportQuery("/path", settings)
	if resp != "/path?from=2013-11-22&to=2013-11-22" {
		t.Fatal("Expected /path?from=2013-11-22&to=2013-11-22, got", resp)
	}

	settings.From = nil
	resp = GetReportQuery("/path", settings)
	if resp != "/path?to=2013-11-22" {
		t.Fatal("Expected /path?to=2013-11-22, got", resp)
	}

	settings = settings.ReportGrouping("group by")
	resp = GetReportQuery("/path", settings)
	if resp != "/path?grouping=group+by&to=2013-11-22" {
		t.Fatal("Expected /path?grouping=group+by&to=2013-11-22, got", resp)
	}

	settings.To = nil
	resp = GetReportQuery("/path", settings)
	if resp != "/path?grouping=group+by" {
		t.Fatal("Expected /path?grouping=group+by, got", resp)
	}
}

func TestGetReportFrom(t *testing.T) {
	now := time.Date(2013, time.November, 22, 0, 0, 0, 0, time.UTC)
	settings := ReportFrom(now)
	if settings.From == nil || *settings.From != now {
		t.Fatal("Expected now", settings.From)
	}
	if settings.To != nil {
		t.Fatal("Expected nil", settings.To)
	}
	if settings.Grouping != nil {
		t.Fatal("Expected nil", settings.Grouping)
	}
}

func TestGetReportTo(t *testing.T) {
	now := time.Date(2013, time.November, 22, 0, 0, 0, 0, time.UTC)
	settings := ReportTo(now)
	if settings.From != nil {
		t.Fatal("Expected nil", settings.From)
	}
	if settings.To == nil || *settings.To != now {
		t.Fatal("Expected now", settings.To)
	}
	if settings.Grouping != nil {
		t.Fatal("Expected nil", settings.Grouping)
	}
}

func TestGetReportGrouping(t *testing.T) {
	settings := ReportGrouping("group by")
	if settings.From != nil {
		t.Fatal("Expected nil", settings.From)
	}
	if settings.To != nil {
		t.Fatal("Expected nil", settings.To)
	}
	if settings.Grouping == nil || *settings.Grouping != "group by" {
		t.Fatal("Expected group by", settings.Grouping)
	}
}

func TestGetVodUsage(t *testing.T) {
	expectedStatus := http.StatusOK
	returnBody := true

	mux := http.NewServeMux()
	mux.HandleFunc("/reports/vod", func(w http.ResponseWriter, r *http.Request) {
		if expectedStatus != http.StatusOK {
			w.WriteHeader(expectedStatus)
			return
		}

		if !returnBody {
			return
		}

		fmt.Fprintln(w, `{
  "total": {
    "encoded_minutes": 6,
    "billable_minutes": 8
  },
  "statistics": [
    {
      "grouping": "zencoder",
      "collected_on": "2011-08-10",
      "encoded_minutes": 4,
      "billable_minutes": 5
    }, {
      "grouping": null,
      "collected_on": "2011-08-10",
      "encoded_minutes": 2,
      "billable_minutes": 3
    }
  ]
}`)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	details, err := zc.GetVodUsage(nil)
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if details == nil {
		t.Fatal("Expected details")
	}

	if details.Total.EncodedMinutes != 6 {
		t.Fatal("Expected 6, got", details.Total.EncodedMinutes)
	}

	if details.Total.BillableMinutes != 8 {
		t.Fatal("Expected 8, got", details.Total.BillableMinutes)
	}

	if len(details.Statistics) != 2 {
		t.Fatal("Expected 2 statistics, got", len(details.Statistics))
	}

	if details.Statistics[0].Grouping != "zencoder" {
		t.Fatal("Expected zencoder, got", details.Statistics[0].Grouping)
	}

	if details.Statistics[0].CollectedOn != "2011-08-10" {
		t.Fatal("Expected 2011-08-10, got", details.Statistics[0].CollectedOn)
	}

	if details.Statistics[0].EncodedMinutes != 4 {
		t.Fatal("Expected 4, got", details.Statistics[0].EncodedMinutes)
	}

	if details.Statistics[0].BillableMinutes != 5 {
		t.Fatal("Expected 5, got", details.Statistics[0].BillableMinutes)
	}

	expectedStatus = http.StatusInternalServerError

	details, err = zc.GetVodUsage(nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}

	expectedStatus = http.StatusOK
	returnBody = false

	details, err = zc.GetVodUsage(nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}

	srv.Close()
	returnBody = false

	details, err = zc.GetVodUsage(nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}
}

func TestGetLiveUsage(t *testing.T) {
	expectedStatus := http.StatusOK
	returnBody := true

	mux := http.NewServeMux()
	mux.HandleFunc("/reports/live", func(w http.ResponseWriter, r *http.Request) {
		if expectedStatus != http.StatusOK {
			w.WriteHeader(expectedStatus)
			return
		}

		if !returnBody {
			return
		}

		fmt.Fprintln(w, `{
  "total": {
    "stream_hours": 5,
    "billable_stream_hours": 6,
    "encoded_hours": 5,
    "billable_encoded_hours": 6,
    "total_hours": 10,
    "total_billable_hours": 12
  },
  "statistics": [
    {
      "grouping": "zencoder",
      "collected_on": "2011-08-10",
      "stream_hours": 3,
      "billable_stream_hours": 4,
      "encoded_hours": 3,
      "billable_encoded_hours": 4,
      "total_hours": 6,
      "total_billable_hours": 8
    }, {
      "grouping": null,
      "collected_on": "2011-08-10",
      "stream_hours": 2,
      "billable_stream_hours": 2,
      "encoded_hours": 2,
      "billable_encoded_hours": 2,
      "total_hours": 4,
      "total_billable_hours": 4
    }
  ]
}`)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	details, err := zc.GetLiveUsage(nil)
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if details == nil {
		t.Fatal("Expected details")
	}

	if details.Total.StreamHours != 5 {
		t.Fatal("Expected 5, got", details.Total.StreamHours)
	}
	if details.Total.BillableStreamHours != 6 {
		t.Fatal("Expected 6, got", details.Total.BillableStreamHours)
	}
	if details.Total.EncodedHours != 5 {
		t.Fatal("Expected 5, got", details.Total.EncodedHours)
	}
	if details.Total.BillableEncodedHours != 6 {
		t.Fatal("Expected 6, got", details.Total.BillableEncodedHours)
	}
	if details.Total.TotalHours != 10 {
		t.Fatal("Expected 10, got", details.Total.TotalHours)
	}
	if details.Total.TotalBillableHours != 12 {
		t.Fatal("Expected 12, got", details.Total.TotalBillableHours)
	}

	if len(details.Statistics) != 2 {
		t.Fatal("Expected 2 statistics, got", len(details.Statistics))
	}

	if details.Statistics[0].Grouping != "zencoder" {
		t.Fatal("Expected zencoder, got", details.Statistics[0].Grouping)
	}
	if details.Statistics[0].CollectedOn != "2011-08-10" {
		t.Fatal("Expected 2011-08-10, got", details.Statistics[0].CollectedOn)
	}
	if details.Statistics[0].StreamHours != 3 {
		t.Fatal("Expected 3, got", details.Statistics[0].StreamHours)
	}
	if details.Statistics[0].BillableStreamHours != 4 {
		t.Fatal("Expected 4, got", details.Statistics[0].BillableStreamHours)
	}
	if details.Statistics[0].EncodedHours != 3 {
		t.Fatal("Expected 3, got", details.Statistics[0].EncodedHours)
	}
	if details.Statistics[0].BillableEncodedHours != 4 {
		t.Fatal("Expected 4, got", details.Statistics[0].BillableEncodedHours)
	}
	if details.Statistics[0].TotalHours != 6 {
		t.Fatal("Expected 6, got", details.Statistics[0].TotalHours)
	}
	if details.Statistics[0].TotalBillableHours != 8 {
		t.Fatal("Expected 8, got", details.Statistics[0].TotalBillableHours)
	}

	expectedStatus = http.StatusInternalServerError

	details, err = zc.GetLiveUsage(nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}

	expectedStatus = http.StatusOK
	returnBody = false

	details, err = zc.GetLiveUsage(nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}

	srv.Close()
	returnBody = false

	details, err = zc.GetLiveUsage(nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}
}

func TestGetUsage(t *testing.T) {
	expectedStatus := http.StatusOK
	returnBody := true

	mux := http.NewServeMux()
	mux.HandleFunc("/reports/all", func(w http.ResponseWriter, r *http.Request) {
		if expectedStatus != http.StatusOK {
			w.WriteHeader(expectedStatus)
			return
		}

		if !returnBody {
			return
		}

		fmt.Fprintln(w, `{
  "total": {
    "live": {
      "stream_hours": 5,
      "billable_stream_hours": 6,
      "encoded_hours": 5,
      "billable_encoded_hours": 6,
      "total_hours": 10,
      "total_billable_hours": 12
    },
    "vod": {
      "encoded_minutes": 6,
      "billable_minutes": 8
    }
  },
  "statistics": {
    "live": [
      {
        "grouping": "zencoder",
        "collected_on": "2011-08-10",
        "stream_hours": 3,
        "billable_stream_hours": 4,
        "encoded_hours": 3,
        "billable_encoded_hours": 4,
        "total_hours": 6,
        "total_billable_hours": 8
      }, {
        "grouping": null,
        "collected_on": "2011-08-10",
        "stream_hours": 2,
        "billable_stream_hours": 2,
        "encoded_hours": 2,
        "billable_encoded_hours": 2,
        "total_hours": 4,
        "total_billable_hours": 4
      }
    ],
    "vod": [
      {
        "grouping": "zencoder",
        "collected_on": "2011-08-10",
        "encoded_minutes": 4,
        "billable_minutes": 5
      }, {
        "grouping": null,
        "collected_on": "2011-08-10",
        "encoded_minutes": 2,
        "billable_minutes": 3
      }
    ]
  }
}`)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	details, err := zc.GetUsage(nil)
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if details == nil {
		t.Fatal("Expected details")
	}

	if details.Total.Live.StreamHours != 5 {
		t.Fatal("Expected 5, got", details.Total.Live.StreamHours)
	}
	if details.Total.Live.BillableStreamHours != 6 {
		t.Fatal("Expected 6, got", details.Total.Live.BillableStreamHours)
	}
	if details.Total.Live.EncodedHours != 5 {
		t.Fatal("Expected 5, got", details.Total.Live.EncodedHours)
	}
	if details.Total.Live.BillableEncodedHours != 6 {
		t.Fatal("Expected 6, got", details.Total.Live.BillableEncodedHours)
	}
	if details.Total.Live.TotalHours != 10 {
		t.Fatal("Expected 10, got", details.Total.Live.TotalHours)
	}
	if details.Total.Live.TotalBillableHours != 12 {
		t.Fatal("Expected 12, got", details.Total.Live.TotalBillableHours)
	}
	if details.Total.Vod.EncodedMinutes != 6 {
		t.Fatal("Expected 6, got", details.Total.Vod.EncodedMinutes)
	}
	if details.Total.Vod.BillableMinutes != 8 {
		t.Fatal("Expected 8, got", details.Total.Vod.BillableMinutes)
	}
	if len(details.Statistics.Live) != 2 {
		t.Fatal("Expected 2 live stats, got", len(details.Statistics.Live))
	}

	if details.Statistics.Live[0].Grouping != "zencoder" {
		t.Fatal("Expected zencoder, got", details.Statistics.Live[0].Grouping)
	}
	if details.Statistics.Live[0].CollectedOn != "2011-08-10" {
		t.Fatal("Expected 2011-08-10, got", details.Statistics.Live[0].CollectedOn)
	}
	if details.Statistics.Live[0].StreamHours != 3 {
		t.Fatal("Expected 3, got", details.Statistics.Live[0].StreamHours)
	}
	if details.Statistics.Live[0].BillableStreamHours != 4 {
		t.Fatal("Expected 4, got", details.Statistics.Live[0].BillableStreamHours)
	}
	if details.Statistics.Live[0].EncodedHours != 3 {
		t.Fatal("Expected 3, got", details.Statistics.Live[0].EncodedHours)
	}
	if details.Statistics.Live[0].BillableEncodedHours != 4 {
		t.Fatal("Expected 4, got", details.Statistics.Live[0].BillableEncodedHours)
	}
	if details.Statistics.Live[0].TotalHours != 6 {
		t.Fatal("Expected 6, got", details.Statistics.Live[0].TotalHours)
	}
	if details.Statistics.Live[0].TotalBillableHours != 8 {
		t.Fatal("Expected 8, got", details.Statistics.Live[0].TotalBillableHours)
	}

	if len(details.Statistics.Vod) != 2 {
		t.Fatal("Expected 2 VOD stats, got", len(details.Statistics.Vod))
	}

	if details.Statistics.Vod[0].Grouping != "zencoder" {
		t.Fatal("Expected zencoder, got", details.Statistics.Vod[0].Grouping)
	}
	if details.Statistics.Vod[0].CollectedOn != "2011-08-10" {
		t.Fatal("Expected 2011-08-10, got", details.Statistics.Vod[0].CollectedOn)
	}
	if details.Statistics.Vod[0].EncodedMinutes != 4 {
		t.Fatal("Expected 4, got", details.Statistics.Vod[0].EncodedMinutes)
	}
	if details.Statistics.Vod[0].BillableMinutes != 5 {
		t.Fatal("Expected 5, got", details.Statistics.Vod[0].BillableMinutes)
	}

	expectedStatus = http.StatusInternalServerError

	details, err = zc.GetUsage(nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}

	expectedStatus = http.StatusOK
	returnBody = false

	details, err = zc.GetUsage(nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}

	srv.Close()
	returnBody = false

	details, err = zc.GetUsage(nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}
}
