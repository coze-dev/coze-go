package util

import (
	"encoding/binary"
	"os"
)

func WritePCMToWavFile(file string, audioPCMData []byte) error {
	outFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// WAV 文件头信息
	var (
		chunkID       = []byte{'R', 'I', 'F', 'F'}
		format        = []byte{'W', 'A', 'V', 'E'}
		subchunk1ID   = []byte{'f', 'm', 't', ' '}
		subchunk1Size = uint32(16) // PCM
		audioFormat   = uint16(1)  // PCM = 1 (线性量化)
		numChannels   = uint16(1)  // Mono = 1, Stereo = 2
		sampleRate    = uint32(24000)
		byteRate      = sampleRate * uint32(numChannels) * uint32(audioFormat) // SampleRate * NumChannels * BitsPerSample/8
		blockAlign    = numChannels * uint16(audioFormat)
		bitsPerSample = uint16(16)
		subchunk2ID   = []byte{'d', 'a', 't', 'a'}
	)

	// 预留空间写入 ChunkSize 和 Subchunk2Size
	if _, err := outFile.Seek(44, 0); err != nil {
		return err
	}

	// 模拟音频数据
	for i := 0; i < len(audioPCMData)-1; i += 2 {
		err := binary.Write(outFile, binary.LittleEndian, audioPCMData[i:i+2])
		if err != nil {
			return err
		}
	}

	// 获取文件大小
	fileInfo, err := outFile.Stat()
	if err != nil {
		return err
	}

	// 计算 ChunkSize 和 Subchunk2Size
	fileSize := fileInfo.Size()
	chunkSize := uint32(fileSize - 8)
	subchunk2Size := uint32(fileSize - 44)

	// 回写 WAV 文件头
	if _, err := outFile.Seek(0, 0); err != nil {
		return err
	}

	headers := [][]interface{}{
		{chunkID, chunkSize, format},
		{subchunk1ID, subchunk1Size, audioFormat, numChannels, sampleRate, byteRate, blockAlign, bitsPerSample},
		{subchunk2ID, subchunk2Size},
	}

	for _, headerSection := range headers {
		for _, headerField := range headerSection {
			if err := binary.Write(outFile, binary.LittleEndian, headerField); err != nil {
				return err
			}
		}
	}
	return nil
}
