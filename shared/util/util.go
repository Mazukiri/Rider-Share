package util

import "fmt"

// GetRandomAvatar returns a random avatar URL from the randomuser.com API
func GetRandomAvatar(index int) string {
	return fmt.Sprintf("https://randomuser.com/api/portraits/lego/%d.jpg", index)
}
