package helpers

import (
	"bufio"
	"os"
)

func ReadFileLns(name string) ([][]byte, error) {
	var ret [][]byte

	fle, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fle.Close()

	scanner := bufio.NewScanner(fle)
	for scanner.Scan() {
		line := scanner.Bytes()
		ret = append(ret, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}
