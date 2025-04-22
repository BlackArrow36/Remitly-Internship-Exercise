package config

import "database/sql"

type SwiftCode struct {
	SWIFTCode     string         `db:"swift_code"`
	BankName      string         `db:"bank_name"`
	IsHeadquarter bool           `db:"is_headquarter"`
	BranchOf      sql.NullString `db:"branch_of"`
	CountryISO2   string         `db:"country_iso2"`
	CountryName   string         `db:"country_name"`
	City          string         `db:"city"`
	Address       string         `db:"address"`
	TimeZone      string         `db:"time_zone"`
}

type SwiftCodeCountry struct {
	SWIFTCode     string         `db:"swift_code"`
	BankName      string         `db:"bank_name"`
	IsHeadquarter bool           `db:"is_headquarter"`
	BranchOf      sql.NullString `db:"branch_of"`
	CountryISO2   string         `db:"country_iso2"`
	CountryName   string         `db:"country_name"`
	Address       string         `db:"address"`
}
