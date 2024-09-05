from .message import Message
from .types import MessageType

AGENCY_FIELD_LENGTH = 3

class NoMoreBetsMessage(Message):
    CODE = MessageType.NO_MORE_BETS
    DATA_LENGTH = AGENCY_FIELD_LENGTH

    def data_length(self, header) -> int:
        return self.DATA_LENGTH

    def parse(self, header, message_bytes) -> list:
        agency = message_bytes[:AGENCY_FIELD_LENGTH].decode("utf-8").rstrip()
        return [agency]
