package dao

func SelectStudentByUnifiedCode(unifiedCode string) error {
	_, err := Get(unifiedCode)
	if err != nil {
		return err
	}
	return nil
}
