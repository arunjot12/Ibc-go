package wasm

import (
	"math"
)

var MaxWasmSize = int(math.Pow(2, 24))

func validateWasmCode(code []byte) (bool, error) {
	if len(code) == 0 {
		return false, ErrWasmEmptyCode
	}
	if len(code) > MaxWasmSize {
		return false, ErrWasmCodeTooLarge
	}

	return true, nil
}
