from typing import List

from .message import Message
from .types import MessageType

AGENCY_FIELD_LENGTH = 3
FIRST_NAME_FIELD_LENGTH = 30
LAST_NAME_FIELD_LENGTH = 20
DOCUMENT_FIELD_LENGTH = 4
BIRTHDATE_FIELD_LENGTH = 10
NUMBER_FIELD_LENGTH = 2


class SingleBetMessage(Message):
    CODE = MessageType.SINGLE_BET
    DATA_LENGTH = sum((
        AGENCY_FIELD_LENGTH,
        FIRST_NAME_FIELD_LENGTH,
        LAST_NAME_FIELD_LENGTH,
        DOCUMENT_FIELD_LENGTH,
        BIRTHDATE_FIELD_LENGTH,
        NUMBER_FIELD_LENGTH,
    ))

    def data_length(self, header) -> int:
        return self.DATA_LENGTH

    def parse(self, header, message_bytes) -> list:
        read_bytes = 0
        agency = message_bytes[read_bytes:read_bytes+AGENCY_FIELD_LENGTH].decode("utf-8").rstrip()
        read_bytes += AGENCY_FIELD_LENGTH
        first_name = message_bytes[read_bytes:read_bytes+FIRST_NAME_FIELD_LENGTH].decode("utf-8").rstrip()
        read_bytes += FIRST_NAME_FIELD_LENGTH
        last_name = message_bytes[read_bytes:read_bytes+LAST_NAME_FIELD_LENGTH].decode("utf-8").rstrip()
        read_bytes += LAST_NAME_FIELD_LENGTH
        document = int.from_bytes(message_bytes[read_bytes:read_bytes+DOCUMENT_FIELD_LENGTH], "big")
        read_bytes += DOCUMENT_FIELD_LENGTH
        birthdate = message_bytes[read_bytes:read_bytes+BIRTHDATE_FIELD_LENGTH].decode("utf-8").rstrip()
        read_bytes += BIRTHDATE_FIELD_LENGTH
        number = int.from_bytes(message_bytes[read_bytes:read_bytes+NUMBER_FIELD_LENGTH], "big")
        read_bytes += NUMBER_FIELD_LENGTH
        return [agency, first_name, last_name, document, birthdate, number]
