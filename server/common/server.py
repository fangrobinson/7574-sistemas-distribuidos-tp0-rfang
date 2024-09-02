import socket
import logging
import signal

from common.connection.message_receiver import MessageReceiver
from common.connection.types import MessageType
from common.utils import Bet, store_bets

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._run = True
        signal.signal(signal.SIGTERM, self.__handle_sigterm)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while self._run:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except OSError as _:
                if not self._run:
                    break
                raise

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # TODO: Modify the receive to avoid short-reads
            msg = MessageReceiver.recv(client_sock)
            if msg[0] == MessageType.SINGLE_BET.value:
                agency, first_name, last_name, document, birthdate, number = msg[1:]
                bet = Bet(agency, first_name, last_name, document, birthdate, number)
                store_bets([bet])
                logging.info(f'action: apuesta_almacenada | result: success | dni: ${document} | numero: ${number}')
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # TODO: Modify the send to avoid short-writes
            client_sock.send("{}\n".format(msg).encode('utf-8'))
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def __handle_sigterm(self, signum, frame):
        # TODO : add compliant logging
        self._server_socket.close()
        self._run = False
