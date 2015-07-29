package main

func retry(f func() error) error {
	var err error
	for {
		err = f()
		if err == nil {
			break
		}
	}
	return err
}
