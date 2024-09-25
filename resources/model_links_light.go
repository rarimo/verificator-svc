/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type LinksLight struct {
	Key
	Attributes LinksLightAttributes `json:"attributes"`
}
type LinksLightRequest struct {
	Data     LinksLight `json:"data"`
	Included Included   `json:"included"`
}

type LinksLightListRequest struct {
	Data     []LinksLight `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustLinksLight - returns LinksLight from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustLinksLight(key Key) *LinksLight {
	var linksLight LinksLight
	if c.tryFindEntry(key, &linksLight) {
		return &linksLight
	}
	return nil
}
