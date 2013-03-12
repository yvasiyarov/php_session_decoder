package php_session_decoder

import (
	"bytes"
//	"errors"
//	"fmt"
//	"strconv"
//	"strings"
)


type PhpEncoder struct {
	dest     bytes.Buffer
	data     *PhpSessionData
}

func NewPhpEncoder(sessionData *PhpSessionData) *PhpEncoder {
        var buffer bytes.Buffer
	e := &PhpEncoder{
		dest:   buffer,
		data:   sessionData,
	}
	return e
}

