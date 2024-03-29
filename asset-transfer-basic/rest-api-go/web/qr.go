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

	var cert CertificateInfo
	if err := json.NewDecoder(r.Body).Decode(&cert); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate PDF
	pdfData, err := genPDF(cert)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=certificate.pdf")

	if _, err := w.Write(pdfData); err != nil {
		http.Error(w, "Failed to write PDF data to response", http.StatusInternalServerError)
		return
	}
}

func genPDF(cert CertificateInfo) ([]byte, error) {

	var pdfBuffer bytes.Buffer
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.ImageOptions("./img/template1.png", 0, 0, 297, 210, false, gofpdf.ImageOptions{}, 0, "")

	if cert.BeginDate != "" || cert.EndDate != "" {
		pdf.SetFont("Helvetica", "", 18)
		pdf.Text(70, 195, "Begin_Date: "+cert.BeginDate)

		pdf.SetFont("Helvetica", "", 18)
		pdf.Text(145, 195, "End_Date: "+cert.EndDate)
	}

	pdf.SetFont("Helvetica", "", 18)
	addTextCentered(pdf, "IssuerDate: "+cert.IssuerDate, 180, 18)

	pdf.SetFont("Helvetica", "B", 36)
	addTextCentered(pdf, cert.StudentName, 110, 36)

	pdf.SetFont("Helvetica", "", 20)
	addTextCentered(pdf, cert.Course, 150, 20)

	pdf.SetFont("Helvetica", "", 20)
	addTextCentered(pdf, cert.Transaction, 160, 20)

	pdf.SetFont("Helvetica", "", 15)
	pdf.Text(70, 170, cert.Issuer)
	pdf.Text(190, 170, cert.EndorserName)

	if err := pdf.Output(&pdfBuffer); err != nil {
		return nil, err
	}

	return pdfBuffer.Bytes(), nil
}
