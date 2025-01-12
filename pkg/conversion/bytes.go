package conversion

// MBToB converts megabytes to bytes
func MBToB(mb uint32) uint64 {
	return uint64(mb) * 1024 * 1024
}

// BToMB converts bytes to megabytes
func BToMB(bytes uint64) uint32 {
	return uint32(bytes / 1024 / 1024)
}
