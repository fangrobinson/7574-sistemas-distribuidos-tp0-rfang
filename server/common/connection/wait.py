from .message import Message
from .types import MessageType


class WaitMessage(Message):
    CODE = MessageType.WAIT

    def parse(self, header, message_bytes) -> list:
        return []

    @classmethod
    def encode(cls, *args, **kwargs) -> bytes:
        return cls.CODE.value.to_bytes(1, "big")
