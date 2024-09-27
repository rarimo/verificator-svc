/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UserParams struct {
	Key
	Attributes UserParamsAttributes `json:"attributes"`
}
type UserParamsRequest struct {
	Data     UserParams `json:"data"`
	Included Included   `json:"included"`
}

type UserParamsListRequest struct {
	Data     []UserParams `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustUserParams - returns UserParams from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustUserParams(key Key) *UserParams {
	var userParams UserParams
	if c.tryFindEntry(key, &userParams) {
		return &userParams
	}
	return nil
}
