package php_session_decoder;
/*
import (
    "errors"
    "fmt"
)
*/

/*
func encode() {
}
*/


type PhpValue interface{}

type PhpSessionData map[string]PhpValue

type PhpDecoder struct {
    source string
    position int
    data *PhpSessionData
}

func NewPhpDecoder(phpSession string) (*PhpDecoder) {
    sessionData := make(PhpSessionData)
    d := &PhpDecoder{
        source: phpSession, 
        position: 0,
        data: &sessionData,
    }
    return d
}

func (decoder *PhpDecoder)Decode() (*PhpSessionData, error) {
    return decoder.data, nil
}
