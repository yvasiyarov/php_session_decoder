package php_session_decoder

import (
	"bytes"
	"fmt"

	"github.com/yvasiyarov/php_session_decoder/php_serialize"
)

type PhpEncoder struct {
	data    PhpSession
	encoder *php_serialize.Serializer
}

func NewPhpEncoder(data PhpSession) *PhpEncoder {
	return &PhpEncoder{
		data:    data,
		encoder: php_serialize.NewSerializer(),
	}
}

func (self *PhpEncoder) SetSerializedEncodeFunc(f php_serialize.SerializedEncodeFunc) {
	self.encoder.SetSerializedEncodeFunc(f)
}

func (self *PhpEncoder) Encode() (string, error) {
	if self.data == nil {
		return "", nil
	}
	var (
		err error
		val string
	)
	buf := bytes.NewBuffer([]byte{})

	for k, v := range self.data {
		buf.WriteString(k)
		buf.WriteRune(SEPARATOR_VALUE_NAME)
		if val, err = self.encoder.Encode(v); err != nil {
			err = fmt.Errorf("php_session: error during encode value for %q: %v", k, err)
			break
		}
		buf.WriteString(val)
	}

	return buf.String(), err
}
