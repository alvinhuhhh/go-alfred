package util

func Remove(a []string, remove string) []string {
	n := 0
	for _, x := range a {
		if x != remove {
			a[n] = x
			n++
		}
	}
	a = a[:n]
	return a
}