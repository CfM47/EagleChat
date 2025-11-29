# Informe sobre EagleChat

                           .a@@@@@#########@@@@a.
                       .a@@######@@@mm@@mm######@@@a.
                  .a####@@@@@@@@@@@@@@@@@@@mm@@##@@v;%%,.
               .a###v@@@@@@@@vvvvvvvvvvvvvv@@@@#@v;%%%vv%%,
            .a##vv@@@@@@@@vv%%%%;S,  .S;%%vv@@#v;%%'/%vvvv%;
          .a##@v@@@@@vv%%vvvvvv%%;SssS;%%vvvv@v;%%./%vvvvvv%;
        ,a##vv@@@vv%%%@@@@@@@@@@@@mmmmmmmmmvv;%%%%vvvvvvvvv%;
        .a##@@@@@@@@@@@@@@@@@@@@@@@mmmmmvv;%%%%%vvvvvvvvvvv%;
       ###vv@@@v##@v@@@@@@@@@@mmv;%;%;%;%;%;%;%;%;%;%;%,%vv%'
      a#vv@@@@v##v@@@@@@@@###@@@@@%v%v%v%v%v%v%v%      ;%%;'
     ',a@@@@@@@v@@@@@@@@v###v@@@nvnvnvnvnvnvnvnv'     .%;'
     a###@@@@@@@###v@@@v##v@@@mnmnmnmnmnmnmnmnmn.     ;'
    ,###vv@@@@v##v@@@@@@v@@@@v##v@@@@@v###v@@@##@.
    ###vv@@@@@@v@@###v@@@@@@@@v@@@@@@v##v@@@v###v@@.


## Arquitectura

### Roles

Esta aplicación tiene dos roles, `clients` y `id managers`.

#### Clientes (Clients)

- Estos contienen la aplicación de mensajería, donde se muestran los mensajes propios y los enviados hacia el cliente. 
- En estos se almacenan los mensajes en tránsito y los que llegaron a su destino, esta información no se guarda en ningún otro tipo de nodo.

Nótese que los clientes pueden tener mensajes (encriptados) para otros clientes, cacheados temporalmente.

#### ID Managers

- Estos se encargan de almacenar los identificadores únicos de los clientes, sus nombres de usuario, sus IPs y sus llaves públicas.
- Permiten el registro de nuevos clientes, aceptando un nombre de usuario y una llave pública, y devolviendo su nuevo identificador.
- Se encargan además de tener un registro de identificadores de mensajes pendientes (no los mensajes en sí), que contienen la información de los remitentes y de los clientes que los tienen cacheados, esta información se guarda temporalmente.
- Permiten la petición de una cantidad de clientes conectados aleatorios, para el cacheo de mensajes distribuida por la red.

### Distribución de servicios entre redes de docker

El sistema no requiere ninguna distribución específica, solo debe existir un `id manager` en cada red para lograr el descubrimiento entre clientes, y su registro en el sistema.

### Organización alto nivel

La idea es sencilla, los clientes inicialmente se registran con un `id manager`, luego, para poder comunicarse con otro cliente, piden su información a un `id manager`, a través de su identificador único, y al recibir la IP, comienza una comunicación directa con el cliente objetivo para enviar sus mensajes.

## Procesos que corren en el sistema

En este sistema, al estar desarrollado en go, e implementarse la concurrencia a través de `goroutines`, el uso de hilos o procesos está manejado por el runtime de go, no directamente por los desarrolladores.

### En ID Managers

Estos servicios son mayormente reactivos, es decir, sirven peticiones iniciadas por los clientes. La única excepción es el intercambio de datos entre `id managers`, un `id manager` puede iniciar una comunicación con otro, para lograr la sincronía de los datos. Por tanto, además de los hilos que deben correr constantemente para el mantenimiento del servidor HTTP necesario para reaccionar a las peticiones de los clientes, en cada `id manager` existe una `goroutine` (se menciona esto, pues el manejo de esto con hilos o procesos es manejado por go) que, después de un intervalo de tiempo periódico, o cuando su base de datos es modificada, inicia comunicación con cierto número de `id managers` aleatorios.

### En Clientes

En los clientes existen varias `groutines` (que pueden estar implementadas como hilos o procesos) corriendo de forma concurrente.

