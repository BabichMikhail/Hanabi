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

func (stat *Stat) setSheetPoints(excel *xlsx.File, histData map[int]int) error {
	sh, err := excel.AddSheet("points")
	if err != nil {
		return err
	}

	for i, points := range stat.Values {
		sh.Cell(i, 0).SetFloat(points)

		count, _ := histData[int(points)]
		histData[int(points)] = count + 1
	}
	return nil
}

func (stat *Stat) setSheetHistogram(excel *xlsx.File, histData map[int]int) error {
	sh, err := excel.AddSheet("histogram")
	if err != nil {
		return err
	}

	for points, count := range histData {
		sh.Cell(points, 0).SetInt(points)
		sh.Cell(points, 1).SetInt(count)
	}
	return nil
}

func (stat *Stat) SaveValues(path string) error {
	excel := xlsx.NewFile()

	histData := map[int]int{}
	for i := 0; i <= 25; i++ {
		histData[i] = 0
	}

	if err := stat.setSheetPoints(excel, histData); err != nil {
		return err
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
