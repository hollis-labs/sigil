package config

// ApplyDefaults fills in default values for optional fields.
func ApplyDefaults(page *Page) {
	if page.Sigil == "" {
		page.Sigil = "1.0"
	}
	if page.Kind == "" {
		page.Kind = "page"
	}
	if page.Overlay == "" {
		page.Overlay = "page"
	}
	applyComponentDefaults(&page.Layout)
}

func applyComponentDefaults(c *Component) {
	for i := range c.Children {
		applyComponentDefaults(&c.Children[i])
	}
}
