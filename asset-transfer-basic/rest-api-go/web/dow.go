package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func GenerateCertificatesanddowload(w http.ResponseWriter, r *http.Request) {

	var cert CertificateInfo

	if err := json.NewDecoder(r.Body).Decode(&cert); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate PDF
	pdfData, err := generatePDF(cert)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=certificate.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfData)))

	if _, err := io.Copy(w, bytes.NewReader(pdfData)); err != nil {
		http.Error(w, "Failed to write PDF to response", http.StatusInternalServerError)
		return
	}
}
