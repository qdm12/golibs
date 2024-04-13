package crypto

func (c *Crypto) NewSalt() (salt [saltSize]byte, err error) {
	const nBytes = 32
	saltSlice, err := c.random.GenerateRandomBytes(nBytes)
	if err != nil {
		return salt, err
	}
	copy(salt[:], saltSlice)
	return salt, nil
}
