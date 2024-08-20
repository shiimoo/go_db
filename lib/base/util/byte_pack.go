package util

const packBytesLimit = 8 //65535

func SubPack(bs []byte) [][]byte {
	subs := make([][]byte, 0)
	if bs == nil {
		return subs
	}
	length := len(bs)
	startIndex := 0
	endIndex := 0
	for startIndex < length-1 {
		endIndex = startIndex + packBytesLimit
		if endIndex >= length {
			endIndex = length
		}
		subs = append(subs, bs[startIndex:endIndex])
		startIndex = endIndex
	}
	return subs
}
