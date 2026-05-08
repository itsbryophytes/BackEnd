package dto

type ConfirmDocumentRequest struct {
	PatientName string                 `json:"patient_name"`
	ReportDate  string                 `json:"report_date"`
	LabName     string                 `json:"lab_name"`
	Metrics     map[string]interface{} `json:"metrics"`
}