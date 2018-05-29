package main

type IntStrPair struct {
	Key int
	Val string
}

func IntStrPairsHasKey(pairs []IntStrPair, key int) bool {
	for _, pair := range pairs {
		if pair.Key == key {
			return true
		}
	}
	return false
}

type StrStrPair struct {
	Key string
	Val string
}

func StrStrPairsHasKey(pairs []StrStrPair, key string) bool {
	for _, pair := range pairs {
		if pair.Key == key {
			return true
		}
	}
	return false
}
