/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Parameters struct {
	Key
	Attributes ParametersAttributes `json:"attributes"`
}
type ParametersResponse struct {
	Data     Parameters `json:"data"`
	Included Included   `json:"included"`
}

type ParametersListResponse struct {
	Data     []Parameters `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustParameters - returns Parameters from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustParameters(key Key) *Parameters {
	var parameters Parameters
	if c.tryFindEntry(key, &parameters) {
		return &parameters
	}
	return nil
}
