# Email Templates Guide

The application uses a dynamic template system that automatically discovers email templates from the `assets/templates/email/` directory.

## Quick Start

1. **Create MJML Template**:
   ```bash
   touch assets/templates/email/template4.mjml
   ```

2. **Write MJML Content** (see example below)

3. **Convert to HTML**:
   ```bash
   make mjml
   ```

4. **Template is automatically available** - no code changes needed!

## Creating New Email Templates

### Step 1: Create MJML Template

Create a new MJML file in `assets/templates/email/`:

```bash
# Create a new MJML file
touch assets/templates/email/template4.mjml
```

### Step 2: Write MJML Content

```mjml
<mjml>
  <mj-head>
    <mj-title>Template 4</mj-title>
  </mj-head>
  <mj-body>
    <mj-section>
      <mj-column>
        <mj-text font-size="20px" color="#333" font-family="helvetica">
          Hallo {{.RecipientName}},
        </mj-text>
        <mj-text color="#666" font-family="helvetica">
          Hier ist Ihr PDF-Dokument: {{.PDFName}}
        </mj-text>
        <mj-text color="#666" font-family="helvetica">
          Mit freundlichen Grüßen,<br/>
          {{.SenderName}}
        </mj-text>
      </mj-column>
    </mj-section>
  </mj-body>
</mjml>
```

### Step 3: Convert to HTML

```bash
# Convert all MJML templates to HTML
make mjml

# Or convert a specific template
mjml assets/templates/email/template4.mjml -o assets/templates/email/template4.html
```

### Step 4: Test Your Template

1. Start the application: `make dev`
2. Upload a PDF file with email in filename
3. Select your new template from the dropdown
4. Check the sent email in Mailpit: `http://localhost:8025`

## Template Variables

Email templates support the following variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `{{.RecipientName}}` | Name of the email recipient | "John Doe" |
| `{{.SenderName}}` | Name of the email sender | "PDF Emailer" |
| `{{.PDFName}}` | Name of the uploaded PDF file | "document_2025-01-15.pdf" |
| `{{.Date}}` | Current date in YYYY-MM-DD format | "2025-01-15" |

## Template Discovery

The application automatically:

- ✅ Scans `assets/templates/email/` directory for `.html` files
- ✅ Discovers new templates without restart
- ✅ Validates template names during upload
- ✅ Provides helpful error messages with available templates

## Template Validation

When uploading files, the system validates:

- Template name is not empty
- Template file exists in the email templates directory
- Provides clear error messages if template is invalid

### Example Error Messages

```
"template 'nonexistent' not found - available templates: template1, template2, template3"
```

## MJML Conversion

### Installing MJML CLI

```bash
# Install MJML CLI globally
npm install -g mjml

# Verify installation
mjml --version
```

### Converting Templates

```bash
# Convert all MJML templates to HTML
make mjml

# Convert a specific template
mjml assets/templates/email/template4.mjml -o assets/templates/email/template4.html

# Convert with custom options
mjml --config.beautify true --config.minify false template4.mjml -o template4.html
```

### MJML Configuration

Create a `.mjmlrc` file in the project root for custom MJML settings:

```json
{
  "beautify": true,
  "minify": false,
  "validationLevel": "soft"
}
```

## Development Workflow

1. **Create MJML file** in `assets/templates/email/`
2. **Convert to HTML** using `make mjml`
3. **Test template** by uploading a PDF and selecting it
4. **View email** in Mailpit at `http://localhost:8025`
5. **Iterate** and improve the template

## Template Best Practices

### Design Guidelines

- **Use MJML**: MJML provides responsive email templates that work across all email clients
- **Keep it Simple**: Avoid complex layouts that might not render consistently
- **Test Thoroughly**: Always test templates in Mailpit before production use
- **Use Variables**: Leverage the available template variables for personalization
- **Consistent Naming**: Use descriptive template names (e.g., `invoice_template`, `receipt_template`)

### Email Client Compatibility

- **Mobile First**: Design for mobile devices first
- **Fallback Fonts**: Use web-safe fonts as fallbacks
- **Image Handling**: Always provide alt text for images
- **Link Styling**: Use inline styles for links

### Performance Tips

- **Minimize Images**: Use optimized images and consider text alternatives
- **Efficient CSS**: Use MJML's built-in responsive features
- **Test Loading**: Check template loading times in different email clients

## Advanced Templates

### Conditional Content

```mjml
<mj-text>
  {{if .RecipientName}}
    Hallo {{.RecipientName}},
  {{else}}
    Hallo,
  {{end}}
</mj-text>
```

### Custom Styling

```mjml
<mj-text 
  font-size="16px" 
  color="#333" 
  font-family="Arial, sans-serif"
  line-height="24px">
  Custom styled text
</mj-text>
```

### Responsive Design

```mjml
<mj-section>
  <mj-column width="50%">
    <mj-text>Left column</mj-text>
  </mj-column>
  <mj-column width="50%">
    <mj-text>Right column</mj-text>
  </mj-column>
</mj-section>
```

## Troubleshooting

### Common Issues

1. **Template not appearing in dropdown**
   ```bash
   # Check if HTML file exists
   ls -la assets/templates/email/
   
   # Rebuild templates
   make mjml
   
   # Restart application
   make dev
   ```

2. **MJML conversion errors**
   ```bash
   # Check MJML syntax
   mjml --validate template4.mjml
   
   # Check MJML version
   mjml --version
   ```

3. **Email not rendering correctly**
   - Test in Mailpit first
   - Check email client compatibility
   - Use simpler layouts for better compatibility

### Getting Help

- [MJML Documentation](https://mjml.io/documentation/)
- [Email Client Compatibility](https://www.campaignmonitor.com/css/)
- [MJML Playground](https://mjmlio.github.io/mjml-app/)

## Examples

### Simple Template

```mjml
<mjml>
  <mj-head>
    <mj-title>Simple Template</mj-title>
  </mj-head>
  <mj-body>
    <mj-section>
      <mj-column>
        <mj-text font-size="20px" color="#333">
          Hallo {{.RecipientName}},
        </mj-text>
        <mj-text color="#666">
          Hier ist Ihr PDF-Dokument: {{.PDFName}}
        </mj-text>
        <mj-text color="#666">
          Mit freundlichen Grüßen,<br/>
          {{.SenderName}}
        </mj-text>
      </mj-column>
    </mj-section>
  </mj-body>
</mjml>
```

### Professional Template

```mjml
<mjml>
  <mj-head>
    <mj-title>Professional Template</mj-title>
  </mj-head>
  <mj-body background-color="#f4f4f4">
    <mj-section background-color="#ffffff" padding="20px">
      <mj-column>
        <mj-text font-size="24px" color="#333" font-weight="bold">
          PDF Dokument
        </mj-text>
        <mj-divider border-color="#e0e0e0" />
        <mj-text font-size="16px" color="#666" line-height="24px">
          Sehr geehrte/r {{.RecipientName}},
        </mj-text>
        <mj-text font-size="16px" color="#666" line-height="24px">
          anbei erhalten Sie das gewünschte PDF-Dokument: <strong>{{.PDFName}}</strong>
        </mj-text>
        <mj-text font-size="16px" color="#666" line-height="24px">
          Bei Fragen stehen wir Ihnen gerne zur Verfügung.
        </mj-text>
        <mj-divider border-color="#e0e0e0" />
        <mj-text font-size="14px" color="#999">
          Mit freundlichen Grüßen<br/>
          {{.SenderName}}<br/>
          {{.Date}}
        </mj-text>
      </mj-column>
    </mj-section>
  </mj-body>
</mjml>
```
