package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sort"

	"github.com/kmjayadeep/totpm/pkg/configfile"
	"github.com/kmjayadeep/totpm/pkg/types"
)

func (c *Cli) refreshCache() {
	u := **c.server
	u.Path = path.Join(u.Path, "api/site")

	config, err := configfile.Read()
	c.handleError(err)

	client := &http.Client{}

	req, err := http.NewRequest("GET", u.String(), nil)
	c.handleError(err)

	req.Header.Set("x-access-token", config.Token)

	res, err := client.Do(req)
	c.handleError(err)

	if res.StatusCode == http.StatusUnauthorized {
		c.handleError(fmt.Errorf("Invalid token... please login again"))
	}

	if res.StatusCode != http.StatusOK {
		c.handleError(fmt.Errorf("Unable to list OTPs"))
	}

	var otps []types.Otp
	if err := json.NewDecoder(res.Body).Decode(&otps); err != nil {
		c.handleError(fmt.Errorf("Unable to list OTPs"))
	}

	cache := configfile.NewCache()

	sort.Slice(otps, func(i, j int) bool {
		return otps[i].Name < otps[j].Name
	})
	cache.Otps = otps
	c.handleError(cache.Write())

}
