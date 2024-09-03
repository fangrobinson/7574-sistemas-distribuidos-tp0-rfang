# Agencia de Quiniela

### Mensajes

Se revisaron los datos de ejemplo y se decidió a partir de ellos el formato que deben respetar los campos de los mensajes:
 - agencia: se permitira una cadena de 3 bytes.
 - nombre: su longitud máxima fue de 23, por lo que se decidió dejar espacio para hasta 30 caracteres.
 - apellido: su longitud máxima fue de 10, por lo que se decidió dejar espacio para hasta 20 caracteres.
 - dni: su valor máximo fue de 39_999_865, por lo que se decidió admitir hasta el valor 4_294_967_295.
 - fecha de nacimiento: dado a que son valores estandar, se utilizarán 10 bytes para la fecha en isoformat YYYY-MM-DD.
 - numero: dado que no conozco nada sobre quiniela, voy a usar 2 bytes para el número. 


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
