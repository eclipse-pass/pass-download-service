package main

// DoiInfo contains information associated with a DOI, most notably
// the available open access manuscripts
type DoiInfo struct {
	Manuscripts []Manuscript `json:"manuscripts"`
}

// Manuscript describes an open access manuscript that can be
// selected by the user.
type Manuscript struct {
	Location              string `json:"url"`             // Location URI of manuascript (e.g. pdf)
	RepositoryInstitution string `json:"repositoryLabel"` // Readable label for the repository where the article can be found
	Type                  string `json:"type"`            // The MIME type of the manuscript file
	Source                string `json:"source"`          // The API where we found the file
	Name                  string `json:"name"`            // The file name
}
