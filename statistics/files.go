package statistics

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strconv"
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

func (stat *Stat) setSheetPoints(excel *xlsx.File) error {
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
	}
	return nil
}

func (stat *Stat) setSheetHistogramIntInt(excel *xlsx.File, histData map[int]int, shName string) error {
	sh, err := excel.AddSheet(shName)
	if err != nil {
		return err
	}

	for key, value := range histData {
		sh.Cell(key, 0).SetInt(key)
		sh.Cell(key, 1).SetInt(value)
	}
	return nil
}

func (stat *Stat) setSheetHistogramIntFloat(excel *xlsx.File, histData map[int]float64, shName string) error {
	sh, err := excel.AddSheet(shName)
	if err != nil {
		return err
	}

	for key, value := range histData {
		sh.Cell(key, 0).SetInt(key)
		sh.Cell(key, 1).SetFloat(value)
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
	sh.Cell(7, 2).SetString(strconv.Itoa(int(100.0*float64(sum)/float64(len(stat.Games)))) + "%")

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

func (stat *Stat) SaveExcel(path string, saveDistrInExcel bool) error {
	excel := xlsx.NewFile()

	histData := map[int]int{}
	for i := 0; i <= 25; i++ {
		histData[i] = 0
	}

	histDataStep := map[int]int{}
	for i := 0; i < 70; i++ {
		histDataStep[i] = 0
	}

	maxStep := 0
	for i := 0; i < len(stat.Games); i++ {
		histDataStep[stat.Games[i].Step]++
		if stat.Games[i].Step > maxStep {
			maxStep = stat.Games[i].Step
		}
	}

	histDataStepToPoints := map[int]float64{}
	histDataStepToGamesCount := map[int]int{}
	for step, count := range histDataStep {
		if count == 0 && step > maxStep {
			delete(histDataStep, step)
		} else {
			histDataStepToPoints[step] = 0
			histDataStepToGamesCount[step] = 0
		}
	}

	for i := 0; i < len(stat.Games); i++ {
		histDataStepToPoints[stat.Games[i].Step] += float64(stat.Games[i].Points)
		histDataStepToGamesCount[stat.Games[i].Step]++
	}
	for step, _ := range histDataStepToPoints {
		if float64(histDataStepToGamesCount[step]) > 0 {
			histDataStepToPoints[step] /= float64(histDataStepToGamesCount[step])
		}
	}

	for _, g := range stat.Games {
		histData[g.Points] = histData[g.Points] + 1
	}

	if saveDistrInExcel {
		if err := stat.setSheetPoints(excel); err != nil {
			return err
		}
	}

	if err := stat.setSheetHistogramIntInt(excel, histData, "histogram"); err != nil {
		return err
	}

	if err := stat.setSheetHistogramIntInt(excel, histDataStep, "histogramSteps"); err != nil {
		return err
	}

	if err := stat.setSheetHistogramIntFloat(excel, histDataStepToPoints, "histogramStepsToPoints"); err != nil {
		return err
	}

	if err := stat.setSheetInfo(excel); err != nil {
		return err
	}

	return excel.Save(path)
}

func (stat *Stat) SaveToFile(dir, subdir string, saveDistrInExcel bool) error {
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

	err = stat.SaveExcel(path.Join(statDir, "distribution"+subdir+".xlsx"), saveDistrInExcel)
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
