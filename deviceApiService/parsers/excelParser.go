package parsers

import (
	"encoding/json"

	"github.com/tealeg/xlsx"
	"proj.com/apisvc/db/models"
)

func ParseProfileExcel(filePath string, sheetNum int) ([]models.DeviceConfigProfile, error) {
	excelFileName := filePath

	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		return []models.DeviceConfigProfile{}, err
	}

	sheet := xlFile.Sheets[sheetNum]
	var profileList []models.DeviceConfigProfile
	for rowIndex, row := range sheet.Rows {
		if rowIndex == 0 {
			continue
		}

		var profile models.DeviceConfigProfile
		var profileData map[string]string
		for cellIndex, cell := range row.Cells {
			if cellIndex == 0 {
				profile.ID = cell.String()
			} else {
				key := sheet.Rows[0].Cells[cellIndex].String()
				value := cell.String()

				if profileData == nil {
					profileData = make(map[string]string)
				}
				profileData[key] = value
			}
		}
		jsonString, err := json.Marshal(profileData)
		if err != nil {
			return []models.DeviceConfigProfile{}, err
		}
		profile.Data = string(jsonString)
		profileList = append(profileList, profile)
	}

	return profileList, nil
}
