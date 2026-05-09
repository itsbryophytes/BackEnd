package dto

type ConfirmDocumentRequest struct {
	ReportDate  string                 `json:"report_date"`
	LabName     string                 `json:"lab_name"`
	Metrics     map[string]interface{} `json:"metrics"`
}