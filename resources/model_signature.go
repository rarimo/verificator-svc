/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Signature struct {
	Key
	Attributes SignatureAttributes `json:"attributes"`
}
type SignatureRequest struct {
	Data     Signature `json:"data"`
	Included Included  `json:"included"`
}

type SignatureListRequest struct {
	Data     []Signature `json:"data"`
	Included Included    `json:"included"`
	Links    *Links      `json:"links"`
}

// MustSignature - returns Signature from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustSignature(key Key) *Signature {
	var signature Signature
	if c.tryFindEntry(key, &signature) {
		return &signature
	}
	return nil
}
