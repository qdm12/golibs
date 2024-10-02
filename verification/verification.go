package verification

type Verifier struct {
	*Regex
}

func NewVerifier() *Verifier {
	return &Verifier{
		Regex: NewRegex(),
	}
}
