package database

type DBLayer interface {
	GetCallInstanceDomainName(callID string) (domainName string, err error)
	GetLeastUsedInstanceDomainName() (domainName string, err error)

	NewInstance(domainName string) error

	AddConnection(callID, domainName string) error
	DeleteConnection(callID, domainName string) error
}
