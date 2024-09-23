/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type SignatureAttributes struct {
	// Signed message, must be len(message) == 32
	Message string `json:"message"`
	// Signature, must be len(signature) == 64
	Signature string `json:"signature"`
}
