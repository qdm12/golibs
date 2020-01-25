package crypto

func (c *crypto) NewSalt() (salt [saltSize]byte, err error) {
	saltSlice, err := c.random.GenerateRandomBytes(32)
	if err != nil {
		return salt, err
	}
	copy(salt[:], saltSlice)
	return salt, nil
}
