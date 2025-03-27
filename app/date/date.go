package date

import "time"

func FormatDateString(dateString string) string {
	layout := "2006-01-02T15:04:05-07:00"
	t, err := time.Parse(layout, dateString)
	if err != nil {
		return dateString
	}

	outputFormat := "02.01.2006 15:04:05"
	formattedDate := t.Format(outputFormat)

	return formattedDate
}
