package models

type MemberApplication struct {
	ID                   int64   `json:"id"`
	Status               string  `json:"status"`
	NumberOfShares       int64   `json:"numberOfShares"`
	PricePerShare        float64 `json:"pricePerShare"`
	PostalStreet         string  `json:"postalStreet"`
	PostalSuburb         string  `json:"postalSuburb"`
	PostalState          string  `json:"postalState"`
	PostalCountry        string  `json:"postalCountry"`
	PostalPostcode       string  `json:"postalPostcode"`
	ChessHinSrn          string  `json:"chessHinSrn"`
	Phone                string  `json:"phone"`
	Email                string  `json:"email"`
	Applicant1Investor   string  `json:"applicant1Investor"`
	Applicant1TrustName  string  `json:"applicant1TrustName"`
	Applicant1NumberType string  `json:"applicant1NumberType"`
	Applicant1Number     string  `json:"applicant1Number"`
	Applicant1Surname    string  `json:"applicant1Surname"`
	ReferenceNumber      string  `json:"referenceNumber"`
	ManagementFee        float64 `json:"managementFee"`
	ProcessingFee        string  `json:"processingFee"`
	Title                string  `json:"title"`
	Name                 string  `json:"name"`
	Surname              string  `json:"surname"`
	ContactPerson        string  `json:"contactPerson"`
	CompanyName          string  `json:"companyName"`
}