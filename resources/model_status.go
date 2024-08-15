/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Status struct {
	Key
	Attributes StatusAttributes `json:"attributes"`
}
type StatusResponse struct {
	Data     Status   `json:"data"`
	Included Included `json:"included"`
}

type StatusListResponse struct {
	Data     []Status `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustStatus - returns Status from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustStatus(key Key) *Status {
	var status Status
	if c.tryFindEntry(key, &status) {
		return &status
	}
	return nil
}
