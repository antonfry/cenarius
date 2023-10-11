package store

type Store interface {
	SecretLP() LoginWithPasswordRepository
	SecretCard() CreditCardRepository
	SecretText() SecretTextRepository
	SecretBinary() SecretBinaryRepository
	SecretData() SecretDataRepository
	User() UserRepository
	Close()
}
