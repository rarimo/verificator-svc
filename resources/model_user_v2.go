/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UserV2 struct {
	Key
	Attributes UserV2Attributes `json:"attributes"`
}
type UserV2Request struct {
	Data     UserV2   `json:"data"`
	Included Included `json:"included"`
}

type UserV2ListRequest struct {
	Data     []UserV2 `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustUserV2 - returns UserV2 from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustUserV2(key Key) *UserV2 {
	var userV2 UserV2
	if c.tryFindEntry(key, &userV2) {
		return &userV2
	}
	return nil
}