#### Announcer

Este se encarga de anunciar su presencia a los `id managers`, cada cierto intervalo de tiempo, socializa su identificador (su IP viene en la petición) y los mensajes para otros clientes que tiene cacheados.

#### Message Sender

Este se encarga de periódicamente pedir las IPs de los clientes que son remitentes de mensajes que se cachean en el cliente iniciador, y enviar todos dichos mensajes a su destino, eliminándolos del caché si se conoce que llegaron a su objetivo.

#### Receiver

Este es simplemente un servidor HTTP, con todos los hilos o procesos necesarios para su implementación, que se encarga de servir peticiones de otros clientes, esencialmente recibiendo mensajes, tanto para sí mismo, como para otros clientes.

#### Interfaz de Usuario

Para esta por supuesto existe una serie de `goroutines` para manejar cosas como las entradas del usuario y seguir siendo reactiva, la implementación está manejada por la biblioteca [tview](github.com/rivo/tview).

## Comunicación

Todas los mensajes entre nodos se envían a través de APIs REST HTTP. Existiendo entre todas las (3) combinaciones de roles.

### Clientes

También existe una implementación `in-house` de colas de mensajes locales, que es lo que usan los clientes para cachear en la red sus mensajes pendientes: Cuando un cliente no puede enviar directamente un mensaje a su objetivo, este cachea el mensaje como un **mensaje pendiente** en una base de datos local, estos mensajes están encriptados, pero exponen un identificador único para el mensaje, además de el identificador del remitente. Luego:

- Periódicamente, anuncia los mensajes pendientes que tiene cacheados a los `id managers`.
- Periódicamente, pide las IPs de cierta cantidad de clientes conectados, y envía sus mensajes pendientes a estos, efectivamente cacheando dichos mensajes en la red.
- Periódicamente, intenta enviar los mensajes pendientes directamente a sus remitentes, eliminándolos si es exitoso.

La única diferencia entre mensajes pendientes creados por el propio cliente y que otros clientes cachearon en él, es que los propios nunca son eliminados del caché local hasta ser exitosamente enviados al objetivo, los de los otros clientes se eliminan después de cierto tiempo.

### ID Managers

Entre los `id managers`, como se explicó brevemente en la sección de Procesos, existe comunicaciones periódicas y reactivas. Cuando un cliente se registra, el `id manager` que manejó el registro envía este cambio a tantos `id managers` como pueda encontrar, además, periódicamente existen comunicaciones entre `id managers`, donde se comparten sus datos, para lograrse la consistencia eventual. 

## Coordinación

Las especificidades del sistema hacen innecesarios algoritmos complejos de coordinación.

## Nombrado y descubrimiento

El descubrimiento se logra principalmente a través del DNS de docker, para obtener todas las IPs de los `id managers`, lo cual es necesario para clientes intentando registrarse o buscando las IP de otros clientes para el envío de mensajes, o `id managers` buscando otros `id managers` para socializar sus actualizaciones por ejemplo. Los clientes periódicamente se anuncian a tantos `id managers` como pueden, para socializar sus IPs.

Como plan B, los nodos cachean las IPs de `id managers` u otros clientes, en caso de que falle el servicio de DNS, o para no tener que hacer peticiones repetidas a este en intervalos cortos de tiempo.

## Consistencia y replicación

Como se ha explicado anteriormente, los clientes logran replicación enviando sus mensajes pendientes (que pueden ser creados por otros clientes) a clientes aleatorios en la red cuando sus objetivos no son alcanzables. Además, a través de intercambios de datos directos, periódicos y después de actualizaciones, los `id managers` logran la consistencia eventual, y la tolerancia a particionamiento.

El hecho de que cada `id manager` tiene esencialmente una base de datos llave-valor con llaves únicas, donde los valores (aparte de la IP) no cambian, hace el logro de la consistencia entre `id managers` relativamente sencillo. Para los datos de remitentes de mensajes pendientes, estos son incluso más sencillos, solo se tienen llaves, por tanto la consistencia se puede lograr simplemente con uniones entre bases de datos.

## Tolerancia a fallos

el sistema no puede fallar :)

## Seguridad

todo muy safe
