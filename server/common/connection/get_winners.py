from .message import Message
from .types import MessageType

AGENCY_FIELD_LENGTH = 3

class GetWinnersMessage(Message):
    CODE = MessageType.GET_WINNERS
    DATA_LENGTH = AGENCY_FIELD_LENGTH

    def data_length(self, header) -> int:
        return self.DATA_LENGTH

    def parse(self, header, message_bytes) -> list:
        agency = message_bytes[:AGENCY_FIELD_LENGTH].decode("utf-8").rstrip()
        return [agency]

    @classmethod
    def encode(cls, *args, **kwargs) -> bytes:
        return cls.CODE.value.to_bytes(1, "big")
