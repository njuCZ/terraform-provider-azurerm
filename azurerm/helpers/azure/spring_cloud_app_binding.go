package azure

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func SuppressKeyDiff(_, old, new string, _ *schema.ResourceData) bool {
	if len(old) < 3 || len(new) < 3 {
		return false
	}
	return new[0:3] == old[0:3]
}

func FindValueInGeneratedProperties(str, key, delimiter string) string {
	indexOfKey := strings.Index(str, key)
	if indexOfKey == -1 {
		return ""
	}
	str = str[indexOfKey:]
	indexOfDelimiter := strings.Index(str, delimiter)
	if indexOfDelimiter == -1 {
		return ""
	}
	return str[len(key)+1 : indexOfDelimiter]
}
