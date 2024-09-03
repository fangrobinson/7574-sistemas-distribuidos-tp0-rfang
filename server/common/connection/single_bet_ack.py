import logging
from typing import List

from .message import Message
from .types import MessageType


class SingleBetAckMessage(Message):
    CODE = MessageType.SINGLE_BET_ACK
    DATA_LENGTH = 0

    def parse(self, message_bytes) -> List[str]:
        self.validate(message_bytes)
        return []

    def validate(self, message_bytes):
        return

    @classmethod
    def encode(cls, *args, **kwargs) -> bytes:
        return cls.CODE.value.to_bytes(1, "big")