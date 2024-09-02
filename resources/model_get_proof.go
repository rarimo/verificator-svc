/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type GetProof struct {
	Key
	Attributes GetProofAttributes `json:"attributes"`
}
type GetProofRequest struct {
	Data     GetProof `json:"data"`
	Included Included `json:"included"`
}

type GetProofListRequest struct {
	Data     []GetProof `json:"data"`
	Included Included   `json:"included"`
	Links    *Links     `json:"links"`
}

// MustGetProof - returns GetProof from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustGetProof(key Key) *GetProof {
	var getProof GetProof
	if c.tryFindEntry(key, &getProof) {
		return &getProof
	}
	return nil
}
