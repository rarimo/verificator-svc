/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Links struct {
	Key
	Attributes LinksAttributes `json:"attributes"`
}
type LinksRequest struct {
	Data     Links    `json:"data"`
	Included Included `json:"included"`
}

type LinksListRequest struct {
	Data     []Links  `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustLinks - returns Links from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustLinks(key Key) *Links {
	var links Links
	if c.tryFindEntry(key, &links) {
		return &links
	}
	return nil
}
