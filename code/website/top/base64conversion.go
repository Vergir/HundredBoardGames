package top

import (
	"bytes"
	"compress/lzw"
	"encoding/base64"
	"encoding/binary"
	"io"
)

func gamesIdsToBase64(gamesIds []uint) (string, error) {
	var lzwBytes bytes.Buffer

	lzwWriter := lzw.NewWriter(&lzwBytes, lzw.MSB, 8)
	bytesBuffer := make([]byte, 8)
	for _, gameId := range gamesIds {
		binary.BigEndian.PutUint64(bytesBuffer, uint64(gameId))
		lzwWriter.Write(bytesBuffer)
	}
	lzwWriter.Close()

	b64gamesIds := base64.URLEncoding.EncodeToString(lzwBytes.Bytes())

	return b64gamesIds, nil
}

func Base64ToGamesIds(base64string string) ([]uint, error) {
	compressedGamesIdsBytes, err := base64.URLEncoding.DecodeString(base64string)
	if err != nil {
		return nil, err
	}

	lzwReader := lzw.NewReader(bytes.NewBuffer(compressedGamesIdsBytes), lzw.MSB, 8)
	decompressedGamesIdsBytes, err := io.ReadAll(lzwReader)
	if err != nil {
		return nil, err
	}
	lzwReader.Close()

	gamesIdsUint64 := make([]uint64, len(decompressedGamesIdsBytes)/8) //8 = length of uint64

	err = binary.Read(bytes.NewReader(decompressedGamesIdsBytes), binary.BigEndian, gamesIdsUint64)
	if err != nil {
		return nil, err
	}

	gameIds := make([]uint, len(gamesIdsUint64))
	for i, gameIdUint64 := range gamesIdsUint64 {
		gameIds[i] = uint(gameIdUint64)
	}

	return gameIds, nil
}
