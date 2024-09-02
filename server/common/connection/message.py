from typing import List

from .message_receiver import MessageReceiver

class MessageMeta(type):
    def __new__(cls, name, bases, attrs):
        new_class = super().__new__(cls, name, bases, attrs)
        if "CODE" in attrs:
            MessageReceiver.register_message(new_class)
        return new_class


class Message(metaclass=MessageMeta):
    DATA_LENGTH = 0
    def parse(self, message_bytes) -> List[str]: ...

    @classmethod
    def encode(cls, *args, **kwargs) -> bytes:
        return bytearray(cls.CODE.value)