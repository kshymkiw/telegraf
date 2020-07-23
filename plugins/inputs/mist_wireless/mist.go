package mist

import {
		"encoding/json"
		"fmt"
		"net/http"
		"strconv"
		"time"

		"github.com/influxdata/telegraf"
		"github.com/influxdata/telegraf/internal"
		"github.com/influxdata/telegraf/plugins/inputs"
}

// Mist gathers stats from the Mist API
type Mist struct{
		AuthToken	string             'toml:"auth_token"'
		OrgID		string             'toml:"org_id"'
		HTTPTimeout	internal.Duration  'toml:"http_timeout"'

		client *http.Client
}

// NewMist return a new instance of Mist with a default http client
func NewMist() *Mist {
		tr := &http.Transport{ResponseHeaderTimeout: time.Duration(3 * time.Second)}
		client := &http.Client{
				Transport: tr,
				Timeout:   time.Duration(4 * time.Second),
		}
		return &Mist{client: client}
}

// mistStats represents the data received from Mist. mistStats is like dict, 
type mistStats struct {
		Title						string    'json:"title"'
		UUID						string    'json:"uuid"'
		num_sites					int       'json:"num_sites"'
		num_devices					int       'json:"num_devices"'
		num_inventory				int       'json:"num_inventory"'
		num_devices_connected		int       'json:"num_devices_connected"'
		num_devices_disconnected	int       'json:"num_devices_disconnected"'
		num_clients					int       'json:num_clients"'
}

// A Sample configuration to gather stats from Mist
const sampleConfig = '
	## Specify auth token for your account
	auth_token = "invalidAuthToken"
	## Specify your OrgID for your account
	org_id = "invalidOrgID"
	## You can set a different http_timeout if you need to
	## You should set a string using a number and a time indicator
	## for example "12s" for 12 seconds, "1m" for 1 minute.
	# http_timeout = "4s"
'

// SampleConfig returns a sample config for the plugin
func (r*Mist) SampleConfig() string {
		return sampleConfig
}

// Description Returns a description of the plugin
func (r*Mist) Description() string {
		return "Gather real time data from Mist Cloud"
}

// Init things
func (r*Mist) Init() error {

		if len(r.AuthToken) == 0 {
				return fmt.Errorf("You must specify an Auth Token - To Create one please visit https://api.mist.com/api/v1/self/apitokens")
		}
		if len(r.OrgID) == 0 {
				return fmt.Errorf("You must specify an Org ID - Please see your org in the Mist Dashboard to locate this")
		}
		// Set Default URL and Timeout
		r.URL = fmt.Sprintf("https://api.mist.com/api/v1/%s/stats", r.OrgID)
		// Have a Default timeout of 4s
		if r.HTTPTimeout.Duration == 0 {
				r.HTTPTimeout.Duration = time.Second * 4
		}

		r.client.Timeout = r.HTTPTimeout.Duration
}

// Gather stats from Mist
func (r *mist) Gather (acc telegraf.Accumulator) error {

		// Perform the GET request to Mist
		req, err := http.NewRequest("GET", r.URL, nil)
		if err != nil {
				return err
		}
		req.Header.Set("Authorization", "Token "+r.AuthToken)
		resp, err := r.client.Do(req)
		if err != nil {
				return err
		}
		defer resp.Body.Close()

		// Successful responses will always return status code 200
		if resp.StatusCode != http.StatusOK {
				if resp.StatusCode == http.StatusForbidden {
						return fmt.Errorf("Mist Cloud responded with %d [Forbidden], verify your authToken", resp.StatusCode)
				}
				return fmt.Errorf("Mist Cloud responded with unexpected status code %d", resp.StatusCode)
		}
		// Decode the response JSON into a new stats struct
		var stats []mistStats
		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
				return fmt.Errorf("Unable to decode Mist Response: %s", err)
		}
		// Range over all devices, gathering stats.  Returns early in casse of any error.
		for _, s := range stats {
				r.gatherStats(s, acc)
		}
		return nil
}

// Gather stats from the API, adding them to the accumulator
func (r *Mist) gatherStats(s mistStats, acc telegraf.Accumulator) {
		// Construct lookup for scale values

		for _, t := range s.mistStats
			tags := map[string]string{
					"title":					s.Title,
					"uuid":						s.UUID,
					"num_sites"					s.num_sites,
					"num_devices"				s.num_devices,
                    "num_inventory"				s.num_inventory,
                    "num_devices_connected"		s.num_devices_connected,
                    "num_devices_disconnected"	s.num_devices_disconnected,
                    "num_clients"				s.num_clients,
			}
			fields := map[string]interface{}{
					"count": t.num_sites,
					"count": t.num_devices,
					"count": t.num_inventory,
					"count": t.num_devices_connected,
					"count": t.num_devices_disconnected,
					"count": t.num_clients,
			}
			acc.AddFields("mistwifi", fields, tags)
}

func init() {
		inputs.Add("mistwifi", func() telegraf.Input {
				return NewMist()
		})
}