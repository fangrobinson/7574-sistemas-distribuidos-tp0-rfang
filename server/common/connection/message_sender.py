import socket


class MessageSender:
    
    @classmethod
    def send(cls, client_socket: socket.socket, msg):
        data = msg.encode()
        sent = 0
        while sent < len(data):
            _sent = client_socket.send(data[sent:])
            if _sent == 0:
                # TODO: proper logging.
                raise RuntimeError("")            
            sent += _sent
