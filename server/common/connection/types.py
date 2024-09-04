from enum import Enum


class MessageType(Enum):
    SINGLE_BET = 1          # Client -> Server: a single bet
    SINGLE_BET_ACK = 2      # Server -> Client: received bet
    MULTIPLE_BET = 3        # Client -> Server: multiple bets
    MULTIPLE_BET_ACK = 4    # Server -> Client: received bets
