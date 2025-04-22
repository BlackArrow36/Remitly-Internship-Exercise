package controller

import (
	"database/sql"
	"net/http"
	"strings"
	"swift-api/config"

	"github.com/gin-gonic/gin"
)

type SwiftCodeResponse struct {
	Address       string                 `json:"address"`
	BankName      string                 `json:"bankName"`
	CountryISO2   string                 `json:"countryISO2"`
	CountryName   string                 `json:"countryName"`
	IsHeadquarter bool                   `json:"isHeadquarter"`
	SwiftCode     string                 `json:"swiftCode"`
	Branches      []SwiftCodeSubResponse `json:"branches"`
}
type SwiftCodeSubResponse struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName,omitempty"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type CountryResponse struct {
	CountryISO2 string                 `json:"countryISO2"`
	CountryName string                 `json:"countryName"`
	SwiftCodes  []SwiftCodeSubResponse `json:"swiftCodes"`
}

func GetSwiftCode(c *gin.Context) {
	code := c.Param("swift-code")
	isHeadquarter := strings.HasSuffix(code, "XXX")
	if isHeadquarter {
		var swift SwiftCodeResponse

		err := config.DB.QueryRow(`
			SELECT swift_code, bank_name, address, country_iso2, country_name, is_headquarter 
			FROM swift_codes 
			WHERE swift_code = $1`, code).
			Scan(&swift.SwiftCode, &swift.BankName, &swift.Address, &swift.CountryISO2, &swift.CountryName, &swift.IsHeadquarter)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		swift.Branches = []SwiftCodeSubResponse{}
		rows, err := config.DB.Query(`
			SELECT swift_code, bank_name, address, country_iso2, country_name, is_headquarter 
			FROM swift_codes 
			WHERE branch_of = $1`, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var b SwiftCodeSubResponse
			if err := rows.Scan(&b.SwiftCode, &b.BankName, &b.Address, &b.CountryISO2, &b.IsHeadquarter); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			swift.Branches = append(swift.Branches, b)
		}
		c.JSON(http.StatusOK, swift)
	} else {
		var swift SwiftCodeSubResponse

		err := config.DB.QueryRow(`
			SELECT swift_code, bank_name, address, country_iso2, country_name, is_headquarter 
			FROM swift_codes 
			WHERE swift_code = $1`, code).
			Scan(&swift.SwiftCode, &swift.BankName, &swift.Address, &swift.CountryISO2, &swift.CountryName, &swift.IsHeadquarter)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, swift)
	}

}

func GetByCountry(c *gin.Context) {
	iso := strings.ToUpper(c.Param("countryISO2"))
	rows, err := config.DB.Query(`
		SELECT swift_code, bank_name, address, country_iso2, country_name, is_headquarter 
		FROM swift_codes 
		WHERE country_iso2 = $1`, iso)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result CountryResponse
	var swiftCodes []SwiftCodeSubResponse

	for rows.Next() {
		var sc SwiftCodeSubResponse
		err := rows.Scan(&sc.SwiftCode, &sc.BankName, &sc.Address, &sc.CountryISO2, &result.CountryName, &sc.IsHeadquarter)
		if err == nil {
			swiftCodes = append(swiftCodes, sc)
		}
	}

	if len(swiftCodes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No SWIFT codes found for country " + iso})
		return
	}

	result.CountryISO2 = iso
	result.SwiftCodes = swiftCodes
	c.JSON(http.StatusOK, result)
}

func AddSwiftCode(c *gin.Context) {
	var input config.SwiftCodeCountry

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO swift_codes 
		(swift_code, bank_name, is_headquarter, branch_of, country_iso2, country_name, city, address, time_zone)
		VALUES ($1, $2, $3, $4, $5, $6, NULL, $7, NULL)
	`

	var branchOf sql.NullString
	if !input.IsHeadquarter {
		branchOf = sql.NullString{String: input.SWIFTCode[:8] + "XXX", Valid: true}
	}

	_, err := config.DB.Exec(query,
		input.SWIFTCode,
		input.BankName,
		input.IsHeadquarter,
		branchOf,
		input.CountryISO2,
		input.CountryName,
		input.Address,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SWIFT code added successfully"})
}

func DeleteSwiftCode(c *gin.Context) {
	code := c.Param("swift-code")

	result, err := config.DB.Exec("DELETE FROM swift_codes WHERE swift_code = $1", code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SWIFT code deleted successfully"})
}
