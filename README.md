# galias

Built to not have to remember dig through my shell history. I personally use it for commands that I use every day, such as
- Tunneling to different environments
- CDing to directories 
- Inpsecting resource such as k8s pod logs

galias operates from a json based config, an example of such a config 
```json
{
   "config":[
      {
         "name":"ssh",
         "command":"ssh ${}@${} ",
         "description": "ssh to an environment",
         "alias":[
            {
               "name":"dev",
               "variables":[
                  "dev-user",
                  "127.0.0.1"
               ]
            },
            {
               "name":"prod",
               "variables":[
                  "prod-user",
                  "127.0.0.2"
               ]
            }
         ]
      },
      {
         "name":"cd",
         "command":"cd ${}",
         "description": "cd to frequently used environments",
         "alias":[
            {
               "name":"viper",
               "description": "spf13/viper"
               "variables":[
                  "~/code/go/viper"
               ]
            },
            {
               "name":"galias",
               "variables":[
                  "~/code/go/galias"
               ]
            }
         ]
      }
   ]
}
```

Usage based on configuration above
```
galias ssh <dev | prod>
galias cd <viper | galias>
```


# Future work
- Add interpolating variables directly in commmandline
- Methods to add new commands and aliases
- Help text for CLI
- Tests# galias
