package user

// Info describes a user that has been authenticated to the system.
type Info interface {
	// GetName returns the name that uniquely identifies this user among all
	// other active users.
	GetName() string
	// GetUID returns a unique value for a particular user that will change
	// if the user is removed from the system and another user is added with
	// the same name.
	GetUID() string
	// GetGroups returns the names of the groups the user is a member of
	GetGroups() []string

	GetExtra() map[string][]string
}

// DefaultInfo provides a simple user information exchange object
// for components that implement the UserInfo interface.
type DefaultInfo struct {
	Name   string
	UID    string
	Groups []string
	Extra  map[string][]string
}

func (i *DefaultInfo) GetName() string {
	return i.Name
}

func (i *DefaultInfo) GetUID() string {
	return i.UID
}

func (i *DefaultInfo) GetGroups() []string {
	return i.Groups
}

func (i *DefaultInfo) GetExtra() map[string][]string {
	return i.Extra
}
