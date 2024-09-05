# Agencia de Quiniela

### Mensajes

Se revisaron los datos de ejemplo y se decidió a partir de ellos el formato que deben respetar los campos de los mensajes:
 - agencia: se permitira una cadena de 3 bytes.
 - nombre: su longitud máxima fue de 23, por lo que se decidió dejar espacio para hasta 30 caracteres.
 - apellido: su longitud máxima fue de 10, por lo que se decidió dejar espacio para hasta 20 caracteres.
 - dni: su valor máximo fue de 39_999_865, por lo que se decidió admitir hasta el valor 4_294_967_295.
 - fecha de nacimiento: dado a que son valores estandar, se utilizarán 10 bytes para la fecha en isoformat YYYY-MM-DD.
 - numero: dado que no conozco nada sobre quiniela, voy a usar 2 bytes para el número. 

Además, los datos simulados de apuestas en la agencias muestran que la máxima cantidad de apuestas por agencia son `max([26935, 25517, 16013, 9237, 990]) = 26935`.
Por esto se tomará `65535` como la máxima cantidad de ganadores por agencias. 

#### 1 SINGLE_BET

Es un mensaje que envía Client a Server. Contiene la información de una única apuesta.

Los campos del mensaje son:

```
┌─────────┬─────────┬────────────────────────────────┬─────────────────────┬──────────┬───────────────────┬─────────┐
│ código  │ agencia │             nombre             │      apellido       │   dni    │fecha de nacimiento│  número │
│ (1 byte)│(3 bytes)│           (30 bytes)           │     (20 bytes)      │(4 bytes) │    (10 bytes)     │(2 bytes)│
├─────────┼─────────┼────────────────────────────────┼─────────────────────┼──────────┼───────────────────┼─────────┤
│ C       │ AAA     │ NNNNNNNNNNNNNNNNNNNNNNNNNNNNNN │ AAAAAAAAAAAAAAAAAAAA│ DDDDDDDD │ YYYY-MM-DD        │ NN      │
└─────────┴─────────┴────────────────────────────────┴─────────────────────┴──────────┴───────────────────┴─────────┘
```


#### 2 SINGLE_BET_ACK

Es un mensaje que envía del Server a Client. Confirma la recepción de la apuesta.

┌─────────┐
│ código  │
│ (1 byte)│
├─────────┤
│ C       │
└─────────┘




#### 3 MULTIPLE_BET

Es un mensaje que envía Client a Server. Contiene la información de n apuestas.

Los campos del mensaje son:

```
┌─────────┬──────────┬────────────┬────────────────────────────────┬
│ código  │ agencia  │  cantidad  │             nombre             │
│ (1 byte)│(3 bytes) │ (1 byte)   │           (30 bytes)           │
├─────────┼──────────┼────────────┼────────────────────────────────┼
│ C       │ AAA      │ C          │ NNNNNNNNNNNNNNNNNNNNNNNNNNNNNN │
└─────────┴──────────┴────────────┴────────────────────────────────┴

─────────────────────┬──────────┬───────────────────┬─────────┬─ ... ─┐
      apellido       │   dni    │fecha de nacimiento│  número │  ...  │
     (20 bytes)      │(4 bytes) │    (10 bytes)     │(2 bytes)│  ...  │
─────────────────────┼──────────┼───────────────────┼─────────┼─ ... ─┤
 AAAAAAAAAAAAAAAAAAAA│ DDDDDDDD │ YYYY-MM-DD        │ NN      │  ...  │
─────────────────────┴──────────┴───────────────────┴─────────┴─ ... ─┘

```


#### 4 MULTIPLE_BET_ACK

Es un mensaje que envía el Server a Client. Confirma la recepción de todas las apuestas.

┌─────────┐
│ código  │
│ (1 byte)│
├─────────┤
│ C       │
└─────────┘

#### 5 NO_MORE_BETS

Es un mensaje que envía el Client a Server. Confirma que ha enviado todas sus apuestas.

┌─────────┬─────────┐
│ código  │ agencia │
│ (1 byte)│(3 bytes)│
├─────────┼─────────┤
│ C       │ AAA     │
└─────────┴─────────┘

#### 6 NO_MORE_BETS_ACK

Es un mensaje que envía el Server a Client. Confirma que se ha recibido el anuncio de fin de envío.

┌─────────┐
│ código  │
│ (1 byte)│
├─────────┤
│ C       │
└─────────┘

#### 7 GET_WINNERS

Es un mensaje que envía el Client a Server. Confirma que ha enviado todas sus apuestas.

┌─────────┬─────────┐
│ código  │ agencia │
│ (1 byte)│(3 bytes)│
├─────────┼─────────┤
│ C       │ AAA     │
└─────────┴─────────┘

#### 8 WAIT

Es un mensaje que envía el Server a Client. Comunica que aún no se ha realizado el sorteo.

┌─────────┐
│ código  │
│ (1 byte)│
├─────────┤
│ C       │
└─────────┘

#### 9 WINNERS

Es un mensaje que envía el Server a Client. Comunica los documentos ganadores de la agencia solicitante.

┌─────────┬────────────┬────────────┬─ ... ─┐
│ código  │  cantidad  │  documento │  ...  │
│ (1 byte)│ (2 bytes)  │(4 bytes)   │  ...  │
├─────────┼────────────┼────────────┼─ ... ─┤
│ C       │ C          │ NNNN       │  ...  │
└─────────┴────────────┴────────────┴─ ... ─┘
