package uniqueid

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const (
	idLength         = 14
	timePartLength   = 9
	randomPartLength = 5
)

// New creates a new 14-character, uppercase, alphanumeric, time-ordered unique ID.
func New() string {
	now := time.Now().UnixMilli()
	timePart := strings.ToUpper(strconv.FormatInt(now, 36))

	if len(timePart) < timePartLength {
		timePart = strings.Repeat("0", timePartLength-len(timePart)) + timePart
	}

	maxRandom := new(big.Int)
	maxRandom.Exp(big.NewInt(36), big.NewInt(randomPartLength), nil)

	randomInt, err := rand.Int(rand.Reader, maxRandom)
	if err != nil {
		panic(fmt.Sprintf("failed to generate random number for unique ID: %v", err))
	}

	randomPart := strings.ToUpper(randomInt.Text(36))
	paddedRandomPart := strings.Repeat("0", randomPartLength-len(randomPart)) + randomPart

	return timePart + paddedRandomPart
}
