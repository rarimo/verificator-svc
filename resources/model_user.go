/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type User struct {
	Key
	Attributes UserAttributes `json:"attributes"`
}
type UserRequest struct {
	Data     User     `json:"data"`
	Included Included `json:"included"`
}

type UserListRequest struct {
	Data     []User   `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustUser - returns User from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustUser(key Key) *User {
	var user User
	if c.tryFindEntry(key, &user) {
		return &user
	}
	return nil
}
