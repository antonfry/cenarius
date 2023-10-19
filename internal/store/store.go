package store

type Store interface {
	LoginWithPassword() LoginWithPasswordRepository
	CreditCard() CreditCardRepository
	SecretText() SecretTextRepository
	SecretFile() SecretFileRepository
	// SecretData() SecretDataRepository
	User() UserRepository
	Close()
}
