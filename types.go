package main

// https://mholt.github.io/json-to-go/

type DiscussionComments struct {
	Data struct {
		Repository struct {
			Discussion struct {
				Body     string `json:"body"`
				Comments struct {
					Edges []struct {
						Node struct {
							Author struct {
								Login string `json:"login"`
							} `json:"author"`
							Body string `json:"body"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"comments"`
			} `json:"discussion"`
		} `json:"repository"`
	} `json:"data"`
}

type Login struct {
	Data struct {
		Viewer struct {
			User string `json:"login"`
		} `json:"viewer"`
	} `json:"data"`
}
