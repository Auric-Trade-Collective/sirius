package helpers

func SliceUp(buff []byte, delim byte) [][]byte {
	var ret [][]byte

	totalLen := len(buff)

	lastEnd := 0
	for i := lastEnd; i < totalLen; i++ {
		if buff[i] == delim {
			slice := buff[lastEnd:i]
			ret = append(ret, slice)

			lastEnd = i + 1
			continue
		}

		if i + 1 == totalLen {
			slice := buff[lastEnd:]
			ret = append(ret, slice)
			break
		}
	}

	return ret
}
