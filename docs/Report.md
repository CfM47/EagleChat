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

En los clientes existen varias `goroutines` (que pueden estar implementadas como hilos o procesos) corriendo de forma concurrente.

#### Announcer

Este se encarga de anunciar su presencia a los `id managers`, cada cierto intervalo de tiempo, socializa su identificador (su IP viene en la petición) y los mensajes para otros clientes que tiene cacheados.

#### Message Sender

Este se encarga de periódicamente pedir las IPs de los clientes que son destinatarios de mensajes que se cachean en el cliente iniciador, y enviar todos dichos mensajes a su destino, eliminándolos del caché si se conoce que llegaron a su objetivo.

#### Receiver

Este es simplemente un servidor HTTP, con todos los hilos o procesos necesarios para su implementación, que se encarga de servir peticiones de otros clientes, esencialmente recibiendo mensajes, tanto para sí mismo, como para otros clientes.

#### Interfaz de Usuario

Para esta por supuesto existe una serie de `goroutines` para manejar cosas como las entradas del usuario y seguir siendo reactiva, la implementación está manejada por la biblioteca [tview](https://github.com/rivo/tview).

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

La tolerancia a fallos del sistema se logra principalmente mediante la replicación distribuida y aleatoria de la información. Cada mensaje pendiente se encuentra almacenado simultáneamente en $m > k$ clientes aleatorios, lo que implica que el sistema puede tolerar la caída de al menos $k$ de estos nodos sin perder disponibilidad: mientras al menos uno permanezca en línea, el mensaje puede seguir siendo reenviado hacia su objetivo.

Se podrá descubrir nuevos peers mientras exista al menos un id manager activo, los clientes pueden seguir registrándose, anunciándose y resolviendo identificadores hacia IPs. El sistema mantiene consistencia eventual entre los id managers, por lo que la falla de algunos de ellos no afecta la operación general.

El caso más crítico ocurre cuando todos los id managers fallan simultáneamente. En tal escenario, los clientes ya conectados pueden seguir comunicándose con pares cuyo IP haya sido previamente cacheado, pero el sistema no puede resolver nuevos identificadores ni registrar nuevos clientes. En esencia, la red continúa funcionando en un modo degradado: el envío de mensajes entre nodos ya descubiertos sigue siendo posible, pero no puede realizarse nuevo descubrimiento ni mantenimiento del estado global hasta que al menos un id manager vuelva a estar disponible.

## Seguridad

La seguridad del sistema se basa en un modelo distribuido donde los `clients` son los únicos responsables de todas las operaciones criptográficas importantes, mientras que los `id managers` actúan únicamente como directorios públicos de llaves e identidades.

### Identidad y autenticación de clientes

Cada cliente funciona como su propio proveedor de identidad: al iniciarse por primera vez, genera de forma local un par de llaves RSA, donde la llave privada actúa como prueba criptográfica de identidad. Esta clave privada nunca sale del dispositivo y no es compartida con ningún otro nodo. La autenticidad de los mensajes entre clientes se garantiza mediante firmas digitales, y la confidencialidad mediante un esquema híbrido RSA + AES.

El almacenamiento de claves públicas en los `id managers` no otorga a estos servidores autoridad criptográfica alguna: únicamente funcionan como directorios de última IP conocida y llave pública asociada a cada identificador. El sistema adopta un esquema Trust On First Use (TOFU), por el cual cada cliente cachea localmente la llave pública de cualquier usuario con el que se comunica por primera vez. Una vez asociada una llave a una identidad, esta no puede cambiar sin activar alertas de seguridad locales, evitando ataques de sustitución de claves tras el primer contacto.

### Confidencialidad y autenticidad de los mensajes

Todos los mensajes enviados entre clientes viajan dentro de un contenedor criptográfico denominado `SecureEnvelope`, que combina:

- Cifrado simétrico (AES–256) para proteger la carga útil con eficiencia.
- Envoltorio RSA de la clave AES, utilizando la llave pública del destinatario.
- Firma digital sobre la clave envuelta y el ciphertext.

De esta manera, solo el destinatario que posea la correspondiente llave privada puede descifrar la clave simétrica y recuperar el mensaje. Cualquier modificación del contenido o intento de suplantación se detecta inmediatamente durante la verificación de la firma.

### Seguridad entre ID Managers

Los `id managers`, al intercambiar información para su sincronización eventual, deben autenticarse de forma recíproca para impedir ataques donde un atacante se haga pasar por un servidor legítimo. Para ello, el sistema utiliza un modelo jerárquico con una Autoridad Certificadora (CA) interna:

- Cada `id manager` posee su propio par de llaves RSA.
- Antes de desplegarse, su llave pública es firmada por la CA del proyecto, generando un certificado.
- Cada servidor almacena su llave privada, su certificado firmado y la llave pública de la CA.

Durante el intercambio entre `id managers`, estos realizan un handshake similar a TLS, intercambiando certificados y verificando su validez mediante la CA. Además, ejecutan un desafío criptográfico de prueba de posesión de la llave privada, garantizando que ambos extremos son servidores legítimos y no imitadores. Solo tras esta autenticación mutua se inicia la sincronización de datos.

Este mecanismo asegura que un nodo malicioso no pueda integrarse como `id manager` ni alterar la propagación de estado del sistema.

### Supervivencia ante nodos no confiables

Gracias a la replicación distribuida, incluso si un número significativo de clientes o `id managers` se comportan de forma maliciosa o fallan, el sistema mantiene propiedades esenciales:

- Los clientes nunca confían en datos críticos provenientes de otros nodos sin validarlos criptográficamente.
- Los `id managers` no pueden leer mensajes (que siempre viajan cifrados de extremo a extremo).
- La replicación aleatoria de mensajes pendientes en múltiples clientes reduce la probabilidad de pérdida de información ante fallas o comportamientos adversarios.
- El modelo TOFU mitiga ataques de suplantación después del primer intercambio.

Estos mecanismos permiten que la red continúe operando de forma segura en entornos parcialmente comprometidos, y proporcionan una capa razonable de defensa sin introducir complejidad innecesaria en la arquitectura.
