package php_session_decoder

import "github.com/yvasiyarov/php_session_decoder/php_serialize"

const SEPARATOR_VALUE_NAME = '|'

type PhpSession map[string]php_serialize.PhpValue
