package importer

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"swift-api/config" // <- adjust to your actual module name
)

type SwiftCode struct {
	SWIFTCode     string
	BankName      string
	IsHeadquarter bool
	BranchOf      sql.NullString
	CountryISO2   string
	CountryName   string
	City          string
	Address       string
	TimeZone      string
}

func ParseCSV(filename string) ([]SwiftCode, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var swiftCodes []SwiftCode
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}

		swift := record[1]
		isHQ := strings.HasSuffix(swift, "XXX")
		branchOf := sql.NullString{Valid: false}
		if !isHQ {
			branchOf = sql.NullString{String: swift[:8] + "XXX", Valid: true}
		}

		swiftCodes = append(swiftCodes, SwiftCode{
			SWIFTCode:     swift,
			BankName:      record[3],
			IsHeadquarter: isHQ,
			BranchOf:      branchOf,
			CountryISO2:   strings.ToUpper(record[0]),
			CountryName:   strings.ToUpper(record[6]),
			City:          record[5],
			Address:       record[4],
			TimeZone:      record[7],
		})
	}
	return swiftCodes, nil
}

func InsertSwiftCodes(data []SwiftCode) error {
	query := `
		INSERT INTO swift_codes 
		(swift_code, bank_name, is_headquarter, branch_of, country_iso2, country_name, city, address, time_zone)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	for _, code := range data {
		_, err := config.DB.Exec(query,
			code.SWIFTCode,
			code.BankName,
			code.IsHeadquarter,
			code.BranchOf,
			code.CountryISO2,
			code.CountryName,
			code.City,
			code.Address,
			code.TimeZone,
		)
		if err != nil {
			return fmt.Errorf("insert failed for %s: %w", code.SWIFTCode, err)
		}
	}
	return nil
}
