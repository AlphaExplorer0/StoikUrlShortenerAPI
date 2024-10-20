package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_hashAndCut_produces_expected_output_format(t *testing.T) {

	url := "https://fr.wikipedia.org/wiki/Uniform_Resource_Locator"
	tiny := hashAndCut(url, 7)
	fmt.Printf("%s\n", tiny)
	assert.Equal(t, 7, len(tiny))
}
