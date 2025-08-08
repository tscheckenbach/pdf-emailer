# Email Templates Directory

This directory contains email templates for the PDF Emailer application.

## Structure

- `template1.mjml` - Template 1 for specific use cases
- `template2.mjml` - Template 2 for specific use cases
- `template3.mjml` - Template 3 for specific use cases

## MJML Conversion

To convert MJML templates to HTML, use the provided scripts:

```bash
# Install MJML CLI (if not already installed)
./scripts/install-mjml.sh

# Convert all templates using Makefile
make mjml

# Or convert manually
cd assets/templates/email
for file in *.mjml; do
    mjml "$file" -o "${file%.mjml}.html"
done
```

## Usage

Email templates will be used to send PDF files to recipients with custom formatting and messaging.

## Template Variables

Templates can use the following variables:
- `{{.RecipientName}}` - Name of the email recipient
- `{{.SenderName}}` - Name of the sender
- `{{.Message}}` - Custom message content
- `{{.PDFName}}` - Name of the attached PDF file
- `{{.Date}}` - Current date
