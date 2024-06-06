package decoder

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func binaryToText(binaryString string) (string, error) {
	// Split the binary string into individual binary numbers
	binaryArray := strings.Split(binaryString, " ")

	var result strings.Builder

	// Iterate over each binary number
	for _, binary := range binaryArray {
		// Convert the binary number to an integer
		decimal, err := strconv.ParseInt(binary, 2, 64)
		if err != nil {
			return "", err
		}

		// Convert the integer to its corresponding ASCII character
		result.WriteByte(byte(decimal))
	}

	return result.String(), nil
}

func decodeFile(fileName string) (string, error) {
	file, _ := os.Open(filepath.Join("pkg/decoder/files", fileName))
	fileBytes, _ := io.ReadAll(file)

	str := string(fileBytes)

	defer file.Close()
	return binaryToText(str)
}

func Decode() error {
	files, err := os.ReadDir("pkg/decoder/files")

	filesDecodedMap := make(map[int]string)

	if err != nil {
		return errors.New("error while reading the dir")
	}

	for _, file := range files {
		fileName := file.Name()
		bytes, _ := base64.StdEncoding.DecodeString(fileName)
		convertedName := string(bytes)
		integerVal, _ := strconv.Atoi(convertedName)
		filesDecodedMap[integerVal] = fileName
	}

	decodedStrBuilder := strings.Builder{}
	for i := 0; i < len(files); i++ {
		fileName, ok := filesDecodedMap[i]

		if !ok {
			continue
		}

		decodedStr, decodeErr := decodeFile(fileName)

		if decodeErr != nil {
			continue
		}

		decodedStrBuilder.WriteString(decodedStr)
	}

	fmt.Println(decodedStrBuilder.String())
	return nil

}
