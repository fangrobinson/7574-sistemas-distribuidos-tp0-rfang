# Informe - TP0: Docker + Comunicaciones + Concurrencia

|           Alumno                | Padrón |
|---------------------------------|--------|
|         Fang, Robinson          |  97009 |

|           Corrector             |
|---------------------------------|
|       Franco Barreneche         |


# Objetivo

El objetivo de este trabajo práctico es diseñar e implementar de forma iterativa una solución respetando la arquitectura cliente-servidor que modele correctamente múltiples agencias de lotería comunicando apuestas al encargado del sorteo. El sorteo se realiza una vez el mínimo de agencias esperado termina de comunicar los partiicpantes del sorteo.

Se partió a partir de un echo server y cliente provistos por la cátedra: [repositorio](https://github.com/7574-sistemas-distribuidos/tp0-base).

Se buscó familiarizarse con las distintas herramientas requeridas: `make`, `docker`, `python`, `golang`, `despliegue de servicios configurables`, como así también automatización de tareas (e.g. la generación de un `docker-compose` válido) y el diseño de un protocolo de comunicación orientado a bytes.

# Solución

## Servidor

Está implementado en python. El mismo atiende de forma concurrente conexiones de clientes que son abiertas y cerradas por cada para de mensaje intercambiado. Los procesos encargados de atender a los clientes comparten recursos para solucionar las peticiones y además protegen el acceso a secciones críticas. Los mecanismos de sincronización utilizados se encuentran descritos en el siguiente apartado [Sincronización](./sincronizacion.md).

Dado que para cada petición se atiende la misma en un proceso nuevo, ante la llegada de un `SIGTERM` se notifica a los mismos que deben terminar, para que puedan liberar sus recursos y concluir la ejecución de forma graceful.

#### Consideraciones y puntos de mejora.

Se aceptan conexiones para cada mensaje. Dado que se utiliza TCP hay un overhead importante de mensajes intercambiados por el protocolo para establecer las mismas.

No existe limpieza de procesos creados durante la ejecución del servidor. Podrían ir removiéndose los recursos a medida que van terminando y no una única vez ante el cierre de la aplicación.

## Protocolo

El detalle de los mensajes intercambiados del protocolo puede consultarse en el siguiente apartado [Protocolo](./protocolo.md).

Paticularidades del protocolo:
 - todos los mensajes comienzan con el código de mensaje en cuestión
 - todos los mensajes enviados por clientes contienen además un campo de ID de cliente
 - los campos de los mensajes son de largo fijo
   - cada campo de texto plano tendrá padding para completar el espacio debidamente.
   - dos mensajes (`multiple_bets` y `winners`) en el segmento de datos pueden contener `n` bets o winners.

#### Consideraciones

Dado que se crean conexiones para cada mensaje no existe concepto de `sesión` en los mensajes intercambiados. Por esto, los mensajes de los clientes contienen el código de identificación de la agencia que podría ahorrarse si ante cada `sesión` se anunciara la agencia en cuestión mediante un mensaje diferente.

## Cliente

El cliente está implementado en golang. La ejecución de un cliente puede dividirse en las siguientes partes:

 - envío de apuestas
 - aviso de fin de envío
 - consulta de ganadores
 - comunicado de los resultados


#### envío de apuestas

Las apuestas se leen de un archivo que contiene los datos simulados de clientes de la agencia. En un mismo mensaje se comunican `n` participantes con `n` configurable. Se ha considerado que no deben enviarse más de 8kbs, por lo que el número de participantes por mensaje no debe ser mayor a 121.

#### aviso de fin de envío

Una vez enviados todos los participantes se procede a comunicar esto al servidor y a esperar que esto sea reconocido.

#### consulta de ganadores

La consulta de ganadores se realiza mediante polling. Ante un aviso de espera, se procede a esperar y repetir el pedido.

#### comunicado de los resultados

Aunque el servidor debería enviar únicamente números de documentos de ganadores de la agencia, se verifica esta información. La lista de ganadores se presume mucho menor a la cantidad de participantes de la agencia, por lo que se convierte esta en un mapa y se itera luego sobre los participantes para chequear pertenencia y así eficientizar el cómputo de los ganadores.

#### Consideraciones

El archivo se mantiene abierto durante el envío de apuestas, por lo que si quisiera añadir concurrencia en esta sección debería añadirse sincronización al estado de procesamiento del archivo.

# Cómo correr los servicios

Para correr los servicios es necesario descomprimir la data incluida en ./data/dataset.zip o se deben proveer `.csv` propios ubicados en `./data/dataset/` con nombres que sigan la convención `agency-${NRO_DE_AGENCIA}.csv`.

Se pueden declarar los servicios, redes y recursos mediante un [script](./generar-compose.sh). El mismo debe ejecutarse de la siguiente manera:

```bash
./generar-compose.sh FILENAME_DOCKER_COMPOSE CANTIDAD_DE_CLIENTES
```

Para iniciar los servicios se deben utilizar los comandos definidos en [Makefile](./Makefile). 

* **make \<target\>**:
Los target imprescindibles para iniciar y detener el sistema son **docker-compose-up** y **docker-compose-down**, siendo los restantes targets de utilidad para el proceso de _debugging_ y _troubleshooting_.

Los targets disponibles son:
* **docker-compose-up**: Inicializa el ambiente de desarrollo (buildear docker images del servidor y cliente, inicializar la red a utilizar por docker, etc.) y arranca los containers de las aplicaciones que componen el proyecto.
* **docker-compose-down**: Realiza un `docker-compose stop` para detener los containers asociados al compose y luego realiza un `docker-compose down` para destruir todos los recursos asociados al proyecto que fueron inicializados. Se recomienda ejecutar este comando al finalizar cada ejecución para evitar que el disco de la máquina host se llene.
* **docker-compose-logs**: Permite ver los logs actuales del proyecto. Acompañar con `grep` para lograr ver mensajes de una aplicación específica dentro del compose.
* **docker-image**: Buildea las imágenes a ser utilizadas tanto en el servidor como en el cliente. Este target es utilizado por **docker-compose-up**, por lo cual se lo puede utilizar para testear nuevos cambios en las imágenes antes de arrancar el proyecto.
* **build**: Compila la aplicación cliente para ejecución en el _host_ en lugar de en docker. La compilación de esta forma es mucho más rápida pero requiere tener el entorno de Golang instalado en la máquina _host_.
