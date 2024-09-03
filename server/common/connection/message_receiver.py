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
        code = int.from_bytes(cls._recv(client_socket, cls.CODE_LENGHT), "big", signed=False)
        if code not in cls.MESSAGE_TYPES:
            logging.warning("Unknown message type:", code)
            return ["UNKNOWN"]
        msg_type = cls.MESSAGE_TYPES[code]
        data = msg_type.parse(cls._recv(client_socket, msg_type.DATA_LENGTH))
        return [code] + data

    @classmethod
    def _recv(cls, client_socket: socket.socket, n: int) -> bytes:
        buffer = bytearray()
        while len(buffer) < n:
            chunk = client_socket.recv(n - len(buffer))
            if not chunk:
                # TODO: proper logging.
                raise RuntimeError("")
            buffer.extend(chunk)
        return buffer