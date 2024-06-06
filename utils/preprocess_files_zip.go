package utils

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"math"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type LidarSingleFile struct {
	DateTime time.Time
	ProfLen  uint
	Count    uint
	RepRate  uint
	Data     []float64
	Step     float64
	ChType   string
}

type LidarOneMeasurement struct {
	DateTime time.Time
	ProfLen  int
	Count    int
	RepRate  int
	Step     float64
	Dat      []float64
	Dak      []float64
}

func FilterProfile(input []byte) ([]byte, error) {

	i := 0
	filtered := make([]byte, 0)

	for i < len(input)-1 {
		if (input[i] == 0x0d) && (input[i+1] == 0x0a) {
			filtered = append(filtered, input[i+1])
			i += 1
		} else {
			filtered = append(filtered, input[i])
		}
		i += 1
	}
	filtered = append(filtered, input[i])

	return filtered, nil
}

func CombineChannels(a, b []*LidarSingleFile) ([]LidarOneMeasurement, error) {
	ret := make([]LidarOneMeasurement, 0)

	for _, f1 := range a {
		for _, f2 := range b {
			if f1.DateTime == f2.DateTime {
				ret = append(ret, LidarOneMeasurement{
					DateTime: f1.DateTime,
					ProfLen:  int(f1.ProfLen),
					Count:    1,
					RepRate:  int(f1.RepRate),
					Step:     f1.Step,
					Dat:      f1.Data,
					Dak:      f2.Data,
				})
			}
		}
	}
	return ret, nil
}

func MakeAverageProfiles(data []LidarOneMeasurement, seconds int) ([]LidarOneMeasurement, error) {
	// Заполнить
	time.Now().Add(time.Second)
}

func ReadZippedLidarFile(filez *zip.File, step float64) (*LidarSingleFile, error) {
	var ret *LidarSingleFile
	if file, err := filez.Open(); err == nil {
		buffer := make([]byte, filez.UncompressedSize64)
		file.Read(buffer)
		file.Close()

		fileLen := len(buffer) - 18
		log2filelen := math.Log2(float64(fileLen))
		if log2filelen != math.Floor(log2filelen) {
			buffer, err = FilterProfile(buffer)
			if err != nil {
				return nil, err
			}
		}
		data16 := make([]uint16, len(buffer)/2)
		for i := range data16 {
			data16[i] = uint16(binary.LittleEndian.Uint16(buffer[i*2 : (i+1)*2]))
		}
		ret = &LidarSingleFile{}
		ret.DateTime = time.Date(int(data16[0]), time.Month(data16[1]+1), int(data16[2]),
			int(data16[3]), int(data16[4]), int(data16[5]), 0, time.UTC)
		ret.ProfLen = uint(data16[6])
		ret.Count = uint(data16[7])
		ret.RepRate = uint(data16[8])
		ret.Step = step
		recorded_size := ret.ProfLen * ret.Count
		profile := data16[9 : 9+recorded_size]

		ret.Data = make([]float64, ret.ProfLen)

		offset := 0
		for i := 0; i < int(ret.Count); i++ {
			for j := 0; j < int(ret.ProfLen); j++ {
				ret.Data[j] = ret.Data[j] + float64(profile[j+offset])
			}
			offset += int(ret.ProfLen)
		}

		ret.ChType = filepath.Ext(filez.Name)
		ret.Step = step
	}
	return ret, nil
}

// func ReadLidarFile(file_name *multipart.FileHeader, step float64) (*LidarSingleFile, error) {
// 	var ret *LidarSingleFile = nil
// 	if file, err := file_name.Open(); err == nil {

// 		buffer := make([]byte, file_name.Size)
// 		file.Read(buffer)
// 		file.Close()
// 		fileLen := len(buffer) - 18
// 		log2filelen := math.Log2(float64(fileLen))
// 		if log2filelen != math.Floor(log2filelen) {
// 			buffer, _ = FilterProfile(buffer)
// 		}
// 		data16 := make([]uint16, len(buffer)/2)
// 		for i := range data16 {
// 			data16[i] = uint16(binary.LittleEndian.Uint16(buffer[i*2 : (i+1)*2]))
// 		}
// 		ret = &LidarSingleFile{}
// 		ret.DateTime = time.Date(int(data16[0]), time.Month(data16[1]+1), int(data16[2]),
// 			int(data16[3]), int(data16[4]), int(data16[5]), 0, time.UTC)
// 		ret.ProfLen = uint(data16[6])
// 		ret.Count = uint(data16[7])
// 		ret.RepRate = uint(data16[8])
// 		ret.Step = step
// 		recorded_size := ret.ProfLen * ret.Count
// 		profile := data16[9 : 9+recorded_size]

// 		ret.Data = make([]float64, ret.ProfLen)

// 		offset := 0
// 		for i := 0; i < int(ret.Count); i++ {
// 			for j := 0; j < int(ret.ProfLen); j++ {
// 				ret.Data[j] = ret.Data[j] + float64(profile[j+offset])
// 			}
// 			offset += int(ret.ProfLen)
// 		}

// 		ret.ChType = filepath.Ext(file_name.Filename)
// 		ret.Step = step

// 	}
// 	return ret, nil
// }

// func ReadLidarFiles(files_a, files_b *multipart.FileHeader, step float64) ([]LidarOneMeasurement, error) {
// 	ArrayDat := make([]*LidarSingleFile, len(files_a))
// 	ArrayDak := make([]*LidarSingleFile, len(files_b))

// 	for i, file := range files_a {
// 		ArrayDat[i], _ = ReadLidarFile(file, step)
// 	}

// 	for i, file := range files_b {
// 		ArrayDak[i], _ = ReadLidarFile(file, step)
// 	}

// 	return CombineChannels(ArrayDat, ArrayDak)
// }

// TODO: Отрефракторить
func ReadZippedArrayLidarFile(filezip *multipart.FileHeader, step float64) ([]*LidarSingleFile, []*LidarSingleFile, error) {
	ArrayDAT := make([]*LidarSingleFile, 0)
	ArrayDAK := make([]*LidarSingleFile, 0)
	var err error
	var file multipart.File
	if file, err = filezip.Open(); err == nil {
		buf, _ := io.ReadAll(file)
		log.Println(len(buf))
		buffer := bytes.NewReader(buf)

		r1, _ := zip.NewReader(buffer, int64(buffer.Len()))

		for _, f := range r1.File {
			log.Println(f.Name)
			if (!f.FileInfo().IsDir()) && (strings.HasSuffix(f.Name, ".dat") || strings.HasSuffix(f.Name, ".dak")) {
				tmp, _ := ReadZippedLidarFile(f, step)
				if strings.HasSuffix(f.Name, ".dat") {
					ArrayDAT = append(ArrayDAT, tmp)
				} else if strings.HasSuffix(f.Name, ".dak") {
					ArrayDAK = append(ArrayDAK, tmp)
				}
			}
		}
		file.Close()
	} else {
		log.Println("Error", err)
		return nil, nil, errors.New("ошибка чтения файла")
	}

	return ArrayDAT, ArrayDAK, nil
}

func ReadZippedLidarArchive(files *multipart.FileHeader, step float64) ([]LidarOneMeasurement, error) {
	DAT, DAK, err := ReadZippedArrayLidarFile(files, step)
	if err != nil {
		return nil, err
	}
	log.Println(len(DAT), len(DAK))
	return CombineChannels(DAT, DAK)
}
