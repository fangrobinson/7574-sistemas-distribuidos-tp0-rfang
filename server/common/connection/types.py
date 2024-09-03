from enum import Enum


class MessageType(Enum):
    AGENCY = 0
    SINGLE_BET = 1          # Client -> Server: a single bet
    SINGLE_BET_ACK = 2      # Server -> Client: received bet
