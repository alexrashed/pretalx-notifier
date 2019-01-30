package main

// submissions page json
type submissions struct {
	Count    int          `json:"count"`
	Next     string       `json:"next"`
	Previous interface{}  `json:"previous"`
	Results  []submission `json:"results"`
}

// submission json
type submission struct {
	Abstract      string        `json:"abstract"`
	Answers       []interface{} `json:"answers"`
	Code          string        `json:"code"`
	ContentLocale string        `json:"content_locale"`
	Description   string        `json:"description"`
	DoNotRecord   bool          `json:"do_not_record"`
	Duration      string        `json:"duration"`
	//Image         interface{}   `json:"image"`
	//IsFeatured    bool          `json:"is_featured"`
	//Slot          interface{}   `json:"slot"`
	Speakers      []struct {
		//	Avatar    interface{} `json:"avatar"`
		Biography string      `json:"biography"`
		Code      string      `json:"code"`
		Name      string      `json:"name"`
	} `json:"speakers"`
	State          string         `json:"state"`
	SubmissionType submissionType `json:"submission_type"`
	Title string                  `json:"title"`
	//Track interface{} `json:"track"`
}

// submission type names json
type submissionType struct {
	De string `json:"de"`
	En string `json:"en"`
}