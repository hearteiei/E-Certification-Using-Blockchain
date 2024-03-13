package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

type CertificateInfo struct {
	StudentName  string `json:"studentName"`
	Course       string `json:"course"`
	Issuer       string `json:"issuer"`
	EndorserName string `json:"endorserName"`
	BeginDate    string `json:"beginDate"`
	EndDate      string `json:"endDate"`
}

func GenerateCertificates(w http.ResponseWriter, r *http.Request) {
	// Parse JSON from request body
	var certinfos CertificateInfo
	if err := json.NewDecoder(r.Body).Decode(&certinfos); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Temporary directory to store PDFs
	tempDir := "./temp/"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		http.Error(w, "Failed to create temporary directory", http.StatusInternalServerError)
		return
	}

	// Generate PDF and collect file path
	pdfPath := tempDir + certinfos.StudentName + "-" + certinfos.Course + ".pdf"
	if err := generatePDF(certinfos, pdfPath); err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	// Send email with PDF attachment
	if err := sendEmailWithAttachment(pdfPath, certinfos); err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "PDF generated and emailed successfully"})
}

func generatePDF(cert CertificateInfo, filePath string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Add background image
	pdf.ImageOptions("./img/template1.png", 0, 0, 297, 210, false, gofpdf.ImageOptions{}, 0, "")

	// Add recipient name
	pdf.SetFont("Helvetica", "B", 36)
	pdf.Text(145, 110, cert.StudentName)

	// Add course name
	pdf.SetFont("Helvetica", "", 20)
	pdf.Text(145, 150, cert.Course)

	pdf.SetFont("Helvetica", "", 15)
	pdf.Text(88, 170, cert.Issuer)
	pdf.Text(208, 170, cert.EndorserName)

	// Save PDF to file
	if err := pdf.OutputFileAndClose(filePath); err != nil {
		return err
	}
	return nil
}

func sendEmailWithAttachment(pdfPath string, cert CertificateInfo) error {
	// Sender email credentials
	from := "test259492@gmail.com"
	password := "rijq jocq csnq lhmq" // Use an App Password if using Gmail

	// Recipient email address
	to := "uncles1512@gmail.com"

	// Email configuration
	subject := "Certificate PDF: " + cert.Course
	body := fmt.Sprintf("Dear %s,\n\nPlease find attached your certificate for the course: %s.\n\nRegards,\n%s", cert.StudentName, cert.Course, cert.Issuer)

	// SMTP server configuration
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication with SMTP server
	auth := smtp.PlainAuth("", from, password, smtpServer)

	// Open PDF file
	pdfFile, err := os.Open(pdfPath)
	if err != nil {
		return err
	}
	defer pdfFile.Close()

	// Get file info
	fileInfo, err := pdfFile.Stat()
	if err != nil {
		return err
	}

	// Read PDF file content
	pdfData := make([]byte, fileInfo.Size())
	_, err = pdfFile.Read(pdfData)
	if err != nil {
		return err
	}

	// Encode PDF data as base64
	encodedPDF := base64.StdEncoding.EncodeToString(pdfData)

	// Compose email message
	message := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: multipart/mixed; boundary=boundary123456\r\n" +
		"\r\n" +
		"--boundary123456\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		body + "\r\n" +
		"\r\n" +
		"--boundary123456\r\n" +
		"Content-Type: application/pdf; name=\"" + filepath.Base(pdfPath) + "\"\r\n" +
		"Content-Disposition: attachment; filename=\"" + filepath.Base(pdfPath) + "\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" +
		encodedPDF + "\r\n" +
		"--boundary123456--\r\n"

	// Send email
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
