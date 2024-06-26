package web

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"net/url"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

type CertificateInfo struct {
	StudentName  string `json:"studentName"`
	Course       string `json:"course"`
	Issuer       string `json:"issuer"`
	EndorserName string `json:"endorserName"`
	BeginDate    string `json:"beginDate"`
	EndDate      string `json:"endDate"`
	Mail         string `json:"Mail"`
	Transaction  string `json:"transaction"`
	IssuerDate   string `json:"issuerdate"`
}

func GenerateCertificates(w http.ResponseWriter, r *http.Request) {
	// Parse JSON from request body
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

	// Send email
	if err := sendEmailWithAttachment(pdfData, cert); err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	// return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "PDF generated and emailed successfully"})
}

func calculateCenterPosition(pdf *gofpdf.Fpdf, text string, fontSize float64) float64 {
	width := pdf.GetStringWidth(text)
	pageWidth, _ := pdf.GetPageSize()
	return (pageWidth - width) / 2
}

func addTextCentered(pdf *gofpdf.Fpdf, text string, y float64, fontSize float64) {
	x := calculateCenterPosition(pdf, text, fontSize)
	pdf.Text(x, y, text)
}

func generatePDF(cert CertificateInfo) ([]byte, error) {
	values := url.Values{}
	values.Set("studentName", cert.StudentName)
	values.Set("course", cert.Course)
	values.Set("issuer", cert.Issuer)
	values.Set("endorserName", cert.EndorserName)
	values.Set("beginDate", cert.BeginDate)
	values.Set("endDate", cert.EndDate)
	values.Set("mail", cert.Mail)
	values.Set("IssuerDate", cert.IssuerDate)

	queryString := values.Encode()

	var pdfBuffer bytes.Buffer
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.ImageOptions("./img/template1.png", 0, 0, 297, 210, false, gofpdf.ImageOptions{}, 0, "")

	qrLink := fmt.Sprintf("http://localhost:3000/fetch?%s", queryString)
	qrImg, err := qrcode.Encode(qrLink, qrcode.Medium, 256)
	if err != nil {
		log.Fatal("Error generating QR code: ", err)
	}

	qrImgReader := bytes.NewReader(qrImg)

	pdf.RegisterImageReader("qr_code", "png", qrImgReader)
	pdf.Image("qr_code", 10, 160, 40, 40, false, "", 0, "")

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

	if cert.Transaction != "" {
		pdf.SetFont("Helvetica", "", 20)
		addTextCentered(pdf, cert.Transaction, 160, 20)
	}

	pdf.SetFont("Helvetica", "", 15)
	pdf.Text(70, 170, cert.Issuer)
	pdf.Text(190, 170, cert.EndorserName)

	if err := pdf.Output(&pdfBuffer); err != nil {
		return nil, err
	}

	return pdfBuffer.Bytes(), nil
}

func sendEmailWithAttachment(pdfData []byte, cert CertificateInfo) error {

	from := "test259492@gmail.com"
	password := ""

	to := cert.Mail

	// Email configuration
	subject := "Certificate PDF: " + cert.Course
	body := fmt.Sprintf("Dear %s,\n\nPlease find attached your certificate for the course: %s.\n\nRegards,\n%s", cert.StudentName, cert.Course, cert.Issuer)

	encodedPDF := base64.StdEncoding.EncodeToString(pdfData)

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
		"Content-Type: application/pdf\r\n" +
		"Content-Disposition: attachment; filename=\"certificate.pdf\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" +
		encodedPDF + "\r\n" +
		"--boundary123456--\r\n"

	// Send email
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
