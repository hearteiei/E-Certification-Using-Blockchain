package web

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/jung-kurt/gofpdf"
)

type CertificateInfos struct {
	StudentName  string `json:"studentName"`
	Course       string `json:"course"`
	Issuer       string `json:"issuer"`
	EndorserName string `json:"endorserName"`
	BeginDate    string `json:"beginDate"`
	EndDate      string `json:"endDate"`
	Mail         string `json:"Mail"`
	IssuerDate   string `json:"issuerdate"`
}

func GenCertificates(w http.ResponseWriter, r *http.Request) {
	// Parse JSON from request body
	var cert CertificateInfo
	if err := json.NewDecoder(r.Body).Decode(&cert); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate PDF content
	pdfData, err := genPDF(cert)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=certificate.pdf")

	// Write PDF data to response body
	if _, err := w.Write(pdfData); err != nil {
		http.Error(w, "Failed to write PDF data to response", http.StatusInternalServerError)
		return
	}
}

func genPDF(cert CertificateInfo) ([]byte, error) {
	// Generate PDF content using gofpdf library
	var pdfBuffer bytes.Buffer
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Add your PDF generation logic here
	// Add background image
	pdf.ImageOptions("./img/template1.png", 0, 0, 297, 210, false, gofpdf.ImageOptions{}, 0, "")

	// Add recipient name
	if cert.BeginDate != "" || cert.EndDate != "" {
		pdf.SetFont("Helvetica", "", 18)
		pdf.Text(70, 195, "Begin_Date: "+cert.BeginDate)

		pdf.SetFont("Helvetica", "", 18)
		pdf.Text(145, 195, "End_Date: "+cert.EndDate)
	}

	pdf.SetFont("Helvetica", "", 18)
	addTextCentered(pdf, "IssuerDate: "+cert.IssuerDate, 180, 18)

	// Add recipient name centered horizontally
	pdf.SetFont("Helvetica", "B", 36)
	addTextCentered(pdf, cert.StudentName, 110, 36)

	// Add course name centered horizontally
	pdf.SetFont("Helvetica", "", 20)
	addTextCentered(pdf, cert.Course, 150, 20)

	// Add transaction details centered horizontally
	pdf.SetFont("Helvetica", "", 20)
	addTextCentered(pdf, cert.Transaction, 160, 20)

	// Add issuer and endorser names centered horizontally
	pdf.SetFont("Helvetica", "", 15)
	pdf.Text(70, 170, cert.Issuer)
	pdf.Text(190, 170, cert.EndorserName)

	// Output PDF content to buffer
	if err := pdf.Output(&pdfBuffer); err != nil {
		return nil, err
	}

	return pdfBuffer.Bytes(), nil
}
