<?php

class Test {
    public $public = 1;
    protected $protected = 2;
    private $private = 3;
}

class Test1 {
    public $public = 11;
    protected $protected = 12;
    private $private = 13;
}

class Test2 {
    public $public = 21;
    protected $protected = 22;
    private $private = 23;
}

class TestSerializable implements Serializable {
    public function serialize() {
        return "foobar";
    }
    public function unserialize($str) {
        // ...
    }
}

class TestSerializable1 implements Serializable {
    public $foo = 4;
    public $bar = 2;

    public function serialize() {
        return serialize(array('foo' => $this->foo, 'bar' => $this->bar));
    }
    public function unserialize($str) {
        // ...
    }
}

class TestSerializable2 implements Serializable {
    public $foo = 4;
    public $bar = 2;

    public function serialize() {
        return json_encode(array('foo' => $this->foo, 'bar' => $this->bar));
    }
    public function unserialize($str) {
        // ...
    }
}

echo '===============' . PHP_EOL;

serializeMe('Nul', null);
serializeMe('Bool True', true);
serializeMe('Bool False', false);
serializeMe('Int', 42);
serializeMe('Int Minus', -42);
serializeMe('Float', 42.3789);
serializeMe('Float Minus', -42.3789);
serializeMe('String', 'foobar');
serializeMe('Array', array(10, 11, 12));
serializeMe('Array Keys', array('foo' => 4, 'bar' => 2));
serializeMe('Array Array', array('foo' => array(10, 11, 12), 'bar' => 2));
serializeMe('Object', new Test());
serializeMe('Array Object', array(new Test1(), new Test2()));
serializeMe('Serializable Empty Object', new TestSerializable());
serializeMe('Serializable Array Object', new TestSerializable1());
serializeMe('Serializable JSON Object', new TestSerializable2());

echo '===============' . PHP_EOL;

function serializeMe($message, $object) {
    echo $message . PHP_EOL;
    echo json_encode(serialize($object)) . PHP_EOL . PHP_EOL;
}