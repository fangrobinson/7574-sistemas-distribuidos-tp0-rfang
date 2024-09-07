from enum import Enum


class MessageType(Enum):
    SINGLE_BET = 1          # Client -> Server: a single bet
    SINGLE_BET_ACK = 2      # Server -> Client: received bet
    MULTIPLE_BET = 3        # Client -> Server: multiple bets
    MULTIPLE_BET_ACK = 4    # Server -> Client: received bets
    NO_MORE_BETS = 5        # Client -> Server: finished sending bets
    NO_MORE_BETS_ACK = 6    # Server -> Client: notified
    GET_WINNERS = 7         # Client -> Server: request for winner
    WAIT = 8                # Server -> Client: no winner yet
    WINNERS = 9             # Server -> Client: winner number
