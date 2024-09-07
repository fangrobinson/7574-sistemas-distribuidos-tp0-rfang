# Concurrencia y sincronización

El servidor que recibe apuestas de las agencias puede atender mensajes concurrentemente.

Ante una nueva conexión, atiende la comunicación en un proceso nuevo (`multiprocessing.Process`). 

## Apuestas recibidas y persistidas por el servidor

Se ha protegido el acceso a las apuestas persistidas tanto para lectura como para escritura mediante el uso de un lock de acceso exclusivo.

Además, se han añadido otros recursos compartidos que podrían tener que protegerse:
 - agencias que concluyeron envío de apuestas
 - ganadores del sorteo

Ambos son diccionarios compartidos entre los procesos(`multiprocessing.Manager.dict()`).

## Agencias que concluyeron envío de apuestas

Para el caso del mapa de agencias que concluyeron envío se decidió que no era necesario un mecanismo de sincronización. Ya que se utiliza únicamente para agregar agencias nuevas y consultar el largo de las claves del mapa mediante `len(Manager.dict)`. Esta operación utiliza una vista que se actualiza ante otros cambios, pero además dado que se usa solamente para mantener agencias terminadas únicamente puede crecer. Ya la utilización de `Manager()` implica menos performance que utilizar memoria compartida, por lo que se consideró que en el peor de los casos (e.g. lectura de un valor no actualizado) el servidor respondería con `WAIT` si no se ha alcanzado aún el número necesario de agencias.

## Ganadores del sorteo

Por otro lado, acceder al procesamiento de ganadores del sorteo requiere tomar dos locks, por un lado el Lock de las apuestas persistidas y por otro el del mapa de ganadores. El procesamiento de ganadores se realiza una única vez por el primer proceso que haya encontrado suficientes agencias que terminaron su envío, pero no encontraron ganadores del sorteo. Los ganadores se almacenan en memoria (heap) y se reutilizan para cada nueva consulta.
