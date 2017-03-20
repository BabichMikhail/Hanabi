package statistics

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/tealeg/xlsx"
)

func (stat *Stat) SaveStat(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(stat)
	if err != nil {
		return err
	}

	f.WriteString(string(bytes))
	return f.Close()
}

func (stat *Stat) SaveValues(path string) error {
	excel := xlsx.NewFile()
	sheetPoints, err := excel.AddSheet("points")
	if err != nil {
		return err
	}

	histData := map[int]int{}
	for i := 0; i <= 25; i++ {
		histData[i] = 0
	}
	for i, points := range stat.Values {
		sheetPoints.Cell(i, 0).SetFloat(points)

		count, _ := histData[int(points)]
		histData[int(points)] = count + 1
	}

	sheetHistogram, err := excel.AddSheet("histogram")
	if err != nil {
		return err
	}

	for points, count := range histData {
		sheetHistogram.Cell(points, 0).SetInt(points)
		sheetHistogram.Cell(points, 1).SetInt(count)
	}

	return excel.Save(path)
}

func (stat *Stat) SaveToFile(dir, subdir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		return err
	}

	statDir := path.Join(dir, subdir)
	err = os.Mkdir(statDir, os.ModeDir)
	if err != nil {
		return err
	}

	err = stat.SaveStat(path.Join(statDir, "stat.json"))
	if err != nil {
		return err
	}

	err = stat.SaveValues(path.Join(statDir, "distribution.xlsx"))
	return err
}

func ReadFromFile(dir, subdir string) (*Stat, error) {
	bytes, err := ioutil.ReadFile(path.Join(dir, subdir, "stat.json"))
	if err != nil {
		return nil, err
	}

	var stat *Stat
	err = json.Unmarshal(bytes, stat)
	return stat, err
}
