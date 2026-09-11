package main

func FNV1aHash(key string) uint32 {
	var hash uint32 = 2166136261 //offset basis, fixed constant
	for i := 0; i < len(key); i++ {
		hash ^= uint32(key[i])
		hash *= 16777619
	}
	return hash
}

func djb2(key string) uint32 {
	var hash uint32 = 5381
	for i := 0; i < len(key); i++ {
		hash = hash*33 + uint32(key[i])
	}
	return hash
}

func djb2a(key string) uint32 {
	var hash uint32 = 5381
	for i := 0; i < len(key); i++ {
		hash = hash*33 ^ uint32(key[i])
	}
	return hash
}
