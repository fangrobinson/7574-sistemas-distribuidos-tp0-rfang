from .message import Message
from .types import MessageType

AGENCY_FIELD_LENGTH = 3
BETS_AMOUNT_LENGHT = 1

FIRST_NAME_FIELD_LENGTH = 30
LAST_NAME_FIELD_LENGTH = 20
DOCUMENT_FIELD_LENGTH = 4
BIRTHDATE_FIELD_LENGTH = 10
NUMBER_FIELD_LENGTH = 2


class MultipleBetMessage(Message):
    CODE = MessageType.MULTIPLE_BET
    HEADER_LENGTH = AGENCY_FIELD_LENGTH + BETS_AMOUNT_LENGHT
    DATA_LENGTH = sum((
        FIRST_NAME_FIELD_LENGTH,
        LAST_NAME_FIELD_LENGTH,
        DOCUMENT_FIELD_LENGTH,
        BIRTHDATE_FIELD_LENGTH,
        NUMBER_FIELD_LENGTH,
    ))

    def data_length(self, header) -> int:
        _, bets_amount = header
        return bets_amount * self.DATA_LENGTH

    def parse_header(self, message_bytes) -> list:
        agency = message_bytes[0:AGENCY_FIELD_LENGTH].decode("utf-8").rstrip()
        bets_amount = int.from_bytes(message_bytes[AGENCY_FIELD_LENGTH:AGENCY_FIELD_LENGTH+BETS_AMOUNT_LENGHT], "big")
        return [agency, bets_amount]

    def parse(self, header, message_bytes) -> list:
        _, bets_amount = header
        read_bytes = 0
        bets = []
        for _ in range(bets_amount):
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
            bets.append([first_name, last_name, document, birthdate, number])
        return [bets]
