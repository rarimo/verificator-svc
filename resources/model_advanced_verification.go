/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type AdvancedVerification struct {
	Key
	Attributes AdvancedVerificationAttributes `json:"attributes"`
}
type AdvancedVerificationRequest struct {
	Data     AdvancedVerification `json:"data"`
	Included Included             `json:"included"`
}

type AdvancedVerificationListRequest struct {
	Data     []AdvancedVerification `json:"data"`
	Included Included               `json:"included"`
	Links    *Links                 `json:"links"`
}

// MustAdvancedVerification - returns AdvancedVerification from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustAdvancedVerification(key Key) *AdvancedVerification {
	var advancedVerification AdvancedVerification
	if c.tryFindEntry(key, &advancedVerification) {
		return &advancedVerification
	}
	return nil
}
