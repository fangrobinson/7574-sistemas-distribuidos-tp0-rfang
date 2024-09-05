from .message import Message
from .types import MessageType


class WinnersMessage(Message):
    CODE = MessageType.WINNERS

    def parse(self, header, message_bytes) -> list:
        return []

    @classmethod
    def encode(cls, winners: list[int], *args, **kwargs) -> bytes:
        message_bytes = bytearray()
        message_bytes.extend(cls.CODE.value.to_bytes(1, "big"))
        message_bytes.extend(len(winners).to_bytes(2, "big"))
        for winner in winners:
            message_bytes.extend(int(winner).to_bytes(4, "big"))
        return message_bytes
