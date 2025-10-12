# EagleChat

REEEEE 🦅🦅🦅💥💥🦅🦅💥💥🦅

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

## How to run

1 - Add your environment variables `source ./iac/.env`

2 - Build both images:

- `docker build -t eaglechat-client:latest -f iac/client/Dockerfile .`
- `docker build -t eaglechat-id_manager:latest -f iac/id_manager/Dockerfile .`

3 - Init the swarm:

- `docker swarm init`

4 - Deploy the stack:

- `docker stack deploy -c iac/docker-compose.local.yml eaglechat`

With this yo should have the id manager and the client running in containers. For testing you can use.

- `curl http://127.0.0.1:8080/status` for testing the id manager. You should get an `{"status": "ok"}` response
- `docker ps --filter "name=eaglechat_client" --format "{{.Names}}"` will give you the name of the client container, then run `docker attach <name>` to attach to the client's terminal. If everything was done correctly and you have the latest image you should see the tui.
