package main

// DoiInfo contains information associated with a DOI, most notably
// the available open access manuscripts
type DoiInfo struct {
	Manuscripts []Manuscript `json:"manuscripts"`
}

// Manuscript describes an open access manuscript that can be
// selected by the user.
type Manuscript struct {
	Description string `json:"description"` // Description of this manuscript
	Location    string `json:"location"`    // Location URI of manuascript (e.g. pdf)
}
