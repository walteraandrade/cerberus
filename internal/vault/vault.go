package vault

import "time"

type Category struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	ParentID string     `json:"parent_id,omitempty"`
	Children []Category `json:"children,omitempty"`
}

type Entry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url,omitempty"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Notes     string    `json:"notes,omitempty"`
	Category  string    `json:"category,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Vault struct {
	Version    int        `json:"version"`
	Categories []Category `json:"categories"`
	Entries    []Entry    `json:"entries"`
}

func New() *Vault {
	return &Vault{
		Version:    1,
		Categories: []Category{},
		Entries:    []Entry{},
	}
}

func (v *Vault) AddEntry(e Entry) {
	v.Entries = append(v.Entries, e)
}

func (v *Vault) RemoveEntry(id string) bool {
	for i, e := range v.Entries {
		if e.ID == id {
			v.Entries = append(v.Entries[:i], v.Entries[i+1:]...)
			return true
		}
	}
	return false
}

func (v *Vault) FindEntry(id string) *Entry {
	for i := range v.Entries {
		if v.Entries[i].ID == id {
			return &v.Entries[i]
		}
	}
	return nil
}

func (v *Vault) AddCategory(c Category) {
	v.Categories = append(v.Categories, c)
}

func (v *Vault) EntriesByCategory(categoryID string) []Entry {
	var result []Entry
	for _, e := range v.Entries {
		if e.Category == categoryID {
			result = append(result, e)
		}
	}
	return result
}
