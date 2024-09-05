import socket


class MessageSender:
    
    @classmethod
    def send(cls, client_socket: socket.socket, data):
        sent = 0
        while sent < len(data):
            _sent = client_socket.send(data[sent:])
            if _sent == 0:
                raise RuntimeError("Socket was closed")
            sent += _sent
