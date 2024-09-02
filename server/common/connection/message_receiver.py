import logging
import socket


class MessageReceiver:

    CODE_LENGHT = 1
    MESSAGE_TYPES = {}

    @classmethod
    def register_message(cls, message_cls):
        cls.MESSAGE_TYPES[message_cls.CODE.value] = message_cls()
    
    @classmethod
    def recv(cls, client_socket: socket.socket) -> list:
        code = int.from_bytes(client_socket.recv(cls.CODE_LENGHT), "big", signed=False)
        if code not in cls.MESSAGE_TYPES:
            logging.warning("Unknown message type:", code)
            return ["UNKNOWN"]
        msg_type = cls.MESSAGE_TYPES[code]
        data = msg_type.parse(client_socket.recv(msg_type.DATA_LENGTH))
        return [code] + data
