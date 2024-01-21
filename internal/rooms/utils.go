package rooms

func SizeToString(size int32) string {
	var sizeString string
	switch size {
	case 0:
		sizeString = "Tiny"
	case 1:
		sizeString = "Small"
	case 2:
		sizeString = "Medium"
	case 3:
		sizeString = "Large"
	case 4:
		sizeString = "Huge"
	default:
		sizeString = "Invalid"
	}
	return sizeString
}
