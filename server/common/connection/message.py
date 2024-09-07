from .message_receiver import MessageReceiver

class MessageMeta(type):
    def __new__(cls, name, bases, attrs):
        new_class = super().__new__(cls, name, bases, attrs)
        if "CODE" in attrs:
            MessageReceiver.register_message(new_class)
        return new_class


class Message(metaclass=MessageMeta):
    HEADER_LENGTH = 0

    def data_length(self, header) -> int:
        return 0

    def parse_header(self, message_bytes) -> list:
        return list()

    def parse(self, header, message_bytes) -> list: ...

    @classmethod
    def encode(cls, *args, **kwargs) -> bytes:
        return bytearray(cls.CODE.value)
