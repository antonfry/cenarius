package store

type Store interface {
	LoginWithPassword() LoginWithPasswordRepository
	CreditCard() CreditCardRepository
	SecretText() SecretTextRepository
	SecretFile() SecretFileRepository
	User() UserRepository
	Close()
}
