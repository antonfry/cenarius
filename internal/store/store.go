package store

type Store interface {
	LoginWithPassword() LoginWithPasswordRepository
	CreditCard() CreditCardRepository
	SecretText() SecretTextRepository
	SecretBinary() SecretBinaryRepository
	// SecretData() SecretDataRepository
	User() UserRepository
	Close()
}
