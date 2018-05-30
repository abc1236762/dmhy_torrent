package main

type UintStrPair struct {
	Key uint
	Val string
}

func UintStrPairsHasKey(pairs []UintStrPair, key uint) bool {
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
