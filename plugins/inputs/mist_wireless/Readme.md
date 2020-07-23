# Mist Wireless Input Plugin

The mist_wireless plugin will allow you to collect Org level stattistics from the Mist Wireless Dashboard.

### Configuration

```toml
[[inputs.mist_wireless]]
  ## Specify auth token for your account
  auth_token = "invalidAuthToken"
  ## Specify your OrgID for your account
  org_id = "invalidOrgID"
  ## You can set a different http_timeout if you need to
  ## You should set a string using a number and a time indicator
  ## for example "12s" for 12 seconds, "1m" for 1 minute.
  # http_timeout = "4s"
```

#### auth_token

Mist does not offer simple login so you will need to create an authToken from the [Mist REST
API](https://api.mist.com/api/v1/self/apitokens).

#### org_id

You will need your ORG ID from the [Mist Dashboard](https://manage.mist.com/).  On the Menu select Organization->Settings.

#### http_timeout

If you need to increase the HTTP timeout, you can do so here. You can set this
value in seconds. The default value is four (4) seconds.

### Metrics

The Mist REST API docs have good examples of the data that is available,
currently this input only returns Org level stattistics.

- mistwifi
  - tags:
    - Number of Sites for your Org (num_sites)
    - Number of Devices in Use (num_devices)
    - Number of Devices Claimed (num_Devices)
    - Number of Devices Connected (num_devices_connected)
    - Number of Devices Disconnected (num_devices_disconnected)
    - Number of Clients Total (num_clients)
  - fields:
    - Count (unit)

### Example Output

This section shows example output in Line Protocol format.  You can often use
`telegraf --input-filter <plugin-name> --test` or use the `file` output to get
this information.

```
Add Test Here
```

#### Special Thanks

Special thanks to OMGKitteh for all his help, as this was my first Go project.  I couldn't have completed this without his help.  Also special thanks to the #go_nuts channel on irc.freenode.net
