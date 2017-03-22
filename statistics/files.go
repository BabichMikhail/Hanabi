package statistics

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"

	ai "github.com/BabichMikhail/Hanabi/AI"
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

	sh.Cell(0, 0).SetString("Points")
	sh.Cell(0, 1).SetString("Step")
	sh.Cell(0, 2).SetString("Red Tokens")
	for i, g := range stat.Games {
		sh.Cell(i+1, 0).SetInt(g.Points)
		sh.Cell(i+1, 1).SetInt(g.Step)
		sh.Cell(i+1, 2).SetInt(g.RedTokens)

		count, _ := histData[g.Points]
		histData[g.Points] = count + 1
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

func (stat *Stat) setSheetInfo(excel *xlsx.File) error {
	sh, err := excel.AddSheet("info")
	if err != nil {
		return err
	}

	sh.Cell(0, 0).SetString("AI Game")
	sh.Cell(1, 0).SetString("Created")
	sh.Cell(1, 1).SetString(time.Now().Format("15:04:05 02.01.2006"))
	sh.Cell(2, 0).SetString("Players count")
	sh.Cell(2, 1).SetInt(len(stat.AITypes))
	sh.Cell(3, 0).SetString("Games count")
	sh.Cell(3, 1).SetInt(stat.Count)
	sh.Cell(4, 0).SetString("AI Types")
	sh.Cell(5, 0).SetString("AI Names")
	for i, aiType := range stat.AITypes {
		sh.Cell(4, 1+i).SetInt(aiType)
		sh.Cell(5, 1+i).SetString(ai.DefaultUsernamePrefix(aiType))
	}
	sh.Cell(7, 0).SetString("RedTokensFail")
	sum := 0
	for _, g := range stat.Games {
		if g.RedTokens == 3 {
			sum++
		}
	}
	sh.Cell(7, 1).SetInt(sum)

	offset := 9
	sh.Cell(offset, 0).SetString("Medium")
	sh.Cell(offset, 1).SetFloat(stat.Medium)
	sh.Cell(offset+1, 0).SetString("Dispersion")
	sh.Cell(offset+1, 1).SetFloat(stat.Disp)
	sh.Cell(offset+2, 0).SetString("Asymmetry")
	sh.Cell(offset+2, 1).SetFloat(stat.Asym)
	sh.Cell(offset+3, 0).SetString("Kurtosis")
	sh.Cell(offset+3, 1).SetFloat(stat.Kurt)
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

	if err := stat.setSheetHistogram(excel, histData); err != nil {
		return err
	}

	if err := stat.setSheetInfo(excel); err != nil {
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

	err = stat.SaveValues(path.Join(statDir, "distribution"+subdir+".xlsx"))
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
