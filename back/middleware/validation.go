package middleware

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// ValidationRule represents a validation rule
type ValidationRule struct {
	Field    string
	Required bool
	MinLen   int
	MaxLen   int
	Pattern  string
	Type     string // string, int, float, email, url
	Custom   func(interface{}) error
}

// ValidateJSON validates JSON request body against provided rules
func ValidateJSON(rules []ValidationRule) fiber.Handler {
	return func(c fiber.Ctx) error {
		var body map[string]interface{}

		if err := c.Bind().JSON(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid JSON format",
				"details": err.Error(),
			})
		}

		errors := make(map[string]string)

		for _, rule := range rules {
			value, exists := body[rule.Field]

			// Check required fields
			if rule.Required && (!exists || value == nil || value == "") {
				errors[rule.Field] = "This field is required"
				continue
			}

			// Skip validation if field is not present and not required
			if !exists || value == nil {
				continue
			}

			// Validate based on type and rules
			if err := validateField(rule, value); err != nil {
				errors[rule.Field] = err.Error()
			}
		}

		if len(errors) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation failed",
				"details": errors,
			})
		}

		// Store validated body for use in handlers
		c.Locals("validatedBody", body)
		return c.Next()
	}
}

// validateField validates a single field against its rule
func validateField(rule ValidationRule, value interface{}) error {
	// Convert value to string for length and pattern validation
	strValue := fmt.Sprintf("%v", value)

	// Type validation
	switch rule.Type {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("must be a string")
		}
	case "int":
		if _, err := strconv.Atoi(strValue); err != nil {
			return fmt.Errorf("must be a valid integer")
		}
	case "float":
		if _, err := strconv.ParseFloat(strValue, 64); err != nil {
			return fmt.Errorf("must be a valid number")
		}
	case "email":
		if !isValidEmail(strValue) {
			return fmt.Errorf("must be a valid email address")
		}
	case "url":
		if !isValidURL(strValue) {
			return fmt.Errorf("must be a valid URL")
		}
	}

	// Length validation
	if rule.MinLen > 0 && len(strValue) < rule.MinLen {
		return fmt.Errorf("must be at least %d characters long", rule.MinLen)
	}
	if rule.MaxLen > 0 && len(strValue) > rule.MaxLen {
		return fmt.Errorf("must be no more than %d characters long", rule.MaxLen)
	}

	// Pattern validation
	if rule.Pattern != "" {
		matched, err := regexp.MatchString(rule.Pattern, strValue)
		if err != nil {
			return fmt.Errorf("invalid pattern validation")
		}
		if !matched {
			return fmt.Errorf("does not match required pattern")
		}
	}

	// Custom validation
	if rule.Custom != nil {
		if err := rule.Custom(value); err != nil {
			return err
		}
	}

	return nil
}

// SanitizeInput sanitizes input to prevent XSS and injection attacks
func SanitizeInput() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Sanitize query parameters
		queryArgs := c.Request().URI().QueryArgs()
		for key, values := range c.Queries() {
			sanitizedValues := make([]string, len(values))
			for i, value := range values {
				sanitizedValues[i] = sanitizeString(string(value))
			}
			queryArgs.Set(key, strings.Join(sanitizedValues, ","))
		}

		// Sanitize form data
		if c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "PATCH" {
			contentType := c.Get("Content-Type")
			if strings.Contains(contentType, "application/json") {
				// For JSON, we'll sanitize in the handler after parsing
				return c.Next()
			} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
				// Sanitize form values
				c.Request().PostArgs().VisitAll(func(key, value []byte) {
					sanitized := sanitizeString(string(value))
					c.Request().PostArgs().Set(string(key), sanitized)
				})
			}
		}

		return c.Next()
	}
}

// sanitizeString removes potentially dangerous characters
func sanitizeString(input string) string {
	// Remove script tags
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	input = scriptRegex.ReplaceAllString(input, "")

	// Remove other dangerous HTML tags
	dangerousTags := []string{"iframe", "object", "embed", "link", "meta", "style"}
	for _, tag := range dangerousTags {
		tagRegex := regexp.MustCompile(fmt.Sprintf(`(?i)<%s[^>]*>.*?</%s>`, tag, tag))
		input = tagRegex.ReplaceAllString(input, "")
	}

	// Remove javascript: and data: URLs
	jsRegex := regexp.MustCompile(`(?i)javascript:`)
	input = jsRegex.ReplaceAllString(input, "")

	dataRegex := regexp.MustCompile(`(?i)data:`)
	input = dataRegex.ReplaceAllString(input, "")

	// Remove SQL injection patterns
	sqlPatterns := []string{
		`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)`,
		`(?i)(or|and)\s+\d+\s*=\s*\d+`,
		`(?i)'\s*(or|and)\s+'`,
	}
	for _, pattern := range sqlPatterns {
		sqlRegex := regexp.MustCompile(pattern)
		input = sqlRegex.ReplaceAllString(input, "")
	}

	return strings.TrimSpace(input)
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// isValidURL validates URL format
func isValidURL(url string) bool {
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return urlRegex.MatchString(url)
}

// ValidateStruct validates a struct using reflection and tags
func ValidateStruct(s interface{}) map[string]string {
	errors := make(map[string]string)
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	// Handle pointer to struct
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		errors["_general"] = "Invalid data structure"
		return errors
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get validation tags
		required := fieldType.Tag.Get("required") == "true"
		minLen := fieldType.Tag.Get("min_len")
		maxLen := fieldType.Tag.Get("max_len")
		pattern := fieldType.Tag.Get("pattern")

		fieldName := strings.ToLower(fieldType.Name)

		// Check required
		if required && isEmptyValue(field) {
			errors[fieldName] = "This field is required"
			continue
		}

		// Skip validation if field is empty and not required
		if isEmptyValue(field) {
			continue
		}

		// Length validation for strings
		if field.Kind() == reflect.String {
			strVal := field.String()

			if minLen != "" {
				if min, err := strconv.Atoi(minLen); err == nil && len(strVal) < min {
					errors[fieldName] = fmt.Sprintf("Must be at least %d characters long", min)
					continue
				}
			}

			if maxLen != "" {
				if max, err := strconv.Atoi(maxLen); err == nil && len(strVal) > max {
					errors[fieldName] = fmt.Sprintf("Must be no more than %d characters long", max)
					continue
				}
			}

			if pattern != "" {
				if matched, err := regexp.MatchString(pattern, strVal); err == nil && !matched {
					errors[fieldName] = "Does not match required pattern"
					continue
				}
			}
		}

		// Email validation
		if fieldType.Tag.Get("email") == "true" && field.Kind() == reflect.String {
			if !isValidEmail(field.String()) {
				errors[fieldName] = "Must be a valid email address"
			}
		}
	}

	return errors
}

// isEmptyValue checks if a reflect.Value is empty
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	}
	return false
}

// GetValidatedBody retrieves the validated body from context
func GetValidatedBody(c fiber.Ctx) map[string]interface{} {
	if body, ok := c.Locals("validatedBody").(map[string]interface{}); ok {
		return body
	}
	return nil
}
