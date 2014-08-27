<?php

class TestSessionHandler implements SessionHandlerInterface
{
    private $dir;

    public function __construct()
    {
        $this->dir = dirname(__FILE__) . '/';
    }

    public function getFileName($id)
    {
        return $this->dir . 'session_' . $id;
    }

    function open($savePath, $sessionName)
    {
        return true;
    }

    function close()
    {
        return true;
    }

    function read($id)
    {
        $data = @file_get_contents($this->getFileName($id));
        return false === $data ? '' : $data;
    }

    function write($id, $data)
    {
        echo json_encode($data).PHP_EOL;
        return file_put_contents($this->getFileName($id), $data) === false ? false : true;
    }

    function destroy($id)
    {
        return true;
    }

    function gc($maxlifetime)
    {
        return true;
    }
}

class AbcClass {
    public $a = 5;
    private $b = 'private';
    protected $c = 8;
}

class TestObject implements Serializable {
    public $item;
    public function serialize() {
        return serialize(array('item' => $this->item));
    }
    public function unserialize($serialized) {
        return unserialize($serialized);
    }
}

$object = new TestObject();
$object->item = new AbcClass();

$handler = new TestSessionHandler();
session_set_save_handler($handler, true);
session_start();

$_SESSION['object'] = $object;

echo $handler->getFileName(session_id()) . PHP_EOL;