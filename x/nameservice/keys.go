package nameservice

var (
	ValuePrefix = []byte{0x00}
	OwnerPrefix = []byte{0x01}
	PricePrefix = []byte{0x02}
)

func ValueKey(domain string) []byte {
	return append(ValuePrefix, []byte(domain)...)
}

func OwnerKey(domain string) []byte {
	return append(OwnerPrefix, []byte(domain)...)
}

func PriceKey(domain string) []byte {
	return append(PricePrefix, []byte(domain)...)
}
