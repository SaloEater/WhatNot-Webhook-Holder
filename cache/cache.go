package cache

import "fmt"

type Cache[T any] interface {
	Get(string) (T, bool)
	Set(string, T)
	Has(string) bool
	Delete(string)
}

func IdToKey(id int64) string {
	return fmt.Sprintf("%d", id)
}
