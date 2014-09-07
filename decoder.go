package php_session_decoder

import (
	"io"
	"bytes"
	"strings"
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
)

type PhpDecoder struct {
	source	*strings.Reader
	decoder	*php_serialize.UnSerializer
}

func NewPhpDecoder(phpSession string) *PhpDecoder {
	decoder := &PhpDecoder{
		source:		strings.NewReader(phpSession),
		decoder:	php_serialize.NewUnSerializer(""),
	}
	decoder.decoder.SetReader(decoder.source)
	return decoder
}

func (self *PhpDecoder) SetSerializedDecodeFunc(f php_serialize.SerializedDecodeFunc) {
	self.decoder.SetSerializedDecodeFunc(f)
}

func (self *PhpDecoder) Decode() (PhpSession, error) {
	var (
		name	string
		err 	error
		value	php_serialize.PhpValue
	)
	res := make(PhpSession)

	for {
		if name, err = self.readName(); err != nil {
			break
		}
		if value, err = self.decoder.Decode(); err != nil {
			break
		}
		res[name] = value
	}

	if err == io.EOF {
		err = nil
	}
	return res, err
}

func (self *PhpDecoder) readName() (string, error) {
	var (
		token byte
		err   error
	)
	buf := bytes.NewBuffer([]byte{})
	for {
		if token, err = self.source.ReadByte(); err != nil || token == byte(SEPARATOR_VALUE_NAME) {
			break
		} else {
			buf.WriteByte(token)
		}
	}
	return buf.String(), err
}
